package domain

import (
	"time"

	"github.com/google/uuid"
)

type GamePhase string

const (
	GamePhaseInvestigating GamePhase = "investigating"
	GamePhaseFinished      GamePhase = "finished"
)

type Session struct {
	ID            uuid.UUID   `bson:"_id"`
	UserID        uuid.UUID   `bson:"user_id"`
	Crime         Crime       `bson:"crime"`
	Timeline      Timeline    `bson:"timeline"`
	CaseName      string      `bson:"case_name"`
	CaseBrief     string      `bson:"case_brief"`
	ActionPoints  int         `bson:"action_points"`
	Phase         GamePhase   `bson:"phase"`
	GameResult    *GameResult `bson:"game_result,omitempty"`
	CreatedAt     time.Time   `bson:"created_at"`
	FinishedAt    *time.Time  `bson:"finished_at,omitempty"`
	ContentLocale Locale      `bson:"content_locale"`
}

const MaxActionPoints = 5

func (s *Session) CanSpendActionPoint() bool {
	return s.ActionPoints > 0
}

func (s *Session) SpendActionPoints(amount int) bool {
	if s.ActionPoints < amount {
		return false
	}
	s.ActionPoints -= amount
	return true
}

func (s *Session) Finish(result *GameResult) {
	now := time.Now()
	s.GameResult = result
	s.Phase = GamePhaseFinished
	s.FinishedAt = &now
}
