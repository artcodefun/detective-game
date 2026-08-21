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
}

func NewInterrogationCommands(sessions ports.SessionRepository, interrogations ports.InterrogationRepository, chars ports.CharacterRepository, chat ports.ChatMessageRepository, llm ports.LlmService, chronology ports.ChronologyRepository) *InterrogationCommands {
	return &InterrogationCommands{Sessions: sessions, Interrogations: interrogations, Characters: chars, Chat: chat, LLM: llm, Chronology: chronology}
}

func (c *InterrogationCommands) Create(ctx context.Context, actor application.Actor, characterID uuid.UUID) (uuid.UUID, error) {
	session, err := c.Sessions.FindByID(ctx, actor.SessionID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	if session.Phase == domain.GamePhaseFinished {
		return uuid.Nil, application.NewAppError(application.KindConflict, domain.T("error.session_already_finished"))
	}

	active, err := c.Interrogations.FindActiveBySession(ctx, actor.SessionID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	if active != nil {
		return uuid.Nil, application.NewAppError(application.KindConflict, domain.T("error.active_interrogation_exists"))
	}

	char, err := c.Characters.FindCharacterByID(ctx, actor.SessionID, characterID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	if !char.CanInterrogate() {
		return uuid.Nil, application.NewAppError(application.KindConflict, domain.T("error.no_interrogations_left"))
	}

	char.DecrementInterrogation()
	if err := c.Characters.UpdateCharacter(ctx, char); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	inter := domain.NewInterrogation(actor.SessionID, characterID)
	if err := c.Interrogations.CreateInterrogation(ctx, inter); err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	chronology := domain.NewChronologyEntry(domain.ChronologyEventTypeInterrogation, &inter.ID, domain.TWith("chronology.interrogation_started", map[string]any{"character_name": char.Name}), inter.CreatedAt)
	if err := c.Chronology.AppendChronologyEntry(ctx, actor.SessionID, chronology); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	return inter.ID, nil
}

func (c *InterrogationCommands) AddMessage(ctx context.Context, actor application.Actor, interrogationID uuid.UUID, message string) (uuid.UUID, error) {
	message = strings.TrimSpace(message)
	if !domain.ValidateInterrogationQuestion(message) {
		return uuid.Nil, application.NewAppError(application.KindInvalidInput, domain.T("error.invalid_interrogation_question"))
	}

	session, err := c.Sessions.FindByID(ctx, actor.SessionID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	if session.Phase == domain.GamePhaseFinished {
		return uuid.Nil, application.NewAppError(application.KindConflict, domain.T("error.session_already_finished"))
	}

	inter, err := c.Interrogations.FindInterrogationByID(ctx, interrogationID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	if !inter.IsActive() {
		return uuid.Nil, application.NewAppError(application.KindConflict, domain.T("error.interrogation_not_active"))
	}
	messages, err := c.Chat.FindChatByInterrogation(ctx, interrogationID)
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
	if err := c.Chat.AppendChatMessage(ctx, &playerMsg); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	char.ApplyAttitudeDelta(resp.AttitudeDelta)
	char.Memories = append(char.Memories,
		domain.Memory{
			Content:   "Детектив спросил: " + message,
			Timestamp: time.Now().Format(time.RFC3339),
		},
		domain.Memory{
			Content:   char.Name + " ответил: " + resp.Answer,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	)

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
	if err := c.Chat.AppendChatMessage(ctx, &npcMsg); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	if err := c.Characters.UpdateCharacter(ctx, char); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	details := make([]domain.NotebookEntry, 0, len(resp.Statements))
	for _, statement := range resp.Statements {
		details = append(details, domain.NotebookEntry{
			ID:          uuid.New(),
			Type:        domain.NotebookEntryTypeStatement,
			CharacterID: &inter.CharacterID,
			Description: statement,
			Timestamp:   time.Now(),
		})
	}
	if len(details) > 0 {
		if err := c.Chronology.AppendNotebookEntries(ctx, actor.SessionID, interrogationID, details); err != nil {
			return uuid.Nil, application.WrapError(err)
		}
	}
	if inter.ShouldCompleteAfterQuestion(messages) {
		inter.Complete()
		if err := c.Interrogations.UpdateInterrogation(ctx, inter); err != nil {
			return uuid.Nil, application.WrapError(err)
		}
	}

	return npcMsg.ID, nil
}

func (c *InterrogationCommands) Complete(ctx context.Context, actor application.Actor, interrogationID uuid.UUID) error {
	inter, err := c.Interrogations.FindInterrogationByID(ctx, interrogationID)
	if err != nil {
		return application.WrapError(err)
	}
	if !inter.IsActive() {
		return application.NewAppError(application.KindConflict, domain.T("error.interrogation_not_active"))
	}

	inter.Complete()
	return application.WrapError(c.Interrogations.UpdateInterrogation(ctx, inter))
}
