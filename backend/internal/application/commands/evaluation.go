package commands

import (
	"context"

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

	result := &domain.GameResult{
		PlayerReport:      report,
		Breakdown:         feedback.Breakdown,
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
