package http

import "net/http"

func NewRouter(h *Handlers) http.Handler {
	apiMux := http.NewServeMux()
	registerUserRoute := func(pattern string, handler http.HandlerFunc) {
		apiMux.Handle(pattern, withAuthentication(h.Authentication, h.Translator, handler))
	}
	registerSessionRoute := func(pattern string, handler http.HandlerFunc) {
		apiMux.Handle(pattern, withAuthentication(h.Authentication, h.Translator, withSessionContext(h.Session, h.Translator, handler)))
	}

	apiMux.HandleFunc("POST /v1/auth/anonymous", h.RegisterAnonymous)
	apiMux.HandleFunc("POST /v1/app/version", h.CheckVersion)

	registerUserRoute("POST /v1/sessions", h.CreateSession)
	registerUserRoute("GET /v1/sessions/history", h.ListHistory)
	registerUserRoute("GET /v1/sessions/current", h.GetSession)
	registerUserRoute("GET /v1/sessions/{id}", h.GetSessionByID)

	registerSessionRoute("GET /v1/evidence", h.ListEvidence)
	registerSessionRoute("GET /v1/evidence/{evId}", h.GetEvidence)
	registerSessionRoute("GET /v1/characters", h.ListCharacters)
	registerSessionRoute("GET /v1/characters/{charId}", h.GetCharacter)
	registerSessionRoute("GET /v1/chronology", h.GetChronology)
	registerSessionRoute("GET /v1/interrogations/active", h.GetActiveInterrogation)
	registerSessionRoute("POST /v1/interrogations", h.CreateInterrogation)
	registerSessionRoute("GET /v1/interrogations/{interId}", h.GetInterrogation)
	registerSessionRoute("POST /v1/interrogations/{interId}/messages", h.AddInterrogationMessage)
	registerSessionRoute("GET /v1/interrogations/{interId}/messages", h.GetInterrogationMessages)
	registerSessionRoute("PATCH /v1/interrogations/{interId}/complete", h.CompleteInterrogation)
	registerSessionRoute("PATCH /v1/chronology/{chronId}/notes/{noteId}", h.UpdateNotebookEntry)
	registerSessionRoute("POST /v1/actions/dna-analysis", h.DNAAnalysis)
	registerSessionRoute("POST /v1/actions/fingerprints", h.FingerprintsCheck)
	registerSessionRoute("POST /v1/actions/alibi-check", h.AlibiCheck)
	registerSessionRoute("POST /v1/actions/camera-review", h.CameraReview)
	registerSessionRoute("POST /v1/actions/call-history", h.CallHistory)
	registerSessionRoute("POST /v1/actions/transactions", h.TransactionCheck)
	registerSessionRoute("GET /v1/reports", h.ListReports)
	registerSessionRoute("POST /v1/reports", h.SubmitReport)

	docsMux := http.NewServeMux()
	docsMux.HandleFunc("GET /v1/openapi.yaml", openapiSpecHandler)
	docsMux.HandleFunc("GET /v1/docs", swaggerDocsHandler)

	topMux := http.NewServeMux()
	topMux.Handle("/v1/docs", docsMux)
	topMux.Handle("/v1/openapi.yaml", docsMux)
	topMux.Handle("/", apiMux)
	return topMux
}
