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
	ID          uuid.UUID          `bson:"_id"`
	SessionID   uuid.UUID          `bson:"session_id"`
	CharacterID uuid.UUID          `bson:"character_id"`
	Phase       InterrogationPhase `bson:"phase"`
	CreatedAt   time.Time          `bson:"created_at"`
	CompletedAt *time.Time         `bson:"completed_at,omitempty"`
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
