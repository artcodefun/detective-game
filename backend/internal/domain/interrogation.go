package domain

import (
	"time"

	"github.com/google/uuid"
)

type InterrogationPhase string

const (
	InterrogationActive    InterrogationPhase = "active"
	InterrogationCompleted InterrogationPhase = "completed"
)

type Interrogation struct {
	ID          uuid.UUID          `json:"id"`
	SessionID   uuid.UUID          `json:"session_id"`
	CharacterID uuid.UUID          `json:"character_id"`
	Phase       InterrogationPhase `json:"phase"`
	CreatedAt   time.Time          `json:"created_at"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
}

func NewInterrogation(sessionID, characterID uuid.UUID) *Interrogation {
	return &Interrogation{
		ID:          uuid.New(),
		SessionID:   sessionID,
		CharacterID: characterID,
		Phase:       InterrogationActive,
		CreatedAt:   time.Now(),
	}
}

func (i *Interrogation) IsActive() bool {
	return i.Phase == InterrogationActive
}

func (i *Interrogation) Complete() {
	now := time.Now()
	i.Phase = InterrogationCompleted
	i.CompletedAt = &now
}
