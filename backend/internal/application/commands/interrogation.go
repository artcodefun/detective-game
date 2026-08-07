package commands

import (
	"context"
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
}

func NewInterrogationCommands(sessions ports.SessionRepository, interrogations ports.InterrogationRepository, chars ports.CharacterRepository, chat ports.ChatMessageRepository, llm ports.LlmService) *InterrogationCommands {
	return &InterrogationCommands{Sessions: sessions, Interrogations: interrogations, Characters: chars, Chat: chat, LLM: llm}
}

func (c *InterrogationCommands) Create(ctx context.Context, actor application.Actor, characterID uuid.UUID) (*domain.Interrogation, error) {
	active, err := c.Interrogations.FindActiveBySession(ctx, actor.SessionID)
	if err != nil {
		return nil, application.WrapError(err)
	}
	if active != nil {
		return nil, application.NewAppError(application.KindConflict, "active_interrogation_exists")
	}

	char, err := c.Characters.FindCharacterByID(ctx, actor.SessionID, characterID)
	if err != nil {
		return nil, application.WrapError(err)
	}
	if !char.CanInterrogate() {
		return nil, application.NewAppError(application.KindConflict, "no_interrogations_left")
	}

	char.DecrementInterrogation()
	if err := c.Characters.UpdateCharacter(ctx, char); err != nil {
		return nil, application.WrapError(err)
	}

	inter := domain.NewInterrogation(actor.SessionID, characterID)
	if err := c.Interrogations.CreateInterrogation(ctx, inter); err != nil {
		return nil, application.WrapError(err)
	}

	return inter, nil
}

func (c *InterrogationCommands) AddMessage(ctx context.Context, actor application.Actor, interrogationID uuid.UUID, message string) (uuid.UUID, error) {
	inter, err := c.Interrogations.FindInterrogationByID(ctx, interrogationID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	if !inter.IsActive() {
		return uuid.Nil, application.NewAppError(application.KindConflict, "interrogation_not_active")
	}

	char, err := c.Characters.FindCharacterByID(ctx, actor.SessionID, inter.CharacterID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	playerMsg := domain.ChatMessage{
		ID:              uuid.New(),
		SessionID:       actor.SessionID,
		InterrogationID: interrogationID,
		FromUser:        true,
		Text:            message,
		Timestamp:       time.Now(),
	}
	if err := c.Chat.AppendChatMessage(ctx, &playerMsg); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	resp, err := c.LLM.RespondInInterrogation(ctx, *char, message)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	char.ApplyAttitudeDelta(resp.AttitudeDelta)

	npcMsg := domain.ChatMessage{
		ID:              uuid.New(),
		SessionID:       actor.SessionID,
		InterrogationID: interrogationID,
		FromUser:        false,
		Text:            resp.Answer,
		Statements:      resp.Statements,
		AttitudeDelta:   resp.AttitudeDelta,
		Timestamp:       time.Now(),
	}
	if err := c.Chat.AppendChatMessage(ctx, &npcMsg); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	if err := c.Characters.UpdateCharacter(ctx, char); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	return npcMsg.ID, nil
}

func (c *InterrogationCommands) Complete(ctx context.Context, actor application.Actor, interrogationID uuid.UUID) error {
	inter, err := c.Interrogations.FindInterrogationByID(ctx, interrogationID)
	if err != nil {
		return application.WrapError(err)
	}
	if !inter.IsActive() {
		return application.NewAppError(application.KindConflict, "interrogation_not_active")
	}

	inter.Complete()
	return application.WrapError(c.Interrogations.UpdateInterrogation(ctx, inter))
}
