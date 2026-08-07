package commands

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type NotebookCommands struct {
	Chronology ports.ChronologyRepository
}

func NewNotebookCommands(chronology ports.ChronologyRepository) *NotebookCommands {
	return &NotebookCommands{Chronology: chronology}
}

func (c *NotebookCommands) UpdateNotebookEntry(ctx context.Context, actor application.Actor, chronologyID, entryID uuid.UUID, tags []domain.NoteTag, note *string) error {
	return application.WrapError(c.Chronology.UpdateChronologyEntry(ctx, actor.SessionID, chronologyID, entryID, tags, note))
}
