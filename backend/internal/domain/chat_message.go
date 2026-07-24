package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID              uuid.UUID `bson:"_id"`
	SessionID       uuid.UUID `bson:"session_id"`
	InterrogationID uuid.UUID `bson:"interrogation_id"`
	FromUser        bool      `bson:"from_user"`
	Text            string    `bson:"text"`
	Statements      []string  `bson:"statements,omitempty"`
	AttitudeDelta   int       `bson:"attitude_delta"`
	Timestamp       time.Time `bson:"timestamp"`
}
