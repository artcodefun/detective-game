package dtos

import (
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type Evidence struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	IconAsset           string    `json:"icon_asset"`
	DetailedDescription string    `json:"detailed_description"`
	Type                string    `json:"type"`
}

func EvidenceFromReadModel(value *readmodels.Evidence) *Evidence {
	if value == nil {
		return nil
	}
	return &Evidence{ID: value.ID, Name: value.Name, Description: value.Description, IconAsset: value.IconAsset, DetailedDescription: value.DetailedDescription, Type: value.Type}
}

func EvidenceListFromReadModels(values []*readmodels.Evidence) []*Evidence {
	result := make([]*Evidence, len(values))
	for i, value := range values {
		result[i] = EvidenceFromReadModel(value)
	}
	return result
}
