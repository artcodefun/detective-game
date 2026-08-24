package queries

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type ChatQueries struct {
	Chat       ports.ChatMessageReadRepository
	Characters ports.CharacterReadRepository
}

func NewChatQueries(chat ports.ChatMessageReadRepository, chars ports.CharacterReadRepository) *ChatQueries {
	return &ChatQueries{Chat: chat, Characters: chars}
}

func (q *ChatQueries) GetChatMessage(ctx context.Context, actor application.Actor, messageID uuid.UUID) (*readmodels.ChatMessage, error) {
	m, err := q.Chat.GetChatMessage(ctx, actor.SessionID, messageID)
	return m, application.WrapError(err)
}

func (q *ChatQueries) ListChatByInterrogation(ctx context.Context, actor application.Actor, interrogationID uuid.UUID) ([]*readmodels.ChatMessage, error) {
	m, err := q.Chat.ListChatByInterrogation(ctx, actor.SessionID, interrogationID)
	return m, application.WrapError(err)
}

func (q *ChatQueries) GetInterrogation(ctx context.Context, actor application.Actor, interrogationID uuid.UUID) (*readmodels.Interrogation, error) {
	i, err := q.Characters.GetInterrogation(ctx, actor.SessionID, interrogationID)
	return i, application.WrapError(err)
}

func (q *ChatQueries) GetActiveInterrogation(ctx context.Context, actor application.Actor) (*readmodels.Interrogation, error) {
	i, err := q.Characters.GetActiveInterrogation(ctx, actor.SessionID)
	return i, application.WrapError(err)
}
