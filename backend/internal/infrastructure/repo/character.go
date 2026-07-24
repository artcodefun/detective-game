package repo

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CharacterRepo struct {
	coll *mongo.Collection
}

func NewCharacterRepo(db *mongo.Database) *CharacterRepo {
	return &CharacterRepo{coll: db.Collection("characters")}
}

func (r *CharacterRepo) CreateCharacter(ctx context.Context, character *domain.Character) error {
	_, err := r.coll.InsertOne(ctx, character)
	return err
}

func (r *CharacterRepo) FindCharacterBySessionAndID(ctx context.Context, sessionID uuid.UUID, prototypeID int) (*domain.Character, error) {
	var c domain.Character
	err := r.coll.FindOne(ctx, bson.M{"session_id": sessionID, "prototype.id": prototypeID}).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("find character: %w", err)
	}
	return &c, nil
}

func (r *CharacterRepo) FindCharacterByID(ctx context.Context, sessionID uuid.UUID, characterID uuid.UUID) (*domain.Character, error) {
	var c domain.Character
	err := r.coll.FindOne(ctx, bson.M{"session_id": sessionID, "_id": characterID}).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("find character: %w", err)
	}
	return &c, nil
}

func (r *CharacterRepo) FindCharactersBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.Character, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, fmt.Errorf("find characters: %w", err)
	}
	defer cursor.Close(ctx)

	var chars []*domain.Character
	for cursor.Next(ctx) {
		var c domain.Character
		if err := cursor.Decode(&c); err != nil {
			return nil, fmt.Errorf("decode character: %w", err)
		}
		chars = append(chars, &c)
	}
	if chars == nil {
		chars = make([]*domain.Character, 0)
	}
	return chars, nil
}

func (r *CharacterRepo) UpdateCharacter(ctx context.Context, character *domain.Character) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": character.ID}, character)
	return err
}
