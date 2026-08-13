package http

import (
	"encoding/json"
	"net/http"

	"github.com/artcodefun/detective-game/backend/internal/application"
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/artcodefun/detective-game/backend/internal/interfaces/http/dtos"
	"github.com/google/uuid"
)

type Handlers struct {
	Scenario      application.ScenarioCommands
	Interrogation application.InterrogationCommands
	Evaluation    application.EvaluationCommands
	Actions       application.ActionCommands
	Notebook      application.NotebookCommands

	Session    application.SessionQueries
	Character  application.CharacterQueries
	Evidence   application.EvidenceQueries
	Chronology application.ChronologyQueries
	Chat       application.ChatQueries

	Translator ports.Translator
}

// POST /api/v1/sessions
func (h *Handlers) CreateSession(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	sessionID, err := h.Scenario.CreateSession(r.Context(), actor)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"session_id": sessionID.String()})
}

// GET /api/v1/sessions/history
func (h *Handlers) ListHistory(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	history, err := h.Session.ListHistory(r.Context(), actor)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	if history == nil {
		history = make([]*readmodels.Session, 0)
	}
	writeJSON(w, http.StatusOK, dtos.SessionsFromReadModels(history))
}

// GET /api/v1/sessions/current
func (h *Handlers) GetSession(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	session, err := h.Session.GetSession(r.Context(), actor)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.SessionFromReadModel(session))
}

// GET /api/v1/sessions/{id}
func (h *Handlers) GetSessionByID(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	sessionID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid session id")
		return
	}
	session, err := h.Session.GetSessionByID(r.Context(), actor, sessionID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.SessionFromReadModel(session))
}

// GET /api/v1/evidence
func (h *Handlers) ListEvidence(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	ev, err := h.Evidence.ListEvidence(r.Context(), actor)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.EvidenceListFromReadModels(ev))
}

// GET /api/v1/evidence/{evId}
func (h *Handlers) GetEvidence(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	evID, err := uuid.Parse(r.PathValue("evId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid evidence id")
		return
	}
	ev, err := h.Evidence.GetEvidence(r.Context(), actor, evID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.EvidenceFromReadModel(ev))
}

// GET /api/v1/characters
func (h *Handlers) ListCharacters(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	chars, err := h.Character.ListCharacters(r.Context(), actor)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.CharactersFromReadModels(chars))
}

// GET /api/v1/characters/{charId}
func (h *Handlers) GetCharacter(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	charID, err := uuid.Parse(r.PathValue("charId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid character id")
		return
	}
	char, err := h.Character.GetCharacter(r.Context(), actor, charID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.CharacterFromReadModel(char))
}

// GET /api/v1/chronology
func (h *Handlers) GetChronology(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	chron, err := h.Chronology.GetChronology(r.Context(), actor)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.ChronologyEntriesFromReadModel(chron, h.Translator))
}

// PATCH /api/v1/chronology/{chronId}/notes/{noteId}
func (h *Handlers) UpdateNotebookEntry(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())

	chronID, err := uuid.Parse(r.PathValue("chronId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid chronology id")
		return
	}
	entryID, err := uuid.Parse(r.PathValue("noteId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid note id")
		return
	}

	var body struct {
		Tags []domain.NoteTag `json:"tags"`
		Note *string          `json:"note,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if err := h.Notebook.UpdateNotebookEntry(r.Context(), actor, chronID, entryID, body.Tags, body.Note); err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GET /api/v1/interrogations/active
func (h *Handlers) GetActiveInterrogation(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	inter, err := h.Chat.GetActiveInterrogation(r.Context(), actor)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.InterrogationFromReadModel(inter))
}

// POST /api/v1/interrogations
func (h *Handlers) CreateInterrogation(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())

	var body struct {
		CharacterID uuid.UUID `json:"character_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	interID, err := h.Interrogation.Create(r.Context(), actor, body.CharacterID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}

	inter, err := h.Chat.GetInterrogation(r.Context(), actor, interID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}

	writeJSON(w, http.StatusCreated, dtos.InterrogationFromReadModel(inter))
}

// GET /api/v1/interrogations/{interId}
func (h *Handlers) GetInterrogation(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())

	interID, err := uuid.Parse(r.PathValue("interId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid interrogation id")
		return
	}

	inter, err := h.Chat.GetInterrogation(r.Context(), actor, interID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}

	writeJSON(w, http.StatusOK, dtos.InterrogationFromReadModel(inter))
}

