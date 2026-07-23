package application

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type ScenarioCommands interface {
	CreateSession(ctx context.Context) (uuid.UUID, error)
}

type InterrogationCommands interface {
	Create(ctx context.Context, actor Actor, characterID uuid.UUID) (*domain.Interrogation, error)
	AddMessage(ctx context.Context, actor Actor, interrogationID uuid.UUID, message string) (messageID uuid.UUID, err error)
	Complete(ctx context.Context, actor Actor, interrogationID uuid.UUID) error
}

type EvaluationCommands interface {
	SubmitReport(ctx context.Context, actor Actor, report domain.FinalReport) error
}

type ActionCommands interface {
	DNAAnalysis(ctx context.Context, actor Actor, evidenceID uuid.UUID) (reportID uuid.UUID, err error)
	FingerprintsCheck(ctx context.Context, actor Actor, evidenceID uuid.UUID) (reportID uuid.UUID, err error)
	AlibiCheck(ctx context.Context, actor Actor, characterID uuid.UUID, alibiText string) (reportID uuid.UUID, err error)
	CameraReview(ctx context.Context, actor Actor) (reportID uuid.UUID, err error)
	CallHistory(ctx context.Context, actor Actor, characterID uuid.UUID) (reportID uuid.UUID, err error)
	TransactionCheck(ctx context.Context, actor Actor, characterID uuid.UUID) (reportID uuid.UUID, err error)
}

type NotebookCommands interface {
	UpdateNotebookEntry(ctx context.Context, actor Actor, chronologyID, entryID uuid.UUID, tags []domain.NoteTag, note *string) error
}
