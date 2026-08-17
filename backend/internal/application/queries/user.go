package queries

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/google/uuid"
)

type UserQueries struct {
	Users ports.UserRepository
}

func NewUserQueries(users ports.UserRepository) *UserQueries {
	return &UserQueries{Users: users}
}

func (q *UserQueries) Authenticate(ctx context.Context, accessToken string) (uuid.UUID, error) {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return uuid.Nil, application.ErrUnauthorized
	}

	tokenHash := sha256.Sum256([]byte(accessToken))
	user, err := q.Users.FindUserByAccessTokenHash(ctx, base64.RawURLEncoding.EncodeToString(tokenHash[:]))
	if errors.Is(err, ports.ErrNotFound) {
		return uuid.Nil, application.ErrUnauthorized
	}
	if err != nil {
		return uuid.Nil, application.WrapError(err)
	}
	return user.ID, nil
}
