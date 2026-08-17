package ports

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/domain"
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
	Breakdown         domain.ScoreBreakdown
	BreakdownDetails  domain.ScoreBreakdownDetails
	MissedFacts       []string
}

type LlmService interface {
	GenerateScenario(ctx context.Context, locale domain.Locale) (*ScenarioOutput, error)
	RespondInInterrogation(ctx context.Context, locale domain.Locale, character domain.Character, playerMessage string) (*LlmInterrogationResponse, error)
	EvaluateReport(ctx context.Context, locale domain.Locale, playerReport domain.FinalReport, groundTruth domain.Crime) (*LlmFeedbackResponse, error)
	RunAction(ctx context.Context, locale domain.Locale, actionName string, crime domain.Crime, timeline domain.Timeline, evidence *domain.Evidence, character *domain.Character, alibiText *string) (string, error)
}
