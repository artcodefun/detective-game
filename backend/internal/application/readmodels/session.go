package readmodels

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID   `json:"id"`
	UserID       uuid.UUID   `json:"user_id"`
	Crime        Crime       `json:"crime"`
	Timeline     Timeline    `json:"timeline"`
	CaseName     string      `json:"case_name"`
	CaseBrief    string      `json:"case_brief"`
	ActionPoints int         `json:"action_points"`
	Phase        string      `json:"phase"`
	GameResult   *GameResult `json:"game_result,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
}
