package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type InterrogationPhase string

const (
	InterrogationActive    InterrogationPhase = "active"
	InterrogationCompleted InterrogationPhase = "completed"

	MaxInterrogationQuestions      = 10
	MaxInterrogationQuestionLength = 600
)

type Interrogation struct {
	ID          uuid.UUID          `bson:"_id"`
	SessionID   uuid.UUID          `bson:"session_id"`
	CharacterID uuid.UUID          `bson:"character_id"`
	Phase       InterrogationPhase `bson:"phase"`
	CreatedAt   time.Time          `bson:"created_at"`
	CompletedAt *time.Time         `bson:"completed_at,omitempty"`
}

func NewInterrogation(sessionID, characterID uuid.UUID) *Interrogation {
	return &Interrogation{
		ID:          uuid.New(),
		SessionID:   sessionID,
		CharacterID: characterID,
		Phase:       InterrogationActive,
		CreatedAt:   time.Now(),
	}
}

func (i *Interrogation) IsActive() bool {
	return i.Phase == InterrogationActive
}

func (i *Interrogation) QuestionsAsked(messages []*ChatMessage) int {
	count := 0
	for _, message := range messages {
		if message.FromUser {
			count++
		}
	}
	return count
}

func (i *Interrogation) CanAskQuestion(messages []*ChatMessage) bool {
	return i.IsActive() && i.QuestionsAsked(messages) < MaxInterrogationQuestions
}

func (i *Interrogation) ShouldCompleteAfterQuestion(messages []*ChatMessage) bool {
	return i.QuestionsAsked(messages)+1 >= MaxInterrogationQuestions
}

func (i *Interrogation) Complete() {
	now := time.Now()
	i.Phase = InterrogationCompleted
	i.CompletedAt = &now
}

func ValidateInterrogationQuestion(question string) bool {
	return strings.TrimSpace(question) != "" && utf8.RuneCountInString(question) <= MaxInterrogationQuestionLength
}
