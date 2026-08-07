package queries

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type SessionQueries struct {
	ReadStore ports.SessionReadRepository
}

func NewSessionQueries(readStore ports.SessionReadRepository) *SessionQueries {
	return &SessionQueries{ReadStore: readStore}
}

func (q *SessionQueries) GetSession(ctx context.Context, actor application.Actor) (*readmodels.Session, error) {
	return q.ReadStore.GetSession(ctx, actor.UserID)
}

func (q *SessionQueries) GetSessionByID(ctx context.Context, actor application.Actor, sessionID uuid.UUID) (*readmodels.Session, error) {
	return q.ReadStore.GetSessionByID(ctx, actor.UserID, sessionID)
}

func (q *SessionQueries) ListHistory(ctx context.Context, actor application.Actor) ([]*readmodels.Session, error) {
	return q.ReadStore.ListHistory(ctx, actor.UserID)
}

func (q *SessionQueries) GetGameResult(ctx context.Context, actor application.Actor) (*readmodels.GameResult, error) {
	return q.ReadStore.GetGameResult(ctx, actor.SessionID)
}
