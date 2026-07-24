package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID              uuid.UUID `json:"id" bson:"_id"`
	SessionID       uuid.UUID `json:"session_id" bson:"session_id"`
	InterrogationID uuid.UUID `json:"interrogation_id" bson:"interrogation_id"`
	FromUser        bool      `json:"from_user" bson:"from_user"`
	Text            string    `json:"text" bson:"text"`
	Statements      []string  `json:"statements,omitempty" bson:"statements,omitempty"`
	AttitudeDelta   int       `json:"attitude_delta" bson:"attitude_delta"`
	Timestamp       time.Time `json:"timestamp" bson:"timestamp"`
}
