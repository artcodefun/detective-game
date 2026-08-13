package http

import "net/http"

func NewRouter(h *Handlers) http.Handler {
	apiMux := http.NewServeMux()

	registerUserRoute(apiMux, "POST /api/v1/sessions", h.CreateSession)
	registerUserRoute(apiMux, "GET /api/v1/sessions/history", h.ListHistory)
	registerUserRoute(apiMux, "GET /api/v1/sessions/current", h.GetSession)
	registerUserRoute(apiMux, "GET /api/v1/sessions/{id}", h.GetSessionByID)

	registerSessionRoute(apiMux, "GET /api/v1/evidence", h.ListEvidence)
	registerSessionRoute(apiMux, "GET /api/v1/evidence/{evId}", h.GetEvidence)
	registerSessionRoute(apiMux, "GET /api/v1/characters", h.ListCharacters)
	registerSessionRoute(apiMux, "GET /api/v1/characters/{charId}", h.GetCharacter)
	registerSessionRoute(apiMux, "GET /api/v1/chronology", h.GetChronology)
	registerSessionRoute(apiMux, "GET /api/v1/interrogations/active", h.GetActiveInterrogation)
	registerSessionRoute(apiMux, "POST /api/v1/interrogations", h.CreateInterrogation)
	registerSessionRoute(apiMux, "GET /api/v1/interrogations/{interId}", h.GetInterrogation)
	registerSessionRoute(apiMux, "POST /api/v1/interrogations/{interId}/messages", h.AddInterrogationMessage)
	registerSessionRoute(apiMux, "GET /api/v1/interrogations/{interId}/messages", h.GetInterrogationMessages)
	registerSessionRoute(apiMux, "PATCH /api/v1/interrogations/{interId}/complete", h.CompleteInterrogation)
	registerSessionRoute(apiMux, "PATCH /api/v1/chronology/{chronId}/notes/{noteId}", h.UpdateNotebookEntry)
	registerSessionRoute(apiMux, "POST /api/v1/actions/dna-analysis", h.DNAAnalysis)
	registerSessionRoute(apiMux, "POST /api/v1/actions/fingerprints", h.FingerprintsCheck)
	registerSessionRoute(apiMux, "POST /api/v1/actions/alibi-check", h.AlibiCheck)
	registerSessionRoute(apiMux, "POST /api/v1/actions/camera-review", h.CameraReview)
	registerSessionRoute(apiMux, "POST /api/v1/actions/call-history", h.CallHistory)
	registerSessionRoute(apiMux, "POST /api/v1/actions/transactions", h.TransactionCheck)
	registerSessionRoute(apiMux, "GET /api/v1/reports/{reportId}", h.GetReport)
	registerSessionRoute(apiMux, "POST /api/v1/reports", h.SubmitReport)

	docsMux := http.NewServeMux()
	docsMux.HandleFunc("GET /api/v1/openapi.yaml", openapiSpecHandler)
	docsMux.HandleFunc("GET /api/v1/docs", swaggerDocsHandler)

	topMux := http.NewServeMux()
	topMux.Handle("/api/v1/docs", docsMux)
	topMux.Handle("/api/v1/openapi.yaml", docsMux)
	topMux.Handle("/", apiMux)
	return topMux
}

func registerUserRoute(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	mux.Handle(pattern, withUserContext(handler))
}

func registerSessionRoute(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	mux.Handle(pattern, withUserContext(withSessionContext(handler)))
}
