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
	Sessions ports.SessionRepository
	Reports  ports.ActionReportRepository
	LLM      ports.LlmService
}

func NewActionCommands(sessions ports.SessionRepository, reports ports.ActionReportRepository, llm ports.LlmService) *ActionCommands {
	return &ActionCommands{Sessions: sessions, Reports: reports, LLM: llm}
}

type actionRequest struct {
	name        string
	evidenceID  *uuid.UUID
	characterID *uuid.UUID
	alibiText   *string
}

func (c *ActionCommands) executeAction(ctx context.Context, actor application.Actor, req actionRequest) (uuid.UUID, error) {
	session, err := c.Sessions.FindByID(ctx, actor.SessionID)
	if err != nil {
		return uuid.Nil, err
	}

	cost := actionCost(req.name)
	if !session.SpendActionPoints(cost) {
		return uuid.Nil, application.NewAppError(application.KindConflict, "not_enough_action_points")
	}

	if err := c.Sessions.Update(ctx, session); err != nil {
		return uuid.Nil, err
	}

	body, err := c.LLM.RunAction(ctx, req.name, req.evidenceID, req.characterID, req.alibiText)
	if err != nil {
		return uuid.Nil, err
	}

	title := actionTitle(req.name)
	desc := body
	if len(desc) > 80 {
		desc = desc[:80] + "..."
	}

	report := domain.ActionReport{
		ID:          uuid.New(),
		Title:       title,
		Description: desc,
		Body:        body,
		EvidenceID:  req.evidenceID,
		CharacterID: req.characterID,
		Timestamp:   time.Now(),
	}

	if err := c.Reports.AppendReport(ctx, actor.SessionID, &report); err != nil {
		return uuid.Nil, err
	}

	return report.ID, nil
}

func (c *ActionCommands) DNAAnalysis(ctx context.Context, actor application.Actor, evidenceID uuid.UUID) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{name: "dna_analysis", evidenceID: &evidenceID})
}

func (c *ActionCommands) FingerprintsCheck(ctx context.Context, actor application.Actor, evidenceID uuid.UUID) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{name: "fingerprints", evidenceID: &evidenceID})
}

func (c *ActionCommands) AlibiCheck(ctx context.Context, actor application.Actor, characterID uuid.UUID, alibiText string) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{name: "alibi_check", characterID: &characterID, alibiText: &alibiText})
}

func (c *ActionCommands) CameraReview(ctx context.Context, actor application.Actor) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{name: "camera_review"})
}

func (c *ActionCommands) CallHistory(ctx context.Context, actor application.Actor, characterID uuid.UUID) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{name: "call_history", characterID: &characterID})
}

func (c *ActionCommands) TransactionCheck(ctx context.Context, actor application.Actor, characterID uuid.UUID) (uuid.UUID, error) {
	return c.executeAction(ctx, actor, actionRequest{name: "transaction_check", characterID: &characterID})
}

func actionCost(name string) int {
	switch name {
	case "dna_analysis", "fingerprints", "alibi_check":
		return 1
	default:
		return 2
	}
}

func actionTitle(name string) string {
	switch name {
	case "dna_analysis":
		return "Анализ ДНК"
	case "fingerprints":
		return "Отпечатки пальцев"
	case "alibi_check":
		return "Проверка алиби"
	case "camera_review":
		return "Записи с камер"
	case "call_history":
		return "История звонков"
	case "transaction_check":
		return "Банковские операции"
	}
	return name
}
