package repo

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (r *ChronologyRepo) AppendNotebookEntries(ctx context.Context, sessionID, originID uuid.UUID, entries []domain.NotebookEntry) error {
	if len(entries) == 0 {
		return nil
	}
	result, err := r.coll.UpdateOne(
		ctx,
		bson.M{"session_id": sessionID, "origin_id": originID},
		bson.M{"$push": bson.M{"details": bson.M{"$each": entries}}},
	)
	if err != nil {
		return fmt.Errorf("append notebook entries: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("append notebook entries: %w", ports.ErrNotFound)
	}
	return nil
}

func (r *ChronologyRepo) FindChronologyBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.ChronologyEntry, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"session_id": sessionID}, options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}))
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
	filter := bson.M{"session_id": sessionID, "id": chronologyID, "details.id": entryID}
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
		return fmt.Errorf("update chronology entry: %w", ports.ErrNotFound)
	}
	return nil
}
