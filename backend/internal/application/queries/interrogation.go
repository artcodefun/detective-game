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
	return q.Chat.GetChatMessage(ctx, messageID)
}

func (q *ChatQueries) ListChatByInterrogation(ctx context.Context, actor application.Actor, interrogationID uuid.UUID) ([]*readmodels.ChatMessage, error) {
	return q.Chat.ListChatByInterrogation(ctx, interrogationID)
}

func (q *ChatQueries) GetInterrogation(ctx context.Context, actor application.Actor, interrogationID uuid.UUID) (*readmodels.Interrogation, error) {
	return q.Characters.GetInterrogation(ctx, interrogationID)
}
