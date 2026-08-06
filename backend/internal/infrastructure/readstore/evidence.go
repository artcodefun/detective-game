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

type EvidenceReadRepo struct {
	coll *mongo.Collection
}

func NewEvidenceReadRepo(db *mongo.Database) *EvidenceReadRepo {
	return &EvidenceReadRepo{coll: db.Collection("evidence")}
}

func (r *EvidenceReadRepo) ListEvidence(ctx context.Context, sessionID uuid.UUID) ([]*readmodels.Evidence, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, fmt.Errorf("list evidence: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*readmodels.Evidence
	for cursor.Next(ctx) {
		var e domain.Evidence
		if err := cursor.Decode(&e); err != nil {
			return nil, fmt.Errorf("decode evidence: %w", err)
		}
		items = append(items, readmodels.EvidenceFromDomain(&e))
	}
	if items == nil {
		items = make([]*readmodels.Evidence, 0)
	}
	return items, nil
}

func (r *EvidenceReadRepo) GetEvidence(ctx context.Context, sessionID uuid.UUID, evidenceID uuid.UUID) (*readmodels.Evidence, error) {
	var e domain.Evidence
	err := r.coll.FindOne(ctx, bson.M{"session_id": sessionID, "_id": evidenceID}).Decode(&e)
	if err != nil {
		return nil, fmt.Errorf("find evidence: %w", err)
	}
	return readmodels.EvidenceFromDomain(&e), nil
}
