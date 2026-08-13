package queries

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
)

type ChronologyQueries struct {
	Chronology ports.ChronologyReadRepository
}

func NewChronologyQueries(chronology ports.ChronologyReadRepository) *ChronologyQueries {
	return &ChronologyQueries{Chronology: chronology}
}

func (q *ChronologyQueries) GetChronology(ctx context.Context, actor application.Actor) (*readmodels.Chronology, error) {
	c, err := q.Chronology.GetChronology(ctx, actor.SessionID)
	return c, application.WrapError(err)
}
