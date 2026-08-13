package dtos

import (
	"time"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type ActionReport struct {
	ID          uuid.UUID  `json:"id"`
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Body        string     `json:"body"`
	EvidenceID  *uuid.UUID `json:"evidence_id,omitempty"`
	CharacterID *uuid.UUID `json:"character_id,omitempty"`
	Timestamp   time.Time  `json:"timestamp"`
}

func ActionReportFromReadModel(report *readmodels.ActionReport, translator ports.Translator) ActionReport {
	return ActionReport{ID: report.ID, Type: report.Type, Title: translator.Translate(report.ContentLocale, report.Title), Description: report.Description, Body: report.Body, EvidenceID: report.EvidenceID, CharacterID: report.CharacterID, Timestamp: report.Timestamp}
}
