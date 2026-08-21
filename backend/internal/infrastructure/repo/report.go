package repo

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReportRepo struct {
	coll *mongo.Collection
}

func NewReportRepo(db *mongo.Database) *ReportRepo {
	return &ReportRepo{coll: db.Collection("reports")}
}

func (r *ReportRepo) AppendReport(ctx context.Context, sessionID uuid.UUID, report *domain.ActionReport) error {
	report.SessionID = sessionID
	_, err := r.coll.InsertOne(ctx, report)
	return err
}

func (r *ReportRepo) FindReportByID(ctx context.Context, reportID uuid.UUID) (*domain.ActionReport, error) {
	var report domain.ActionReport
	err := r.coll.FindOne(ctx, bson.M{"_id": reportID}).Decode(&report)
	if err != nil {
		return nil, wrapFindError("find report", err)
	}
	return &report, nil
}

func (r *ReportRepo) FindReportByEvidenceAction(ctx context.Context, sessionID uuid.UUID, actionType domain.ActionType, evidenceID uuid.UUID) (*domain.ActionReport, error) {
	var report domain.ActionReport
	err := r.coll.FindOne(ctx, bson.M{
		"session_id":  sessionID,
		"type":        actionType,
		"evidence_id": evidenceID,
	}).Decode(&report)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, wrapFindError("find report by evidence action", err)
	}
	return &report, nil
}

func (r *ReportRepo) FindReportsBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.ActionReport, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, fmt.Errorf("find reports: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*domain.ActionReport
	for cursor.Next(ctx) {
		var report domain.ActionReport
		if err := cursor.Decode(&report); err != nil {
			return nil, fmt.Errorf("decode report: %w", err)
		}
		items = append(items, &report)
	}
	if items == nil {
		items = make([]*domain.ActionReport, 0)
	}
	return items, nil
}
