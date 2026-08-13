package dtos

import "github.com/artcodefun/detective-game/backend/internal/application/readmodels"

type GameResult struct {
	PlayerReport      FinalReport       `json:"player_report"`
	Breakdown         ScoreBreakdown    `json:"breakdown"`
	NarrativeFeedback string            `json:"narrative_feedback"`
	BreakdownDetails  map[string]string `json:"breakdown_details"`
}

type FinalReport struct {
	Who      string `json:"who"`
	Why      string `json:"why"`
	How      string `json:"how"`
	When     string `json:"when"`
	Evidence string `json:"evidence"`
}

type ScoreBreakdown struct {
	WhoCorrect      bool `json:"who_correct"`
	WhyCorrect      bool `json:"why_correct"`
	HowCorrect      bool `json:"how_correct"`
	WhenCorrect     bool `json:"when_correct"`
	EvidenceCorrect bool `json:"evidence_correct"`
}

func GameResultFromReadModel(value *readmodels.GameResult) *GameResult {
	if value == nil {
		return nil
	}
	return &GameResult{PlayerReport: FinalReport{Who: value.PlayerReport.Who, Why: value.PlayerReport.Why, How: value.PlayerReport.How, When: value.PlayerReport.When, Evidence: value.PlayerReport.Evidence}, Breakdown: ScoreBreakdown{WhoCorrect: value.Breakdown.WhoCorrect, WhyCorrect: value.Breakdown.WhyCorrect, HowCorrect: value.Breakdown.HowCorrect, WhenCorrect: value.Breakdown.WhenCorrect, EvidenceCorrect: value.Breakdown.EvidenceCorrect}, NarrativeFeedback: value.NarrativeFeedback, BreakdownDetails: value.BreakdownDetails}
}
