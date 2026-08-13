package dtos

import (
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type Character struct {
	ID                      uuid.UUID         `json:"id"`
	SessionID               uuid.UUID         `json:"session_id"`
	Name                    string            `json:"name"`
	Age                     int               `json:"age"`
	Profession              string            `json:"profession"`
	Personality             string            `json:"personality"`
	Gender                  string            `json:"gender"`
	Relationships           map[string]string `json:"relationships"`
	Trust                   int               `json:"trust"`
	InterrogationsRemaining int               `json:"interrogations_remaining"`
}

func CharacterFromReadModel(value *readmodels.Character) *Character {
	if value == nil {
		return nil
	}
	return &Character{ID: value.ID, SessionID: value.SessionID, Name: value.Name, Age: value.Age, Profession: value.Profession, Personality: value.Personality, Gender: value.Gender, Relationships: value.Relationships, Trust: value.Trust, InterrogationsRemaining: value.InterrogationsRemaining}
}

func CharactersFromReadModels(values []*readmodels.Character) []*Character {
	result := make([]*Character, len(values))
	for i, value := range values {
		result[i] = CharacterFromReadModel(value)
	}
	return result
}
