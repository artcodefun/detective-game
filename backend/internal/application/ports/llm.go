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
	GenerateScenario(ctx context.Context) (*ScenarioOutput, error)
	RespondInInterrogation(ctx context.Context, character domain.Character, playerMessage string) (*LlmInterrogationResponse, error)
	EvaluateReport(ctx context.Context, playerReport domain.FinalReport, groundTruth domain.Crime) (*LlmFeedbackResponse, error)
	RunAction(ctx context.Context, actionName string, evidenceID *uuid.UUID, characterID *uuid.UUID, alibiText *string) (string, error)
}
