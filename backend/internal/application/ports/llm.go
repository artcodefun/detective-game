package ports

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type ScenarioOutput struct {
	Crime      domain.Crime
	Timeline   domain.Timeline
	CaseName   string
	CaseBrief  string
	Characters []domain.Character
	Evidence   []domain.Evidence
}

type LlmInterrogationResponse struct {
	Answer        string
	AttitudeDelta int
	Statements    []string
}

type LlmFeedbackResponse struct {
	NarrativeFeedback string
	BreakdownDetails  map[string]string
	MissedFacts       []string
}

type LlmService interface {
	GenerateScenario(ctx context.Context, locale domain.Locale) (*ScenarioOutput, error)
	RespondInInterrogation(ctx context.Context, locale domain.Locale, character domain.Character, playerMessage string) (*LlmInterrogationResponse, error)
	EvaluateReport(ctx context.Context, locale domain.Locale, playerReport domain.FinalReport, groundTruth domain.Crime) (*LlmFeedbackResponse, error)
	RunAction(ctx context.Context, locale domain.Locale, actionName string, evidenceID *uuid.UUID, characterID *uuid.UUID, alibiText *string) (string, error)
}
