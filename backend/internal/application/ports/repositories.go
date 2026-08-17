package ports

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	FindUserByAccessTokenHash(ctx context.Context, accessTokenHash string) (*domain.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
	FinishActiveByUserID(ctx context.Context, userID uuid.UUID) error
}

type CharacterRepository interface {
	CreateCharacter(ctx context.Context, character *domain.Character) error
	FindCharacterByID(ctx context.Context, sessionID uuid.UUID, characterID uuid.UUID) (*domain.Character, error)
	FindCharactersBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.Character, error)
	UpdateCharacter(ctx context.Context, character *domain.Character) error
}

type InterrogationRepository interface {
	CreateInterrogation(ctx context.Context, interrogation *domain.Interrogation) error
	FindInterrogationByID(ctx context.Context, id uuid.UUID) (*domain.Interrogation, error)
	FindActiveBySession(ctx context.Context, sessionID uuid.UUID) (*domain.Interrogation, error)
	UpdateInterrogation(ctx context.Context, interrogation *domain.Interrogation) error
}

type ChatMessageRepository interface {
	AppendChatMessage(ctx context.Context, msg *domain.ChatMessage) error
	FindChatByInterrogation(ctx context.Context, interrogationID uuid.UUID) ([]*domain.ChatMessage, error)
	FindChatByID(ctx context.Context, id uuid.UUID) (*domain.ChatMessage, error)
}

type EvidenceRepository interface {
	AppendEvidence(ctx context.Context, sessionID uuid.UUID, evidence *domain.Evidence) error
	FindEvidenceBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.Evidence, error)
	FindEvidenceByID(ctx context.Context, sessionID uuid.UUID, evidenceID uuid.UUID) (*domain.Evidence, error)
}

type ActionReportRepository interface {
	AppendReport(ctx context.Context, sessionID uuid.UUID, report *domain.ActionReport) error
	FindReportByID(ctx context.Context, reportID uuid.UUID) (*domain.ActionReport, error)
	FindReportsBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.ActionReport, error)
}

type ChronologyRepository interface {
	AppendChronologyEntry(ctx context.Context, sessionID uuid.UUID, entry *domain.ChronologyEntry) error
	AppendNotebookEntries(ctx context.Context, sessionID, originID uuid.UUID, entries []domain.NotebookEntry) error
	FindChronologyBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.ChronologyEntry, error)
	UpdateChronologyEntry(ctx context.Context, sessionID uuid.UUID, chronologyID uuid.UUID, entryID uuid.UUID, tags []domain.NoteTag, note *string) error
}
