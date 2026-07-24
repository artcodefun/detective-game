package domain

import "github.com/google/uuid"

type CharacterPrototype struct {
	ID          int    `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Age         int    `json:"age" bson:"age"`
	Profession  string `json:"profession" bson:"profession"`
	ImagePath   string `json:"image_path" bson:"image_path"`
	Personality string `json:"personality" bson:"personality"`
	AudioToneID string `json:"audio_tone_id" bson:"audio_tone_id"`
}

type TrustLevel string

const (
	TrustLevelOpen     TrustLevel = "open"
	TrustLevelReserved TrustLevel = "reserved"
	TrustLevelTense    TrustLevel = "tense"
	TrustLevelClosed   TrustLevel = "closed"
)

type Memory struct {
	ID        uuid.UUID `json:"id" bson:"id"`
	Content   string    `json:"content" bson:"content"`
	IsTrue    bool      `json:"is_true" bson:"is_true"`
	Timestamp string    `json:"timestamp" bson:"timestamp"`
}

type CharacterKnowledge struct {
	KnownFacts   []string `json:"known_facts" bson:"known_facts"`
	PartialFacts []string `json:"partial_facts" bson:"partial_facts"`
	FalseBeliefs []string `json:"false_beliefs" bson:"false_beliefs"`
}

type Character struct {
	ID                      uuid.UUID          `json:"id" bson:"_id"`
	SessionID               uuid.UUID          `json:"session_id" bson:"session_id"`
	Prototype               CharacterPrototype `json:"prototype" bson:"prototype"`
	Knowledge               CharacterKnowledge `json:"knowledge" bson:"knowledge"`
	Secrets                 []string           `json:"secrets" bson:"secrets"`
	Relationships           map[string]string  `json:"relationships" bson:"relationships"`
	Memories                []Memory           `json:"memories" bson:"memories"`
	Trust                   int                `json:"trust" bson:"trust"`
	InterrogationsRemaining int                `json:"interrogations_remaining" bson:"interrogations_remaining"`
}

const (
	MinTrust          = 0
	MaxTrust          = 100
	MaxInterrogations = 3
)

func NewCharacter(prototype CharacterPrototype, sessionID uuid.UUID) Character {
	return Character{
		ID:                      uuid.New(),
		SessionID:               sessionID,
		Prototype:               prototype,
		Trust:                   50,
		InterrogationsRemaining: MaxInterrogations,
	}
}

func (c *Character) TrustLevel() TrustLevel {
	switch {
	case c.Trust >= 75:
		return TrustLevelOpen
	case c.Trust >= 50:
		return TrustLevelReserved
	case c.Trust >= 25:
		return TrustLevelTense
	default:
		return TrustLevelClosed
	}
}

func (c *Character) CanInterrogate() bool {
	return c.InterrogationsRemaining > 0
}

func (c *Character) ApplyAttitudeDelta(delta int) {
	c.Trust = clamp(c.Trust+delta, MinTrust, MaxTrust)
}

func (c *Character) DecrementInterrogation() {
	if c.InterrogationsRemaining > 0 {
		c.InterrogationsRemaining--
	}
}

func (c *Character) AssignToSession(sessionID uuid.UUID) {
	c.SessionID = sessionID
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
