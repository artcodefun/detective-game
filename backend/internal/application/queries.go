package application

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type SessionQueries interface {
	GetSession(ctx context.Context, actor Actor) (*readmodels.Session, error)
	ListHistory(ctx context.Context) ([]*readmodels.Session, error)
	GetGameResult(ctx context.Context, actor Actor) (*readmodels.GameResult, error)
}

type CharacterQueries interface {
	ListCharacters(ctx context.Context, actor Actor) ([]*readmodels.Character, error)
	GetCharacter(ctx context.Context, actor Actor, characterID uuid.UUID) (*readmodels.Character, error)
}

type EvidenceQueries interface {
	ListEvidence(ctx context.Context, actor Actor) ([]*readmodels.Evidence, error)
	GetEvidence(ctx context.Context, actor Actor, evidenceID uuid.UUID) (*readmodels.Evidence, error)
	GetReport(ctx context.Context, actor Actor, reportID uuid.UUID) (*readmodels.ActionReport, error)
	ListReports(ctx context.Context, actor Actor) ([]*readmodels.ActionReport, error)
}

type ChronologyQueries interface {
	GetChronology(ctx context.Context, actor Actor) ([]*readmodels.ChronologyEntry, error)
}

type ChatQueries interface {
	GetChatMessage(ctx context.Context, actor Actor, messageID uuid.UUID) (*readmodels.ChatMessage, error)
	ListChatByInterrogation(ctx context.Context, actor Actor, interrogationID uuid.UUID) ([]*readmodels.ChatMessage, error)
	GetInterrogation(ctx context.Context, actor Actor, interrogationID uuid.UUID) (*readmodels.Interrogation, error)
}
