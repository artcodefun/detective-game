package repo

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepo struct {
	coll *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{coll: db.Collection("users")}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.coll.InsertOne(ctx, user)
	return err
}

func (r *UserRepo) FindUserByAccessTokenHash(ctx context.Context, accessTokenHash string) (*domain.User, error) {
	var user domain.User
	err := r.coll.FindOne(ctx, bson.M{"access_token_hash": accessTokenHash}).Decode(&user)
	if err != nil {
		return nil, wrapFindError("find user by access token hash", err)
	}
	return &user, nil
}

func (r *UserRepo) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "access_token_hash", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}
