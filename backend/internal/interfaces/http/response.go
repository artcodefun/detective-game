package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeAppError(w http.ResponseWriter, ctx context.Context, translator ports.Translator, err error) {
	var appErr application.AppError
	if errors.As(err, &appErr) {
		message := translator.Translate(ActorFromContext(ctx).Locale, appErr.Translation)
		writeError(w, statusForAppError(appErr.Kind), message)
		return
	}
	writeError(w, http.StatusInternalServerError, "internal_error")
}

func statusForAppError(kind application.ErrorKind) int {
	switch kind {
	case application.KindUnauthorized:
		return http.StatusUnauthorized
	case application.KindNotFound:
		return http.StatusNotFound
	case application.KindConflict:
		return http.StatusConflict
	case application.KindForbidden:
		return http.StatusForbidden
	case application.KindInvalidInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
