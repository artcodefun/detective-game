package commands

import (
	"context"
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type ScenarioCommands struct {
	Users      ports.UserRepository
	Sessions   ports.SessionRepository
	LLM        ports.LlmService
	Characters ports.CharacterRepository
	Evidence   ports.EvidenceRepository
	Prototypes ports.CharacterPrototypeRepository
}

func NewScenarioCommands(users ports.UserRepository, sessions ports.SessionRepository, llm ports.LlmService, chars ports.CharacterRepository, ev ports.EvidenceRepository, p ports.CharacterPrototypeRepository) *ScenarioCommands {
	return &ScenarioCommands{Users: users, Sessions: sessions, LLM: llm, Characters: chars, Evidence: ev, Prototypes: p}
}

func (c *ScenarioCommands) CreateSession(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	_, err := c.Users.FindUserByID(ctx, userID)
	if err != nil {
		user := domain.NewUser()
		user.ID = userID
		if err := c.Users.CreateUser(ctx, &user); err != nil {
			return uuid.Nil, err
		}
	}

	prototypes, err := c.Prototypes.GetRandom(ctx, 5)
	if err != nil {
		return uuid.Nil, err
	}

	output, err := c.LLM.GenerateScenario(ctx, prototypes)
	if err != nil {
		return uuid.Nil, err
	}

	sessionID := uuid.New()

	session := &domain.Session{
		ID:           sessionID,
		UserID:       userID,
		Crime:        output.Crime,
		Timeline:     output.Timeline,
		CaseName:     output.CaseName,
		CaseBrief:    output.CaseBrief,
		ActionPoints: domain.MaxActionPoints,
		Phase:        domain.GamePhaseInvestigating,
		CreatedAt:    time.Now(),
	}

	if err := c.Sessions.Create(ctx, session); err != nil {
		return uuid.Nil, err
	}

	for i := range output.Characters {
		output.Characters[i].AssignToSession(sessionID)
		if err := c.Characters.CreateCharacter(ctx, &output.Characters[i]); err != nil {
			return uuid.Nil, err
		}
	}

	for i := range output.Evidence {
		if err := c.Evidence.AppendEvidence(ctx, sessionID, &output.Evidence[i]); err != nil {
			return uuid.Nil, err
		}
	}

	return sessionID, nil
}
