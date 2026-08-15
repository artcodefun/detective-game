package queries

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type SessionQueries struct {
	ReadStore ports.SessionReadRepository
}

func (q *SessionQueries) GetSessionLocale(ctx context.Context, userID, sessionID uuid.UUID) (domain.Locale, error) {
	session, err := q.ReadStore.GetSessionByID(ctx, userID, sessionID)
	if err != nil {
		return "", application.WrapError(err)
	}
	return session.ContentLocale, nil
}

func NewSessionQueries(readStore ports.SessionReadRepository) *SessionQueries {
	return &SessionQueries{ReadStore: readStore}
}

func (q *SessionQueries) GetSession(ctx context.Context, actor application.Actor) (*readmodels.Session, error) {
	s, err := q.ReadStore.GetSession(ctx, actor.UserID)
	return s, application.WrapError(err)
}

func (q *SessionQueries) GetSessionByID(ctx context.Context, actor application.Actor, sessionID uuid.UUID) (*readmodels.Session, error) {
	s, err := q.ReadStore.GetSessionByID(ctx, actor.UserID, sessionID)
	return s, application.WrapError(err)
}

func (q *SessionQueries) ListHistory(ctx context.Context, actor application.Actor) ([]*readmodels.Session, error) {
	h, err := q.ReadStore.ListHistory(ctx, actor.UserID)
	return h, application.WrapError(err)
}

func (q *SessionQueries) GetGameResult(ctx context.Context, actor application.Actor) (*readmodels.GameResult, error) {
	r, err := q.ReadStore.GetGameResult(ctx, actor.SessionID)
	return r, application.WrapError(err)
}
