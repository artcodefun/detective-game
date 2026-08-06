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

type SessionReadRepo struct {
	coll *mongo.Collection
}

func NewSessionReadRepo(db *mongo.Database) *SessionReadRepo {
	return &SessionReadRepo{coll: db.Collection("sessions")}
}

func (r *SessionReadRepo) GetSession(ctx context.Context, sessionID uuid.UUID) (*readmodels.Session, error) {
	var session domain.Session
	err := r.coll.FindOne(ctx, bson.M{"_id": sessionID}).Decode(&session)
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}
	return readmodels.SessionFromDomain(&session), nil
}

func (r *SessionReadRepo) GetGameResult(ctx context.Context, sessionID uuid.UUID) (*readmodels.GameResult, error) {
	var session domain.Session
	err := r.coll.FindOne(ctx, bson.M{"_id": sessionID}).Decode(&session)
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}
	if session.GameResult == nil {
		return nil, fmt.Errorf("session %s has no result yet", sessionID)
	}
	return readmodels.GameResultFromDomain(session.GameResult), nil
}

func (r *SessionReadRepo) ListHistory(ctx context.Context, userID uuid.UUID) ([]*readmodels.Session, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*readmodels.Session
	for cursor.Next(ctx) {
		var s domain.Session
		if err := cursor.Decode(&s); err != nil {
			return nil, fmt.Errorf("decode session: %w", err)
		}
		items = append(items, readmodels.SessionFromDomain(&s))
	}
	if items == nil {
		items = make([]*readmodels.Session, 0)
	}
	return items, nil
}
