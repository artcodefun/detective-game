package dtos

import (
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type Session struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"user_id"`
	CaseName      string      `json:"case_name"`
	CaseBrief     string      `json:"case_brief"`
	ActionPoints  int         `json:"action_points"`
	Phase         string      `json:"phase"`
	GameResult    *GameResult `json:"game_result,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	FinishedAt    *time.Time  `json:"finished_at,omitempty"`
	ContentLocale string      `json:"content_locale"`
}

func SessionFromReadModel(value *readmodels.Session) *Session {
	if value == nil {
		return nil
	}
	return &Session{ID: value.ID, UserID: value.UserID, CaseName: value.CaseName, CaseBrief: value.CaseBrief, ActionPoints: value.ActionPoints, Phase: value.Phase, GameResult: GameResultFromReadModel(value.GameResult), CreatedAt: value.CreatedAt, FinishedAt: value.FinishedAt, ContentLocale: string(value.ContentLocale)}
}

func SessionsFromReadModels(values []*readmodels.Session) []*Session {
	result := make([]*Session, len(values))
	for i, value := range values {
		result[i] = SessionFromReadModel(value)
	}
	return result
}
