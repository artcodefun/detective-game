package http

import (
	"context"
	"net/http"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/google/uuid"
)

func NewRouter(h *Handlers) http.Handler {
	apiMux := http.NewServeMux()

	apiMux.HandleFunc("POST /api/v1/sessions", h.CreateSession)
	apiMux.HandleFunc("GET /api/v1/sessions/history", h.ListHistory)
	apiMux.HandleFunc("GET /api/v1/sessions/current", h.GetSession)

	apiMux.HandleFunc("GET /api/v1/evidence", h.ListEvidence)
	apiMux.HandleFunc("GET /api/v1/evidence/{evId}", h.GetEvidence)
	apiMux.HandleFunc("GET /api/v1/characters", h.ListCharacters)
	apiMux.HandleFunc("GET /api/v1/characters/{charId}", h.GetCharacter)
	apiMux.HandleFunc("GET /api/v1/chronology", h.GetChronology)
	apiMux.HandleFunc("POST /api/v1/interrogations", h.CreateInterrogation)
	apiMux.HandleFunc("GET /api/v1/interrogations/{interId}", h.GetInterrogation)
	apiMux.HandleFunc("POST /api/v1/interrogations/{interId}/messages", h.AddInterrogationMessage)
	apiMux.HandleFunc("GET /api/v1/interrogations/{interId}/messages", h.GetInterrogationMessages)
	apiMux.HandleFunc("PATCH /api/v1/interrogations/{interId}/complete", h.CompleteInterrogation)
	apiMux.HandleFunc("PATCH /api/v1/chronology/{chronId}/notes/{noteId}", h.UpdateNotebookEntry)
	apiMux.HandleFunc("POST /api/v1/actions/dna-analysis", h.DNAAnalysis)
	apiMux.HandleFunc("POST /api/v1/actions/fingerprints", h.FingerprintsCheck)
	apiMux.HandleFunc("POST /api/v1/actions/alibi-check", h.AlibiCheck)
	apiMux.HandleFunc("POST /api/v1/actions/camera-review", h.CameraReview)
	apiMux.HandleFunc("POST /api/v1/actions/call-history", h.CallHistory)
	apiMux.HandleFunc("POST /api/v1/actions/transactions", h.TransactionCheck)
	apiMux.HandleFunc("GET /api/v1/reports/{reportId}", h.GetReport)
	apiMux.HandleFunc("POST /api/v1/reports", h.SubmitReport)

	docsMux := http.NewServeMux()
	docsMux.HandleFunc("GET /api/v1/openapi.yaml", openapiSpecHandler)
	docsMux.HandleFunc("GET /api/v1/docs", swaggerDocsHandler)

	topMux := http.NewServeMux()
	topMux.Handle("/api/v1/docs", docsMux)
	topMux.Handle("/api/v1/openapi.yaml", docsMux)
	topMux.Handle("/", sessionMiddleware(apiMux))

	return topMux
}

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(r.Header.Get("X-User-ID"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "missing or invalid X-User-ID header")
			return
		}

		actor := application.Actor{UserID: userID}

		if needsSessionID(r) {
			sessionID, err := uuid.Parse(r.Header.Get("X-Session-ID"))
			if err != nil {
				writeError(w, http.StatusBadRequest, "missing or invalid X-Session-ID header")
				return
			}
			actor.SessionID = sessionID
		}

		r = r.WithContext(context.WithValue(r.Context(), sessionKey, actor))
		next.ServeHTTP(w, r)
	})
}

func needsSessionID(r *http.Request) bool {
	if r.Method == http.MethodPost && r.URL.Path == "/api/v1/sessions" {
		return false
	}
	if r.Method == http.MethodGet && r.URL.Path == "/api/v1/sessions/history" {
		return false
	}
	return true
}
