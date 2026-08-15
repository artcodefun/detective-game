package dtos

import (
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/artcodefun/detective-game/backend/internal/domain"
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

func ChronologyEntriesFromReadModel(entriesModel []readmodels.ChronologyEntry, translator ports.Translator, locale domain.Locale) []ChronologyEntry {
	entries := make([]ChronologyEntry, len(entriesModel))
	for i, entry := range entriesModel {
		details := make([]NotebookEntry, len(entry.Details))
		for j, detail := range entry.Details {
			details[j] = NotebookEntry{ID: detail.ID, Type: detail.Type, CharacterID: detail.CharacterID, Description: detail.Description, UserTags: detail.UserTags, UserNote: detail.UserNote, Timestamp: detail.Timestamp}
		}
		entries[i] = ChronologyEntry{ID: entry.ID, OriginID: entry.OriginID, EventType: entry.EventType, Title: translator.Translate(locale, entry.Title), Timestamp: entry.Timestamp, Details: details}
	}
	return entries
}
