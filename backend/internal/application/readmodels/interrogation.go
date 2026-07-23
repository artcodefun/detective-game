package readmodels

import (
	"time"

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
