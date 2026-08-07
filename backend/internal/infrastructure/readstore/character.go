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
		var c domain.Character
		if err := cursor.Decode(&c); err != nil {
			return nil, fmt.Errorf("decode character: %w", err)
		}
		rm := readmodels.CharacterFromDomain(c)
		items = append(items, &rm)
	}
	if items == nil {
		items = make([]*readmodels.Character, 0)
	}
	return items, nil
}

func (r *CharacterReadRepo) GetCharacter(ctx context.Context, sessionID uuid.UUID, characterID uuid.UUID) (*readmodels.Character, error) {
	var c domain.Character
	err := r.coll.FindOne(ctx, bson.M{"session_id": sessionID, "_id": characterID}).Decode(&c)
	if err != nil {
		return nil, wrapFindError("find character", err)
	}
	rm := readmodels.CharacterFromDomain(c)
	return &rm, nil
}

func (r *CharacterReadRepo) GetInterrogation(ctx context.Context, interrogationID uuid.UUID) (*readmodels.Interrogation, error) {
	var inter domain.Interrogation
	err := r.inter.FindOne(ctx, bson.M{"_id": interrogationID}).Decode(&inter)
	if err != nil {
		return nil, wrapFindError("find interrogation", err)
	}
	return readmodels.InterrogationFromDomain(&inter), nil
}

func (r *CharacterReadRepo) GetActiveInterrogation(ctx context.Context, sessionID uuid.UUID) (*readmodels.Interrogation, error) {
	var inter domain.Interrogation
	err := r.inter.FindOne(ctx, bson.M{"session_id": sessionID, "phase": domain.InterrogationActive}).Decode(&inter)
	if err != nil {
		return nil, wrapFindError("find active interrogation", err)
	}
	return readmodels.InterrogationFromDomain(&inter), nil
}
