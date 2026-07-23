package readmodels

import (
	"time"

	"github.com/google/uuid"
)

type ChronologyEntry struct {
	ID        uuid.UUID       `json:"id"`
	EventType string          `json:"event_type"`
	Title     string          `json:"title"`
	Timestamp time.Time       `json:"timestamp"`
	Details   []NotebookEntry `json:"details"`
}

type NotebookEntry struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	CharacterID *int      `json:"character_id,omitempty"`
	Description string    `json:"description"`
	UserTags    []string  `json:"user_tags"`
	UserNote    *string   `json:"user_note,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

type ActionReport struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Body        string     `json:"body"`
	EvidenceID  *uuid.UUID `json:"evidence_id,omitempty"`
	CharacterID *uuid.UUID `json:"character_id,omitempty"`
	Timestamp   time.Time  `json:"timestamp"`
}
