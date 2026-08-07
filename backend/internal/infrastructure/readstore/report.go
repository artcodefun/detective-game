package readstore

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReportReadRepo struct {
	coll *mongo.Collection
}

func NewReportReadRepo(db *mongo.Database) *ReportReadRepo {
	return &ReportReadRepo{coll: db.Collection("reports")}
}

func (r *ReportReadRepo) GetReport(ctx context.Context, reportID uuid.UUID) (*readmodels.ActionReport, error) {
	var report domain.ActionReport
	err := r.coll.FindOne(ctx, bson.M{"_id": reportID}).Decode(&report)
	if err != nil {
		return nil, wrapFindError("find report", err)
	}
	return readmodels.ActionReportFromDomain(&report), nil
}

func (r *ReportReadRepo) ListReports(ctx context.Context, sessionID uuid.UUID) ([]*readmodels.ActionReport, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*readmodels.ActionReport
	for cursor.Next(ctx) {
		var report domain.ActionReport
		if err := cursor.Decode(&report); err != nil {
			return nil, fmt.Errorf("decode report: %w", err)
		}
		items = append(items, readmodels.ActionReportFromDomain(&report))
	}
	if items == nil {
		items = make([]*readmodels.ActionReport, 0)
	}
	return items, nil
}
