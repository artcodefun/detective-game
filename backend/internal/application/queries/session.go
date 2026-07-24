package queries

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
)

type SessionQueries struct {
	ReadStore ports.SessionReadRepository
}

func NewSessionQueries(readStore ports.SessionReadRepository) *SessionQueries {
	return &SessionQueries{ReadStore: readStore}
}

func (q *SessionQueries) GetSession(ctx context.Context, actor application.Actor) (*readmodels.Session, error) {
	return q.ReadStore.GetSession(ctx, actor.SessionID)
}

func (q *SessionQueries) ListHistory(ctx context.Context, actor application.Actor) ([]*readmodels.Session, error) {
	return q.ReadStore.ListHistory(ctx, actor.UserID)
}

func (q *SessionQueries) GetGameResult(ctx context.Context, actor application.Actor) (*readmodels.GameResult, error) {
	return q.ReadStore.GetGameResult(ctx, actor.SessionID)
}
