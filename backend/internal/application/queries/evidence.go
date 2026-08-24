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
	e, err := q.Evidence.ListEvidence(ctx, actor.SessionID)
	return e, application.WrapError(err)
}

func (q *EvidenceQueries) GetEvidence(ctx context.Context, actor application.Actor, evidenceID uuid.UUID) (*readmodels.Evidence, error) {
	e, err := q.Evidence.GetEvidence(ctx, actor.SessionID, evidenceID)
	return e, application.WrapError(err)
}

func (q *EvidenceQueries) GetReport(ctx context.Context, actor application.Actor, reportID uuid.UUID) (*readmodels.ActionReport, error) {
	r, err := q.Reports.GetReport(ctx, actor.SessionID, reportID)
	return r, application.WrapError(err)
}

func (q *EvidenceQueries) ListReports(ctx context.Context, actor application.Actor) ([]*readmodels.ActionReport, error) {
	r, err := q.Reports.ListReports(ctx, actor.SessionID)
	return r, application.WrapError(err)
}