// POST /api/v1/interrogations/{interId}/messages
func (h *Handlers) AddInterrogationMessage(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())

	interID, err := uuid.Parse(r.PathValue("interId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid interrogation id")
		return
	}

	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if body.Message == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}

	msgID, err := h.Interrogation.AddMessage(r.Context(), actor, interID, body.Message)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}

	msg, err := h.Chat.GetChatMessage(r.Context(), actor, msgID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}

	writeJSON(w, http.StatusOK, dtos.ChatMessageFromReadModel(msg))
}

// GET /api/v1/interrogations/{interId}/messages
func (h *Handlers) GetInterrogationMessages(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())

	interID, err := uuid.Parse(r.PathValue("interId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid interrogation id")
		return
	}

	messages, err := h.Chat.ListChatByInterrogation(r.Context(), actor, interID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	if messages == nil {
		messages = make([]*readmodels.ChatMessage, 0)
	}

	writeJSON(w, http.StatusOK, dtos.ChatMessagesFromReadModels(messages))
}

// PATCH /api/v1/interrogations/{interId}/complete
func (h *Handlers) CompleteInterrogation(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())

	interID, err := uuid.Parse(r.PathValue("interId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid interrogation id")
		return
	}

	if err := h.Interrogation.Complete(r.Context(), actor, interID); err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /api/v1/actions/dna-analysis
func (h *Handlers) DNAAnalysis(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	var body struct {
		EvidenceID uuid.UUID `json:"evidence_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	reportID, err := h.Actions.DNAAnalysis(r.Context(), actor, body.EvidenceID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	result, err := h.Evidence.GetReport(r.Context(), actor, reportID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.ActionReportFromReadModel(result, h.Translator))
}

// POST /api/v1/actions/fingerprints
func (h *Handlers) FingerprintsCheck(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	var body struct {
		EvidenceID uuid.UUID `json:"evidence_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	reportID, err := h.Actions.FingerprintsCheck(r.Context(), actor, body.EvidenceID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	result, err := h.Evidence.GetReport(r.Context(), actor, reportID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.ActionReportFromReadModel(result, h.Translator))
}

// POST /api/v1/actions/alibi-check
func (h *Handlers) AlibiCheck(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	var body struct {
		CharacterID uuid.UUID `json:"character_id"`
		AlibiText   string    `json:"alibi_text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	reportID, err := h.Actions.AlibiCheck(r.Context(), actor, body.CharacterID, body.AlibiText)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	result, err := h.Evidence.GetReport(r.Context(), actor, reportID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.ActionReportFromReadModel(result, h.Translator))
}

// POST /api/v1/actions/camera-review
func (h *Handlers) CameraReview(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	reportID, err := h.Actions.CameraReview(r.Context(), actor)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	result, err := h.Evidence.GetReport(r.Context(), actor, reportID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.ActionReportFromReadModel(result, h.Translator))
}

// POST /api/v1/actions/call-history
func (h *Handlers) CallHistory(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	var body struct {
		CharacterID uuid.UUID `json:"character_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	reportID, err := h.Actions.CallHistory(r.Context(), actor, body.CharacterID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	result, err := h.Evidence.GetReport(r.Context(), actor, reportID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.ActionReportFromReadModel(result, h.Translator))
}

// POST /api/v1/actions/transactions
func (h *Handlers) TransactionCheck(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	var body struct {
		CharacterID uuid.UUID `json:"character_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	reportID, err := h.Actions.TransactionCheck(r.Context(), actor, body.CharacterID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	result, err := h.Evidence.GetReport(r.Context(), actor, reportID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.ActionReportFromReadModel(result, h.Translator))
}

// GET /api/v1/reports/{reportId}
func (h *Handlers) GetReport(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())
	reportID, err := uuid.Parse(r.PathValue("reportId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid report id")
		return
	}
	report, err := h.Evidence.GetReport(r.Context(), actor, reportID)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}
	writeJSON(w, http.StatusOK, dtos.ActionReportFromReadModel(report, h.Translator))
}

// POST /api/v1/reports
func (h *Handlers) SubmitReport(w http.ResponseWriter, r *http.Request) {
	actor := ActorFromContext(r.Context())

	var body struct {
		Who      string `json:"who"`
		Why      string `json:"why"`
		How      string `json:"how"`
		When     string `json:"when"`
		Evidence string `json:"evidence"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	report := domain.FinalReport{
		Who:      body.Who,
		Why:      body.Why,
		How:      body.How,
		When:     body.When,
		Evidence: body.Evidence,
	}

	if err := h.Evaluation.SubmitReport(r.Context(), actor, report); err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}

	result, err := h.Session.GetGameResult(r.Context(), actor)
	if err != nil {
		writeAppError(w, r.Context(), h.Translator, err)
		return
	}

	writeJSON(w, http.StatusOK, dtos.GameResultFromReadModel(result))
}
