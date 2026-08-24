package repo

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoTxManager executes callbacks in a MongoDB multi-document transaction.
// Mongo repositories receive the transaction session through txCtx.
type MongoTxManager struct {
	client *mongo.Client
}

func NewMongoTxManager(client *mongo.Client) *MongoTxManager {
	return &MongoTxManager{client: client}
}

func (m *MongoTxManager) WithTx(ctx context.Context, fn func(txCtx context.Context) error) error {
	session, err := m.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		return nil, fn(txCtx)
	})
	return err
}

var _ ports.TransactionManager = (*MongoTxManager)(nil)
