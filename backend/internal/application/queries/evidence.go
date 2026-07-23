package queries

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
)

type EvidenceQueries struct {
	Evidence ports.EvidenceReadRepository
	Reports  ports.ReportReadRepository
}

func NewEvidenceQueries(ev ports.EvidenceReadRepository, reports ports.ReportReadRepository) *EvidenceQueries {
	return &EvidenceQueries{Evidence: ev, Reports: reports}
}

func (q *EvidenceQueries) ListEvidence(ctx context.Context, actor application.Actor) ([]*readmodels.Evidence, error) {
	return q.Evidence.ListEvidence(ctx, actor.SessionID)
}

func (q *EvidenceQueries) GetEvidence(ctx context.Context, actor application.Actor, evidenceID uuid.UUID) (*readmodels.Evidence, error) {
	return q.Evidence.GetEvidence(ctx, actor.SessionID, evidenceID)
}

func (q *EvidenceQueries) GetReport(ctx context.Context, actor application.Actor, reportID uuid.UUID) (*readmodels.ActionReport, error) {
	return q.Reports.GetReport(ctx, reportID)
}

func (q *EvidenceQueries) ListReports(ctx context.Context, actor application.Actor) ([]*readmodels.ActionReport, error) {
	return q.Reports.ListReports(ctx, actor.SessionID)
}
