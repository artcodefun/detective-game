package repo

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ChronologyRepo struct {
	coll *mongo.Collection
}

func NewChronologyRepo(db *mongo.Database) *ChronologyRepo {
	return &ChronologyRepo{coll: db.Collection("chronology")}
}

func (r *ChronologyRepo) AppendChronologyEntry(ctx context.Context, sessionID uuid.UUID, entry *domain.ChronologyEntry) error {
	entry.SessionID = sessionID
	_, err := r.coll.InsertOne(ctx, entry)
	return err
}

func (r *ChronologyRepo) FindChronologyBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.ChronologyEntry, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, fmt.Errorf("find chronology: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*domain.ChronologyEntry
	for cursor.Next(ctx) {
		var c domain.ChronologyEntry
		if err := cursor.Decode(&c); err != nil {
			return nil, fmt.Errorf("decode chronology entry: %w", err)
		}
		items = append(items, &c)
	}
	if items == nil {
		items = make([]*domain.ChronologyEntry, 0)
	}
	return items, nil
}

func (r *ChronologyRepo) UpdateChronologyEntry(ctx context.Context, sessionID uuid.UUID, chronologyID uuid.UUID, entryID uuid.UUID, tags []domain.NoteTag, note *string) error {
	filter := bson.M{"session_id": sessionID, "_id": chronologyID, "details.id": entryID}
	update := bson.M{
		"$set": bson.M{
			"details.$.user_tags": tags,
			"details.$.user_note": note,
		},
	}
	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update chronology entry: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("entry %s not found in chronology %s", entryID, chronologyID)
	}
	return nil
}
