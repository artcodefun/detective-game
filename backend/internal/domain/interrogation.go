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
	ID          uuid.UUID          `json:"id" bson:"_id"`
	SessionID   uuid.UUID          `json:"session_id" bson:"session_id"`
	CharacterID uuid.UUID          `json:"character_id" bson:"character_id"`
	Phase       InterrogationPhase `json:"phase" bson:"phase"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	CompletedAt *time.Time         `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
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
