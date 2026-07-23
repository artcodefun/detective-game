package commands

import (
	"context"
	"strings"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
)

type EvaluationCommands struct {
	Sessions ports.SessionRepository
	LLM      ports.LlmService
}

func NewEvaluationCommands(sessions ports.SessionRepository, llm ports.LlmService) *EvaluationCommands {
	return &EvaluationCommands{Sessions: sessions, LLM: llm}
}

func (c *EvaluationCommands) SubmitReport(ctx context.Context, actor application.Actor, report domain.FinalReport) error {
	session, err := c.Sessions.FindByID(ctx, actor.SessionID)
	if err != nil {
		return err
	}
	if session.Phase == domain.GamePhaseFinished {
		return application.NewAppError(application.KindConflict, "session_already_finished")
	}

	feedback, err := c.LLM.EvaluateReport(ctx, report, session.Crime)
	if err != nil {
		return err
	}

	breakdown := scoreReport(report, session.Crime)

	result := &domain.GameResult{
		PlayerReport:      report,
		Breakdown:         *breakdown,
		NarrativeFeedback: feedback.NarrativeFeedback,
		BreakdownDetails:  feedback.BreakdownDetails,
		MissedFacts:       feedback.MissedFacts,
	}
	session.Finish(result)

	return c.Sessions.Update(ctx, session)
}

func scoreReport(report domain.FinalReport, truth domain.Crime) *domain.ScoreBreakdown {
	whoCorrect := len(strings.TrimSpace(report.Who)) > 3

	motiveLower := strings.ToLower(truth.Motive)
	whyCorrect := false
	for _, w := range strings.FieldsFunc(motiveLower, func(r rune) bool { return r == ' ' || r == ',' }) {
		if len(w) > 3 && strings.Contains(strings.ToLower(report.Why), w) {
			whyCorrect = true
			break
		}
	}

	howLower := strings.ToLower(truth.Method)
	howWords := strings.Fields(howLower)
	howCorrect := false
	for _, w := range howWords {
		if len(w) > 3 && strings.Contains(strings.ToLower(report.How), w) {
			howCorrect = true
			break
		}
	}

	whenCorrect := strings.Contains(report.When, truth.TimeOfCrime)

	evidenceCorrect := len(strings.TrimSpace(report.Evidence)) > 10

	return &domain.ScoreBreakdown{
		WhoCorrect:      whoCorrect,
		WhyCorrect:      whyCorrect,
		HowCorrect:      howCorrect,
		WhenCorrect:     whenCorrect,
		EvidenceCorrect: evidenceCorrect,
	}
}
