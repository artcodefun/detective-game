package readstore

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ReportReadRepo struct {
	coll *mongo.Collection
}

func NewReportReadRepo(db *mongo.Database) *ReportReadRepo {
	return &ReportReadRepo{coll: db.Collection("reports")}
}

func (r *ReportReadRepo) GetReport(ctx context.Context, sessionID, reportID uuid.UUID) (*readmodels.ActionReport, error) {
	var report domain.ActionReport
	err := r.coll.FindOne(ctx, bson.M{"_id": reportID, "session_id": sessionID}).Decode(&report)
	if err != nil {
		return nil, wrapFindError("find report", err)
	}
	return readmodels.ActionReportFromDomain(&report), nil
}

func (r *ReportReadRepo) ListReports(ctx context.Context, sessionID uuid.UUID) ([]*readmodels.ActionReport, error) {
	cursor, err := r.coll.Find(
		ctx,
		bson.M{"session_id": sessionID},
		options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}),
	)
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
		item := readmodels.ActionReportFromDomain(&report)
		items = append(items, item)
	}
	if items == nil {
		items = make([]*readmodels.ActionReport, 0)
	}
	return items, nil
}
