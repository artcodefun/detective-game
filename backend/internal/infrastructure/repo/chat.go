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

type ChatRepo struct {
	coll *mongo.Collection
}

func NewChatRepo(db *mongo.Database) *ChatRepo {
	return &ChatRepo{coll: db.Collection("chat_messages")}
}

func (r *ChatRepo) AppendChatMessage(ctx context.Context, msg *domain.ChatMessage) error {
	_, err := r.coll.InsertOne(ctx, msg)
	return err
}

func (r *ChatRepo) FindChatByInterrogation(ctx context.Context, sessionID, interrogationID uuid.UUID) ([]*domain.ChatMessage, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"session_id": sessionID, "interrogation_id": interrogationID}, options.Find().SetSort(bson.M{"timestamp": 1}))
	if err != nil {
		return nil, fmt.Errorf("find chat messages: %w", err)
	}
	defer cursor.Close(ctx)

	var msgs []*domain.ChatMessage
	for cursor.Next(ctx) {
		var m domain.ChatMessage
		if err := cursor.Decode(&m); err != nil {
			return nil, fmt.Errorf("decode chat message: %w", err)
		}
		msgs = append(msgs, &m)
	}
	if msgs == nil {
		msgs = make([]*domain.ChatMessage, 0)
	}
	return msgs, nil
}

func (r *ChatRepo) FindChatByID(ctx context.Context, id uuid.UUID) (*domain.ChatMessage, error) {
	var msg domain.ChatMessage
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&msg)
	if err != nil {
		return nil, wrapFindError("find chat message", err)
	}
	return &msg, nil
}
