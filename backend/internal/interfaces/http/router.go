package http

import (
	"context"
	"net/http"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/google/uuid"
)

func NewRouter(h *Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/sessions", h.CreateSession)
	mux.HandleFunc("GET /api/v1/sessions/history", h.ListHistory)
	mux.HandleFunc("GET /api/v1/sessions/current", h.GetSession)

	mux.HandleFunc("GET /api/v1/evidence", h.ListEvidence)
	mux.HandleFunc("GET /api/v1/evidence/{evId}", h.GetEvidence)
	mux.HandleFunc("GET /api/v1/characters", h.ListCharacters)
	mux.HandleFunc("GET /api/v1/characters/{charId}", h.GetCharacter)
	mux.HandleFunc("GET /api/v1/chronology", h.GetChronology)
	mux.HandleFunc("POST /api/v1/interrogations", h.CreateInterrogation)
	mux.HandleFunc("GET /api/v1/interrogations/{interId}", h.GetInterrogation)
	mux.HandleFunc("POST /api/v1/interrogations/{interId}/messages", h.AddInterrogationMessage)
	mux.HandleFunc("GET /api/v1/interrogations/{interId}/messages", h.GetInterrogationMessages)
	mux.HandleFunc("PATCH /api/v1/interrogations/{interId}/complete", h.CompleteInterrogation)
	mux.HandleFunc("PATCH /api/v1/chronology/{chronId}/notes/{noteId}", h.UpdateNotebookEntry)
	mux.HandleFunc("POST /api/v1/actions/dna-analysis", h.DNAAnalysis)
	mux.HandleFunc("POST /api/v1/actions/fingerprints", h.FingerprintsCheck)
	mux.HandleFunc("POST /api/v1/actions/alibi-check", h.AlibiCheck)
	mux.HandleFunc("POST /api/v1/actions/camera-review", h.CameraReview)
	mux.HandleFunc("POST /api/v1/actions/call-history", h.CallHistory)
	mux.HandleFunc("POST /api/v1/actions/transactions", h.TransactionCheck)
	mux.HandleFunc("GET /api/v1/reports/{reportId}", h.GetReport)
	mux.HandleFunc("POST /api/v1/reports", h.SubmitReport)

	return sessionMiddleware(mux)
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
