package dtos

import (
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type Interrogation struct {
	ID          uuid.UUID  `json:"id"`
	SessionID   uuid.UUID  `json:"session_id"`
	CharacterID uuid.UUID  `json:"character_id"`
	Phase       string     `json:"phase"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func InterrogationFromReadModel(value *readmodels.Interrogation) *Interrogation {
	if value == nil {
		return nil
	}
	return &Interrogation{ID: value.ID, SessionID: value.SessionID, CharacterID: value.CharacterID, Phase: value.Phase, CreatedAt: value.CreatedAt, CompletedAt: value.CompletedAt}
}
