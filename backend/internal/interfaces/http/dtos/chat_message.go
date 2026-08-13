package dtos

import (
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
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

func ChatMessageFromReadModel(value *readmodels.ChatMessage) *ChatMessage {
	if value == nil {
		return nil
	}
	return &ChatMessage{ID: value.ID, SessionID: value.SessionID, InterrogationID: value.InterrogationID, FromUser: value.FromUser, Text: value.Text, Statements: value.Statements, AttitudeDelta: value.AttitudeDelta, Timestamp: value.Timestamp}
}

func ChatMessagesFromReadModels(values []*readmodels.ChatMessage) []*ChatMessage {
	result := make([]*ChatMessage, len(values))
	for i, value := range values {
		result[i] = ChatMessageFromReadModel(value)
	}
	return result
}
