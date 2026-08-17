package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/google/uuid"
)

func withAuthentication(users application.UserQueries, translator ports.Translator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := localeFromRequest(r)
		actorContext := context.WithValue(r.Context(), actorContextKey, application.Actor{Locale: locale})
		accessToken, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			writeError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
			return
		}
		userID, err := users.Authenticate(r.Context(), accessToken)
		if err != nil {
			writeAppError(w, actorContext, translator, err)
			return
		}
		actor := application.Actor{UserID: userID, Locale: locale}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), actorContextKey, actor)))
	})
}

func bearerToken(value string) (string, bool) {
	scheme, token, ok := strings.Cut(strings.TrimSpace(value), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
		return "", false
	}
	return strings.TrimSpace(token), true
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
