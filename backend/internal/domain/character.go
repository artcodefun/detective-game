package domain

import "github.com/google/uuid"

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type TrustLevel string

const (
	TrustLevelOpen     TrustLevel = "open"
	TrustLevelReserved TrustLevel = "reserved"
	TrustLevelTense    TrustLevel = "tense"
	TrustLevelClosed   TrustLevel = "closed"
)

type Memory struct {
	ID        uuid.UUID `bson:"id"`
	Content   string    `bson:"content"`
	IsTrue    bool      `bson:"is_true"`
	Timestamp string    `bson:"timestamp"`
}

type CharacterKnowledge struct {
	KnownFacts   []string `bson:"known_facts"`
	PartialFacts []string `bson:"partial_facts"`
	FalseBeliefs []string `bson:"false_beliefs"`
}

type Character struct {
	ID                      uuid.UUID          `bson:"_id"`
	SessionID               uuid.UUID          `bson:"session_id"`
	Name                    string             `bson:"name"`
	Age                     int                `bson:"age"`
	Profession              string             `bson:"profession"`
	Personality             string             `bson:"personality"`
	Gender                  Gender             `bson:"gender"`
	Knowledge               CharacterKnowledge `bson:"knowledge"`
	Secrets                 []string           `bson:"secrets"`
	Relationships           map[string]string  `bson:"relationships"`
	Memories                []Memory           `bson:"memories"`
	Trust                   int                `bson:"trust"`
	InterrogationsRemaining int                `bson:"interrogations_remaining"`
}

const (
	MinTrust          = 0
	MaxTrust          = 100
	MaxInterrogations = 3
)

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
