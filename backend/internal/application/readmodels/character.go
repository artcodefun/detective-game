package readmodels

import "github.com/google/uuid"

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
