package http

import (
	"context"
	"net/http"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/google/uuid"
)

func withUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(r.Header.Get("X-User-ID"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "missing or invalid X-User-ID header")
			return
		}
		actor := application.Actor{UserID: userID, Locale: localeFromRequest(r)}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), actorContextKey, actor)))
	})
}

func withSessionContext(sessions application.SessionQueries, translator ports.Translator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := uuid.Parse(r.Header.Get("X-Session-ID"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "missing or invalid X-Session-ID header")
			return
		}
		actor := ActorFromContext(r.Context())
		locale, err := sessions.GetSessionLocale(r.Context(), actor.UserID, sessionID)
		if err != nil {
			writeAppError(w, r.Context(), translator, err)
			return
		}
		actor.SessionID = sessionID
		actor.SessionContentLocale = locale
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), actorContextKey, actor)))
	})
}
