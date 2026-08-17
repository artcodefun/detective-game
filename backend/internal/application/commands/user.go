package commands

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
)

const accessTokenBytes = 32

type UserCommands struct {
	Users ports.UserRepository
}

func NewUserCommands(users ports.UserRepository) *UserCommands {
	return &UserCommands{Users: users}
}

func (c *UserCommands) RegisterAnonymous(ctx context.Context, device domain.DeviceInfo) (string, error) {
	if !device.IsValid() {
		return "", application.ErrInvalidInput
	}

	rawToken := make([]byte, accessTokenBytes)
	if _, err := rand.Read(rawToken); err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}
	accessToken := base64.RawURLEncoding.EncodeToString(rawToken)
	tokenHash := sha256.Sum256([]byte(accessToken))

	user := domain.NewUser(base64.RawURLEncoding.EncodeToString(tokenHash[:]), device)
	if err := c.Users.CreateUser(ctx, &user); err != nil {
		return "", application.WrapError(err)
	}

	return accessToken, nil
}
