package readmodels

import "github.com/google/uuid"

type CharacterPrototype struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Profession  string `json:"profession"`
	ImagePath   string `json:"image_path"`
	Personality string `json:"personality"`
	AudioToneID string `json:"audio_tone_id"`
}

type Memory struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	IsTrue    bool      `json:"is_true"`
	Timestamp string    `json:"timestamp"`
}

type CharacterKnowledge struct {
	KnownFacts   []string `json:"known_facts"`
	PartialFacts []string `json:"partial_facts"`
	FalseBeliefs []string `json:"false_beliefs"`
}

type Character struct {
	ID                      uuid.UUID          `json:"id"`
	SessionID               uuid.UUID          `json:"session_id"`
	PrototypeID             int                `json:"prototype_id"`
	Name                    string             `json:"name"`
	Age                     int                `json:"age"`
	Profession              string             `json:"profession"`
	ImagePath               string             `json:"image_path"`
	Personality             string             `json:"personality"`
	AudioToneID             string             `json:"audio_tone_id"`
	Knowledge               CharacterKnowledge `json:"knowledge"`
	Secrets                 []string           `json:"secrets"`
	Relationships           map[string]string  `json:"relationships"`
	Memories                []Memory           `json:"memories"`
	Trust                   int                `json:"trust"`
	InterrogationsRemaining int                `json:"interrogations_remaining"`
}
