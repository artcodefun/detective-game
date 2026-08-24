package commands

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

func requireActiveSession(ctx context.Context, sessions ports.SessionRepository, sessionID uuid.UUID) (*domain.Session, error) {
	session, err := sessions.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.Phase == domain.GamePhaseFinished {
		return nil, application.NewAppError(application.KindConflict, domain.T("error.session_already_finished"))
	}
	return session, nil
}
