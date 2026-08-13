package commands

import (
	"context"
	"strings"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
)

type EvaluationCommands struct {
	Sessions   ports.SessionRepository
	LLM        ports.LlmService
	Chronology ports.ChronologyRepository
}

func NewEvaluationCommands(sessions ports.SessionRepository, llm ports.LlmService, chronology ports.ChronologyRepository) *EvaluationCommands {
	return &EvaluationCommands{Sessions: sessions, LLM: llm, Chronology: chronology}
}

func (c *EvaluationCommands) SubmitReport(ctx context.Context, actor application.Actor, report domain.FinalReport) error {
	session, err := c.Sessions.FindByID(ctx, actor.SessionID)
	if err != nil {
		return application.WrapError(err)
	}
	if session.Phase == domain.GamePhaseFinished {
		return application.NewAppError(application.KindConflict, domain.T("error.session_already_finished"))
	}

	feedback, err := c.LLM.EvaluateReport(ctx, session.ContentLocale, report, session.Crime)
	if err != nil {
		return application.WrapError(err)
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

	if err := c.Sessions.Update(ctx, session); err != nil {
		return application.WrapError(err)
	}

	chronology := domain.NewChronologyEntry(domain.ChronologyEventTypeFinalReport, &actor.SessionID, domain.T("chronology.final_report_submitted"), *session.FinishedAt)
	err = c.Chronology.AppendChronologyEntry(ctx, actor.SessionID, chronology)
	if err != nil {
		return application.WrapError(err)
	}
	return nil
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
