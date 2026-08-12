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
	LLM        ports.LlmService
	Chronology ports.ChronologyRepository
}

func NewActionCommands(sessions ports.SessionRepository, reports ports.ActionReportRepository, llm ports.LlmService, chronology ports.ChronologyRepository) *ActionCommands {
	return &ActionCommands{Sessions: sessions, Reports: reports, LLM: llm, Chronology: chronology}
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

	if !session.SpendActionPoints(req.kind.Cost()) {
		return uuid.Nil, application.NewAppError(application.KindConflict, "not_enough_action_points")
	}

	if err := c.Sessions.Update(ctx, session); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	body, err := c.LLM.RunAction(ctx, string(req.kind), req.evidenceID, req.characterID, req.alibiText)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	report := domain.ActionReport{
		ID:          uuid.New(),
		Type:        req.kind,
		Body:        body,
		EvidenceID:  req.evidenceID,
		CharacterID: req.characterID,
		Timestamp:   time.Now(),
	}

	if err := c.Reports.AppendReport(ctx, actor.SessionID, &report); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	chronology := domain.NewChronologyEntry(domain.ChronologyEventTypeFromAction(req.kind), &report.ID, actionTitle(req.kind), report.Timestamp)
	if err := c.Chronology.AppendChronologyEntry(ctx, actor.SessionID, chronology); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	return report.ID, nil
}

func actionTitle(action domain.ActionType) string {
	switch action {
	case domain.ActionTypeDNAAnalysis:
		return "Анализ ДНК"
	case domain.ActionTypeFingerprints:
		return "Проверка отпечатков"
	case domain.ActionTypeAlibiCheck:
		return "Проверка алиби"
	case domain.ActionTypeCameraReview:
		return "Просмотр записей с камер"
	case domain.ActionTypeCallHistory:
		return "История звонков"
	case domain.ActionTypeTransactionCheck:
		return "Проверка транзакций"
	default:
		return string(action)
	}
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
