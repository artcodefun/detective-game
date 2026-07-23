package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID              uuid.UUID `json:"id"`
	SessionID       uuid.UUID `json:"session_id"`
	InterrogationID uuid.UUID `json:"interrogation_id"`
	FromUser        bool      `json:"from_user"`
	Text            string    `json:"text"`
	Statements      []string  `json:"statements,omitempty"`
	AttitudeDelta   int       `json:"attitude_delta"`
	Timestamp       time.Time `json:"timestamp"`
}
