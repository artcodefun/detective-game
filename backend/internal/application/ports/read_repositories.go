package ports

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type SessionReadRepository interface {
	GetSession(ctx context.Context, userID uuid.UUID) (*readmodels.Session, error)
	GetSessionByID(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (*readmodels.Session, error)
	GetGameResult(ctx context.Context, sessionID uuid.UUID) (*readmodels.GameResult, error)
	ListHistory(ctx context.Context, userID uuid.UUID) ([]*readmodels.Session, error)
}

type CharacterReadRepository interface {
	ListCharacters(ctx context.Context, sessionID uuid.UUID) ([]*readmodels.Character, error)
	GetCharacter(ctx context.Context, sessionID uuid.UUID, characterID uuid.UUID) (*readmodels.Character, error)
	GetInterrogation(ctx context.Context, interrogationID uuid.UUID) (*readmodels.Interrogation, error)
	GetActiveInterrogation(ctx context.Context, sessionID uuid.UUID) (*readmodels.Interrogation, error)
}

type ChatMessageReadRepository interface {
	GetChatMessage(ctx context.Context, messageID uuid.UUID) (*readmodels.ChatMessage, error)
	ListChatByInterrogation(ctx context.Context, interrogationID uuid.UUID) ([]*readmodels.ChatMessage, error)
}

type EvidenceReadRepository interface {
	ListEvidence(ctx context.Context, sessionID uuid.UUID) ([]*readmodels.Evidence, error)
	GetEvidence(ctx context.Context, sessionID uuid.UUID, evidenceID uuid.UUID) (*readmodels.Evidence, error)
}

type ReportReadRepository interface {
	GetReport(ctx context.Context, reportID uuid.UUID) (*readmodels.ActionReport, error)
	ListReports(ctx context.Context, sessionID uuid.UUID) ([]*readmodels.ActionReport, error)
}

type ChronologyReadRepository interface {
	GetChronology(ctx context.Context, sessionID uuid.UUID) ([]readmodels.ChronologyEntry, error)
}
