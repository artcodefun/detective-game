package domain

import (
	"time"

	"github.com/google/uuid"
)

type GamePhase string

const (
	GamePhaseIdle          GamePhase = "idle"
	GamePhaseGenerating    GamePhase = "generating"
	GamePhaseInvestigating GamePhase = "investigating"
	GamePhaseWritingReport GamePhase = "writing_report"
	GamePhaseFinished      GamePhase = "finished"
)

type Session struct {
	ID           uuid.UUID   `json:"id"`
	Crime        Crime       `json:"crime"`
	Timeline     Timeline    `json:"timeline"`
	CaseName     string      `json:"case_name"`
	CaseBrief    string      `json:"case_brief"`
	ActionPoints int         `json:"action_points"`
	Phase        GamePhase   `json:"phase"`
	GameResult   *GameResult `json:"game_result,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
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
	s.GameResult = result
	s.Phase = GamePhaseFinished
}
