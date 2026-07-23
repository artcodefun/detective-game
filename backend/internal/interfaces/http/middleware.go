package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/artcodefun/detective-game/backend/internal/application"
)

type contextKey string

const sessionKey contextKey = "actor"

func ActorFromContext(ctx context.Context) application.Actor {
	actor, _ := ctx.Value(sessionKey).(application.Actor)
	return actor
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeAppError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(application.AppError); ok {
		switch appErr.Kind {
		case application.KindNotFound:
			writeError(w, http.StatusNotFound, appErr.Code)
		case application.KindConflict:
			writeError(w, http.StatusConflict, appErr.Code)
		case application.KindForbidden:
			writeError(w, http.StatusForbidden, appErr.Code)
		case application.KindInvalidInput:
			writeError(w, http.StatusBadRequest, appErr.Code)
		default:
			writeError(w, http.StatusInternalServerError, "internal_error")
		}
		return
	}
	writeError(w, http.StatusInternalServerError, "internal_error")
}
