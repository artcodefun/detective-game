package dtos

import (
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type ChronologyEntry struct {
	ID        uuid.UUID       `json:"id"`
	OriginID  *uuid.UUID      `json:"origin_id,omitempty"`
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

func ChronologyEntriesFromReadModel(chronology *readmodels.Chronology, translator ports.Translator) []ChronologyEntry {
	entries := make([]ChronologyEntry, len(chronology.Entries))
	for i, entry := range chronology.Entries {
		details := make([]NotebookEntry, len(entry.Details))
		for j, detail := range entry.Details {
			details[j] = NotebookEntry{ID: detail.ID, Type: detail.Type, CharacterID: detail.CharacterID, Description: detail.Description, UserTags: detail.UserTags, UserNote: detail.UserNote, Timestamp: detail.Timestamp}
		}
		entries[i] = ChronologyEntry{ID: entry.ID, OriginID: entry.OriginID, EventType: entry.EventType, Title: translator.Translate(chronology.ContentLocale, entry.Title), Timestamp: entry.Timestamp, Details: details}
	}
	return entries
}
