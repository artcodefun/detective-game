package commands

import (
	"context"
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type ScenarioCommands struct {
	Sessions   ports.SessionRepository
	LLM        ports.LlmService
	Characters ports.CharacterRepository
	Evidence   ports.EvidenceRepository
	Chronology ports.ChronologyRepository
}

func NewScenarioCommands(sessions ports.SessionRepository, llm ports.LlmService, chars ports.CharacterRepository, ev ports.EvidenceRepository, chronology ports.ChronologyRepository) *ScenarioCommands {
	return &ScenarioCommands{Sessions: sessions, LLM: llm, Characters: chars, Evidence: ev, Chronology: chronology}
}

func (c *ScenarioCommands) CreateSession(ctx context.Context, actor application.Actor) (uuid.UUID, error) {
	userID := actor.UserID
	output, err := c.LLM.GenerateScenario(ctx, actor.Locale)
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	if err := c.Sessions.FinishActiveByUserID(ctx, userID); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	sessionID := uuid.New()

	session := &domain.Session{
		ID:            sessionID,
		UserID:        userID,
		Crime:         output.Crime,
		Timeline:      output.Timeline,
		CaseName:      output.CaseName,
		CaseBrief:     output.CaseBrief,
		ActionPoints:  domain.MaxActionPoints,
		Phase:         domain.GamePhaseInvestigating,
		CreatedAt:     time.Now(),
		ContentLocale: actor.Locale,
	}

	if err := c.Sessions.Create(ctx, session); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	for i := range output.Characters {
		output.Characters[i].AssignToSession(sessionID)
		if err := c.Characters.CreateCharacter(ctx, &output.Characters[i]); err != nil {
			return uuid.Nil, application.WrapError(err)
		}
	}

	for i := range output.Evidence {
		if err := c.Evidence.AppendEvidence(ctx, sessionID, &output.Evidence[i]); err != nil {
			return uuid.Nil, application.WrapError(err)
		}
	}

	chronology := domain.NewChronologyEntry(domain.ChronologyEventTypeCaseStarted, &sessionID, domain.T("chronology.case_started"), session.CreatedAt)
	if err := c.Chronology.AppendChronologyEntry(ctx, sessionID, chronology); err != nil {
		return uuid.Nil, application.WrapError(err)
	}

	return sessionID, nil
}
