package repo

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SessionRepo struct {
	coll *mongo.Collection
}

func NewSessionRepo(db *mongo.Database) *SessionRepo {
	return &SessionRepo{coll: db.Collection("sessions")}
}

func (r *SessionRepo) Create(ctx context.Context, session *domain.Session) error {
	_, err := r.coll.InsertOne(ctx, session)
	return err
}

func (r *SessionRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	var session domain.Session
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
	if err != nil {
		return nil, wrapFindError("find session", err)
	}
	return &session, nil
}

func (r *SessionRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("find sessions by user: %w", err)
	}
	defer cursor.Close(ctx)

	var sessions []*domain.Session
	for cursor.Next(ctx) {
		var s domain.Session
		if err := cursor.Decode(&s); err != nil {
			return nil, fmt.Errorf("decode session: %w", err)
		}
		sessions = append(sessions, &s)
	}
	if sessions == nil {
		sessions = make([]*domain.Session, 0)
	}
	return sessions, nil
}

func (r *SessionRepo) Update(ctx context.Context, session *domain.Session) error {
	_, err := r.coll.ReplaceOne(ctx, bson.M{"_id": session.ID}, session)
	return err
}

func (r *SessionRepo) FinishActiveByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.coll.UpdateMany(ctx,
		bson.M{"user_id": userID, "phase": bson.M{"$ne": "finished"}},
		bson.M{"$set": bson.M{"phase": "finished"}},
	)
	return err
}
