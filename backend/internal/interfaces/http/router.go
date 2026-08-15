package http

import "net/http"

func NewRouter(h *Handlers) http.Handler {
	apiMux := http.NewServeMux()
	registerUserRoute := func(pattern string, handler http.HandlerFunc) {
		apiMux.Handle(pattern, withUserContext(handler))
	}
	registerSessionRoute := func(pattern string, handler http.HandlerFunc) {
		apiMux.Handle(pattern, withUserContext(withSessionContext(h.Session, h.Translator, handler)))
	}

	registerUserRoute("POST /api/v1/sessions", h.CreateSession)
	registerUserRoute("GET /api/v1/sessions/history", h.ListHistory)
	registerUserRoute("GET /api/v1/sessions/current", h.GetSession)
	registerUserRoute("GET /api/v1/sessions/{id}", h.GetSessionByID)

	registerSessionRoute("GET /api/v1/evidence", h.ListEvidence)
	registerSessionRoute("GET /api/v1/evidence/{evId}", h.GetEvidence)
	registerSessionRoute("GET /api/v1/characters", h.ListCharacters)
	registerSessionRoute("GET /api/v1/characters/{charId}", h.GetCharacter)
	registerSessionRoute("GET /api/v1/chronology", h.GetChronology)
	registerSessionRoute("GET /api/v1/interrogations/active", h.GetActiveInterrogation)
	registerSessionRoute("POST /api/v1/interrogations", h.CreateInterrogation)
	registerSessionRoute("GET /api/v1/interrogations/{interId}", h.GetInterrogation)
	registerSessionRoute("POST /api/v1/interrogations/{interId}/messages", h.AddInterrogationMessage)
	registerSessionRoute("GET /api/v1/interrogations/{interId}/messages", h.GetInterrogationMessages)
	registerSessionRoute("PATCH /api/v1/interrogations/{interId}/complete", h.CompleteInterrogation)
	registerSessionRoute("PATCH /api/v1/chronology/{chronId}/notes/{noteId}", h.UpdateNotebookEntry)
	registerSessionRoute("POST /api/v1/actions/dna-analysis", h.DNAAnalysis)
	registerSessionRoute("POST /api/v1/actions/fingerprints", h.FingerprintsCheck)
	registerSessionRoute("POST /api/v1/actions/alibi-check", h.AlibiCheck)
	registerSessionRoute("POST /api/v1/actions/camera-review", h.CameraReview)
	registerSessionRoute("POST /api/v1/actions/call-history", h.CallHistory)
	registerSessionRoute("POST /api/v1/actions/transactions", h.TransactionCheck)
	registerSessionRoute("GET /api/v1/reports", h.ListReports)
	registerSessionRoute("POST /api/v1/reports", h.SubmitReport)

	docsMux := http.NewServeMux()
	docsMux.HandleFunc("GET /api/v1/openapi.yaml", openapiSpecHandler)
	docsMux.HandleFunc("GET /api/v1/docs", swaggerDocsHandler)

	topMux := http.NewServeMux()
	topMux.Handle("/api/v1/docs", docsMux)
	topMux.Handle("/api/v1/openapi.yaml", docsMux)
	topMux.Handle("/", apiMux)
	return topMux
}
