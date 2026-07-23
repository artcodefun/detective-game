package queries

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type CharacterQueries struct {
	Characters ports.CharacterReadRepository
}

func NewCharacterQueries(chars ports.CharacterReadRepository) *CharacterQueries {
	return &CharacterQueries{Characters: chars}
}

func (q *CharacterQueries) ListCharacters(ctx context.Context, actor application.Actor) ([]*readmodels.Character, error) {
	return q.Characters.ListCharacters(ctx, actor.SessionID)
}

func (q *CharacterQueries) GetCharacter(ctx context.Context, actor application.Actor, characterID uuid.UUID) (*readmodels.Character, error) {
	return q.Characters.GetCharacter(ctx, actor.SessionID, characterID)
}
