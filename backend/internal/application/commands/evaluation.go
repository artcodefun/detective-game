package commands

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
)

type EvaluationCommands struct {
	Sessions       ports.SessionRepository
	Interrogations ports.InterrogationRepository
	LLM            ports.LlmService
	Chronology     ports.ChronologyRepository
	TxMgr          ports.TransactionManager
}

func NewEvaluationCommands(sessions ports.SessionRepository, interrogations ports.InterrogationRepository, llm ports.LlmService, chronology ports.ChronologyRepository, txMgr ports.TransactionManager) *EvaluationCommands {
	return &EvaluationCommands{Sessions: sessions, Interrogations: interrogations, LLM: llm, Chronology: chronology, TxMgr: txMgr}
}

func (c *EvaluationCommands) SubmitReport(ctx context.Context, actor application.Actor, report domain.FinalReport) error {
	session, err := requireActiveSession(ctx, c.Sessions, actor.SessionID)
	if err != nil {
		return application.WrapError(err)
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
	if err := c.TxMgr.WithTx(ctx, func(txCtx context.Context) error {
		currentSession, err := requireActiveSession(txCtx, c.Sessions, actor.SessionID)
		if err != nil {
			return err
		}

		activeInterrogation, err := c.Interrogations.FindActiveBySession(txCtx, actor.SessionID)
		if err != nil {
			return err
		}
		if activeInterrogation != nil {
			activeInterrogation.Complete()
			if err := c.Interrogations.UpdateInterrogation(txCtx, activeInterrogation); err != nil {
				return err
			}
		}
		currentSession.Finish(result)
		if err := c.Sessions.Update(txCtx, currentSession); err != nil {
			return err
		}

		chronology := domain.NewChronologyEntry(domain.ChronologyEventTypeFinalReport, &actor.SessionID, domain.T("chronology.final_report_submitted"), *currentSession.FinishedAt)
		return c.Chronology.AppendChronologyEntry(txCtx, actor.SessionID, chronology)
	}); err != nil {
		return application.WrapError(err)
	}
	return nil
}
