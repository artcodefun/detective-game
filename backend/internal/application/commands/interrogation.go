package commands

import (
	"context"
	"strings"
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type InterrogationCommands struct {
	Sessions       ports.SessionRepository
	Interrogations ports.InterrogationRepository
	Characters     ports.CharacterRepository
	Chat           ports.ChatMessageRepository
	LLM            ports.LlmService
	Chronology     ports.ChronologyRepository
	TxMgr          ports.TransactionManager
}

func NewInterrogationCommands(sessions ports.SessionRepository, interrogations ports.InterrogationRepository, chars ports.CharacterRepository, chat ports.ChatMessageRepository, llm ports.LlmService, chronology ports.ChronologyRepository, txMgr ports.TransactionManager) *InterrogationCommands {
	return &InterrogationCommands{Sessions: sessions, Interrogations: interrogations, Characters: chars, Chat: chat, LLM: llm, Chronology: chronology, TxMgr: txMgr}
}

func (c *InterrogationCommands) Create(ctx context.Context, actor application.Actor, characterID uuid.UUID) (uuid.UUID, error) {
	inter := domain.NewInterrogation(actor.SessionID, characterID)
	if err := c.TxMgr.WithTx(ctx, func(txCtx context.Context) error {
		_, err := requireActiveSession(txCtx, c.Sessions, actor.SessionID)
		if err != nil {
			return err
		}

		active, err := c.Interrogations.FindActiveBySession(txCtx, actor.SessionID)
		if err != nil {
			return err
		}
		if active != nil {
			return application.NewAppError(application.KindConflict, domain.T("error.active_interrogation_exists"))
		}

		char, err := c.Characters.FindCharacterByID(txCtx, actor.SessionID, characterID)
		if err != nil {
			return err
		}
		if !char.CanInterrogate() {
			return application.NewAppError(application.KindConflict, domain.T("error.no_interrogations_left"))
		}

		char.DecrementInterrogation()
		if err := c.Characters.UpdateCharacter(txCtx, char); err != nil {
			return err
		}
		if err := c.Interrogations.CreateInterrogation(txCtx, inter); err != nil {
			return err
		}
		chronology := domain.NewChronologyEntry(domain.ChronologyEventTypeInterrogation, &inter.ID, domain.TWith("chronology.interrogation_started", map[string]any{"character_name": char.Name}), inter.CreatedAt)
		return c.Chronology.AppendChronologyEntry(txCtx, actor.SessionID, chronology)
	}); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	return inter.ID, nil
}

func (c *InterrogationCommands) AddMessage(ctx context.Context, actor application.Actor, interrogationID uuid.UUID, message string) (uuid.UUID, error) {
	message = strings.TrimSpace(message)
	if !domain.ValidateInterrogationQuestion(message) {
		return uuid.Nil, application.NewAppError(application.KindInvalidInput, domain.T("error.invalid_interrogation_question"))
	}

	_, err := requireActiveSession(ctx, c.Sessions, actor.SessionID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	inter, err := c.Interrogations.FindInterrogationByID(ctx, actor.SessionID, interrogationID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	if !inter.IsActive() {
		return uuid.Nil, application.NewAppError(application.KindConflict, domain.T("error.interrogation_not_active"))
	}
	messages, err := c.Chat.FindChatByInterrogation(ctx, actor.SessionID, interrogationID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	if !inter.CanAskQuestion(messages) {
		return uuid.Nil, application.NewAppError(application.KindConflict, domain.T("error.interrogation_question_limit_reached"))
	}

	char, err := c.Characters.FindCharacterByID(ctx, actor.SessionID, inter.CharacterID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	resp, err := c.LLM.RespondInInterrogation(ctx, actor.SessionContentLocale, *char, message)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	now := time.Now()
	playerMsg := domain.ChatMessage{
		ID:              uuid.New(),
		SessionID:       actor.SessionID,
		InterrogationID: interrogationID,
		FromUser:        true,
		Text:            message,
		Timestamp:       now,
	}
	npcMsg := domain.ChatMessage{
		ID:              uuid.New(),
		SessionID:       actor.SessionID,
		InterrogationID: interrogationID,
		FromUser:        false,
		Text:            resp.Answer,
		Statements:      resp.Statements,
		AttitudeDelta:   resp.AttitudeDelta,
		Timestamp:       now,
	}
	if err := c.TxMgr.WithTx(ctx, func(txCtx context.Context) error {
		_, err := requireActiveSession(txCtx, c.Sessions, actor.SessionID)
		if err != nil {
			return err
		}
		currentInter, err := c.Interrogations.FindInterrogationByID(txCtx, actor.SessionID, interrogationID)
		if err != nil {
			return err
		}
		messages, err := c.Chat.FindChatByInterrogation(txCtx, actor.SessionID, interrogationID)
		if err != nil {
			return err
		}
		if !currentInter.CanAskQuestion(messages) {
			return application.NewAppError(application.KindConflict, domain.T("error.interrogation_question_limit_reached"))
		}
		currentChar, err := c.Characters.FindCharacterByID(txCtx, actor.SessionID, currentInter.CharacterID)
		if err != nil {
			return err
		}
		if err := c.Chat.AppendChatMessage(txCtx, &playerMsg); err != nil {
			return err
		}
		if err := c.Chat.AppendChatMessage(txCtx, &npcMsg); err != nil {
			return err
		}
		currentChar.ApplyAttitudeDelta(resp.AttitudeDelta)
		currentChar.Memories = append(currentChar.Memories, domain.Memory{Content: "Детектив спросил: " + message, Timestamp: now.Format(time.RFC3339)}, domain.Memory{Content: currentChar.Name + " ответил: " + resp.Answer, Timestamp: now.Format(time.RFC3339)})
		if err := c.Characters.UpdateCharacter(txCtx, currentChar); err != nil {
			return err
		}
		details := make([]domain.NotebookEntry, 0, len(resp.Statements))
		for _, statement := range resp.Statements {
			details = append(details, domain.NotebookEntry{ID: uuid.New(), Type: domain.NotebookEntryTypeStatement, CharacterID: &currentInter.CharacterID, Description: statement, Timestamp: now})
		}
		if len(details) > 0 {
			if err := c.Chronology.AppendNotebookEntries(txCtx, actor.SessionID, interrogationID, details); err != nil {
				return err
			}
		}
		if currentInter.ShouldCompleteAfterQuestion(messages) {
			currentInter.Complete()
			return c.Interrogations.UpdateInterrogation(txCtx, currentInter)
		}
		return nil
	}); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	return npcMsg.ID, nil
}

func (c *InterrogationCommands) Complete(ctx context.Context, actor application.Actor, interrogationID uuid.UUID) error {
	inter, err := c.Interrogations.FindInterrogationByID(ctx, actor.SessionID, interrogationID)
	if err != nil {
		return application.WrapError(err)
	}
	if !inter.IsActive() {
		return application.NewAppError(application.KindConflict, domain.T("error.interrogation_not_active"))
	}

	inter.Complete()
	return application.WrapError(c.Interrogations.UpdateInterrogation(ctx, inter))
}
