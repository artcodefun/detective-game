package readmodels

import (
	"time"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type Session struct {
	ID            uuid.UUID     `json:"id"`
	UserID        uuid.UUID     `json:"user_id"`
	CaseName      string        `json:"case_name"`
	CaseBrief     string        `json:"case_brief"`
	ActionPoints  int           `json:"action_points"`
	Phase         string        `json:"phase"`
	GameResult    *GameResult   `json:"game_result,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	FinishedAt    *time.Time    `json:"finished_at,omitempty"`
	ContentLocale domain.Locale `json:"content_locale"`
}
