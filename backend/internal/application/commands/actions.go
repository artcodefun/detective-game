package commands

import (
	"context"
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type ActionCommands struct {
	Sessions   ports.SessionRepository
	Reports    ports.ActionReportRepository
	Evidence   ports.EvidenceRepository
	Characters ports.CharacterRepository
	LLM        ports.LlmService
	Chronology ports.ChronologyRepository
	TxMgr      ports.TransactionManager
}

func NewActionCommands(sessions ports.SessionRepository, reports ports.ActionReportRepository, evidence ports.EvidenceRepository, characters ports.CharacterRepository, llm ports.LlmService, chronology ports.ChronologyRepository, txMgr ports.TransactionManager) *ActionCommands {
	return &ActionCommands{Sessions: sessions, Reports: reports, Evidence: evidence, Characters: characters, LLM: llm, Chronology: chronology, TxMgr: txMgr}
}

type actionRequest struct {
	kind        domain.ActionType
	evidenceID  *uuid.UUID
	characterID *uuid.UUID
	alibiText   *string
}

func (c *ActionCommands) executeAction(ctx context.Context, actor application.Actor, req actionRequest) (uuid.UUID, error) {
	session, err := c.Sessions.FindByID(ctx, actor.SessionID)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	if session.Phase == domain.GamePhaseFinished {
		return uuid.Nil, application.NewAppError(application.KindConflict, domain.T("error.session_already_finished"))
	}

	if session.ActionPoints < req.kind.Cost() {
		return uuid.Nil, application.NewAppError(application.KindConflict, domain.T("error.not_enough_action_points"))
	}

	var evidence *domain.Evidence
	if req.evidenceID != nil {
		evidence, err = c.Evidence.FindEvidenceByID(ctx, actor.SessionID, *req.evidenceID)
		if err != nil {
			return uuid.Nil, application.WrapError(err)
		}
		existingReport, err := c.Reports.FindReportByEvidenceAction(ctx, actor.SessionID, req.kind, *req.evidenceID)
		if err != nil {
			return uuid.Nil, application.WrapError(err)
		}
		if existingReport != nil {
			return uuid.Nil, application.NewAppError(application.KindConflict, domain.T("error.evidence_action_already_completed"))
		}
	}
	var character *domain.Character
	if req.characterID != nil {
		character, err = c.Characters.FindCharacterByID(ctx, actor.SessionID, *req.characterID)
		if err != nil {
			return uuid.Nil, application.WrapError(err)
		}
	}
	body, err := c.LLM.RunAction(ctx, session.ContentLocale, string(req.kind), session.Crime, session.Timeline, evidence, character, req.alibiText)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	report := domain.ActionReport{
		ID:          uuid.New(),
		Type:        req.kind,
		Title:       domain.T("action." + string(req.kind)),
		Body:        body,
		EvidenceID:  req.evidenceID,
		CharacterID: req.characterID,
		Timestamp:   time.Now(),
	}

	if err := c.TxMgr.WithTx(ctx, func(txCtx context.Context) error {
		currentSession, err := c.Sessions.FindByID(txCtx, actor.SessionID)
		if err != nil {
			return err
		}
		if currentSession.Phase == domain.GamePhaseFinished {
			return application.NewAppError(application.KindConflict, domain.T("error.session_already_finished"))
		}
		if currentSession.ActionPoints < req.kind.Cost() {
			return application.NewAppError(application.KindConflict, domain.T("error.not_enough_action_points"))
		}
		if req.evidenceID != nil {
			existingReport, err := c.Reports.FindReportByEvidenceAction(txCtx, actor.SessionID, req.kind, *req.evidenceID)
			if err != nil {
				return err
			}
			if existingReport != nil {
				return application.NewAppError(application.KindConflict, domain.T("error.evidence_action_already_completed"))
			}
		}

		currentSession.SpendActionPoints(req.kind.Cost())
		if err := c.Sessions.Update(txCtx, currentSession); err != nil {
			return err
		}
		if err := c.Reports.AppendReport(txCtx, actor.SessionID, &report); err != nil {
			return err
		}

		chronology := domain.NewChronologyEntry(domain.ChronologyEventTypeFromAction(req.kind), &report.ID, report.Title, report.Timestamp)
		return c.Chronology.AppendChronologyEntry(txCtx, actor.SessionID, chronology)
	}); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	return report.ID, nil
}

func (c *ActionCommands) DNAAnalysis(ctx context.Context, actor application.Actor, evidenceID uuid.UUID) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{kind: domain.ActionTypeDNAAnalysis, evidenceID: &evidenceID})
}

func (c *ActionCommands) FingerprintsCheck(ctx context.Context, actor application.Actor, evidenceID uuid.UUID) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{kind: domain.ActionTypeFingerprints, evidenceID: &evidenceID})
}

func (c *ActionCommands) AlibiCheck(ctx context.Context, actor application.Actor, characterID uuid.UUID, alibiText string) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{kind: domain.ActionTypeAlibiCheck, characterID: &characterID, alibiText: &alibiText})
}

func (c *ActionCommands) CameraReview(ctx context.Context, actor application.Actor) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{kind: domain.ActionTypeCameraReview})
}

func (c *ActionCommands) CallHistory(ctx context.Context, actor application.Actor, characterID uuid.UUID) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{kind: domain.ActionTypeCallHistory, characterID: &characterID})
}

func (c *ActionCommands) TransactionCheck(ctx context.Context, actor application.Actor, characterID uuid.UUID) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{kind: domain.ActionTypeTransactionCheck, characterID: &characterID})
}
