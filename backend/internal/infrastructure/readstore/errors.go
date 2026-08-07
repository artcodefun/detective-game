package readstore

import (
	"errors"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func wrapFindError(op string, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("%s: %w", op, ports.ErrNotFound)
	}
	return fmt.Errorf("%s: %w", op, err)
}
