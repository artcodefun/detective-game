package repo

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (r *UserRepo) FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	return &user, nil
}
