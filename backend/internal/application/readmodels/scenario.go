package readmodels

import "github.com/google/uuid"

type TimelineEntry struct {
	Time        string     `json:"time"`
	Event       string     `json:"event"`
	CharacterID *uuid.UUID `json:"character_id,omitempty"`
}

type Timeline struct {
	Entries []TimelineEntry `json:"entries"`
}

type Evidence struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	IconAsset           string    `json:"icon_asset"`
	DetailedDescription string    `json:"detailed_description"`
	Type                string    `json:"type"`
}
