package readstore

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChronologyReadRepo struct {
	coll *mongo.Collection
}

func NewChronologyReadRepo(db *mongo.Database) *ChronologyReadRepo {
	return &ChronologyReadRepo{coll: db.Collection("chronology")}
}

func (r *ChronologyReadRepo) GetChronology(ctx context.Context, sessionID uuid.UUID) ([]*readmodels.ChronologyEntry, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, fmt.Errorf("list chronology: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*readmodels.ChronologyEntry
	for cursor.Next(ctx) {
		var c readmodels.ChronologyEntry
		if err := cursor.Decode(&c); err != nil {
			return nil, fmt.Errorf("decode chronology entry: %w", err)
		}
		items = append(items, &c)
	}
	if items == nil {
		items = make([]*readmodels.ChronologyEntry, 0)
	}
	return items, nil
}
