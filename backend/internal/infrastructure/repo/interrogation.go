package repo

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type InterrogationRepo struct {
	coll *mongo.Collection
}

func (r *InterrogationRepo) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "session_id", Value: 1}},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"phase": domain.InterrogationActive}),
	})
	return err
}

func NewInterrogationRepo(db *mongo.Database) *InterrogationRepo {
	return &InterrogationRepo{coll: db.Collection("interrogations")}
}

func (r *InterrogationRepo) CreateInterrogation(ctx context.Context, interrogation *domain.Interrogation) error {
	_, err := r.coll.InsertOne(ctx, interrogation)
	return err
}

func (r *InterrogationRepo) FindInterrogationByID(ctx context.Context, sessionID, id uuid.UUID) (*domain.Interrogation, error) {
	var inter domain.Interrogation
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "session_id": sessionID}).Decode(&inter)
	if err != nil {
		return nil, wrapFindError("find interrogation", err)
	}
	return &inter, nil
}

func (r *InterrogationRepo) FindActiveBySession(ctx context.Context, sessionID uuid.UUID) (*domain.Interrogation, error) {
	var inter domain.Interrogation
	err := r.coll.FindOne(ctx, bson.M{"session_id": sessionID, "phase": domain.InterrogationActive}).Decode(&inter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find active interrogation: %w", err)
	}
	return &inter, nil
}

func (r *InterrogationRepo) UpdateInterrogation(ctx context.Context, interrogation *domain.Interrogation) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": interrogation.ID}, interrogation)
	return err
}
