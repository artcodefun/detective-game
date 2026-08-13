package http

import (
	"context"

	"github.com/artcodefun/detective-game/backend/internal/application"
)

type contextKey string

const actorContextKey contextKey = "actor"

func ActorFromContext(ctx context.Context) application.Actor {
	actor, _ := ctx.Value(actorContextKey).(application.Actor)
	return actor
}
