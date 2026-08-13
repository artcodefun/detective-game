package application

import (
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type Actor struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	Locale    domain.Locale
}
