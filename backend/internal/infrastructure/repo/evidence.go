package repo

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EvidenceRepo struct {
	coll *mongo.Collection
}

func NewEvidenceRepo(db *mongo.Database) *EvidenceRepo {
	return &EvidenceRepo{coll: db.Collection("evidence")}
}

func (r *EvidenceRepo) AppendEvidence(ctx context.Context, sessionID uuid.UUID, evidence *domain.Evidence) error {
	evidence.SessionID = sessionID
	_, err := r.coll.InsertOne(ctx, evidence)
	return err
}

func (r *EvidenceRepo) FindEvidenceBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.Evidence, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, fmt.Errorf("find evidence: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*domain.Evidence
	for cursor.Next(ctx) {
		var e domain.Evidence
		if err := cursor.Decode(&e); err != nil {
			return nil, fmt.Errorf("decode evidence: %w", err)
		}
		items = append(items, &e)
	}
	if items == nil {
		items = make([]*domain.Evidence, 0)
	}
	return items, nil
}

func (r *EvidenceRepo) FindEvidenceByID(ctx context.Context, sessionID uuid.UUID, evidenceID uuid.UUID) (*domain.Evidence, error) {
	var e domain.Evidence
	err := r.coll.FindOne(ctx, bson.M{"session_id": sessionID, "_id": evidenceID}).Decode(&e)
	if err != nil {
		return nil, wrapFindError("find evidence", err)
	}
	return &e, nil
}
