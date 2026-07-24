package readstore

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CharacterReadRepo struct {
	coll  *mongo.Collection
	inter *mongo.Collection
}

func NewCharacterReadRepo(db *mongo.Database) *CharacterReadRepo {
	return &CharacterReadRepo{
		coll:  db.Collection("characters"),
		inter: db.Collection("interrogations"),
	}
}

func (r *CharacterReadRepo) ListCharacters(ctx context.Context, sessionID uuid.UUID) ([]*readmodels.Character, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, fmt.Errorf("list characters: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*readmodels.Character
	for cursor.Next(ctx) {
		var c readmodels.Character
		if err := cursor.Decode(&c); err != nil {
			return nil, fmt.Errorf("decode character: %w", err)
		}
		items = append(items, &c)
	}
	if items == nil {
		items = make([]*readmodels.Character, 0)
	}
	return items, nil
}

func (r *CharacterReadRepo) GetCharacter(ctx context.Context, sessionID uuid.UUID, characterID uuid.UUID) (*readmodels.Character, error) {
	var c readmodels.Character
	err := r.coll.FindOne(ctx, bson.M{"session_id": sessionID, "_id": characterID}).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("find character: %w", err)
	}
	return &c, nil
}

func (r *CharacterReadRepo) GetInterrogation(ctx context.Context, interrogationID uuid.UUID) (*readmodels.Interrogation, error) {
	var inter readmodels.Interrogation
	err := r.inter.FindOne(ctx, bson.M{"_id": interrogationID}).Decode(&inter)
	if err != nil {
		return nil, fmt.Errorf("find interrogation: %w", err)
	}
	return &inter, nil
}
