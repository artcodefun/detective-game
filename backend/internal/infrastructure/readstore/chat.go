package readstore

import (
	"context"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ChatReadRepo struct {
	coll *mongo.Collection
}

func NewChatReadRepo(db *mongo.Database) *ChatReadRepo {
	return &ChatReadRepo{coll: db.Collection("chat_messages")}
}

func (r *ChatReadRepo) GetChatMessage(ctx context.Context, messageID uuid.UUID) (*readmodels.ChatMessage, error) {
	var msg readmodels.ChatMessage
	err := r.coll.FindOne(ctx, bson.M{"_id": messageID}).Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("find chat message: %w", err)
	}
	return &msg, nil
}

func (r *ChatReadRepo) ListChatByInterrogation(ctx context.Context, interrogationID uuid.UUID) ([]*readmodels.ChatMessage, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"interrogation_id": interrogationID}, options.Find().SetSort(bson.M{"timestamp": 1}))
	if err != nil {
		return nil, fmt.Errorf("list chat messages: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*readmodels.ChatMessage
	for cursor.Next(ctx) {
		var m readmodels.ChatMessage
		if err := cursor.Decode(&m); err != nil {
			return nil, fmt.Errorf("decode chat message: %w", err)
		}
		items = append(items, &m)
	}
	if items == nil {
		items = make([]*readmodels.ChatMessage, 0)
	}
	return items, nil
}
