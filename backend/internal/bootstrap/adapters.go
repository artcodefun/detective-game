package bootstrap

import (
	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/infrastructure/llm"
	"github.com/artcodefun/detective-game/backend/internal/infrastructure/storage"
)

type Adapters struct {
	Users          ports.UserRepository
	Sessions       ports.SessionRepository
	Characters     ports.CharacterRepository
	Interrogations ports.InterrogationRepository
	Chat           ports.ChatMessageRepository
	Evidence       ports.EvidenceRepository
	Reports        ports.ActionReportRepository
	Chronology     ports.ChronologyRepository

	ReadSessions ports.SessionReadRepository
	ReadChars    ports.CharacterReadRepository
	ReadEvidence ports.EvidenceReadRepository
	ReadReports  ports.ReportReadRepository
	ReadChron    ports.ChronologyReadRepository
	ReadChat     ports.ChatMessageReadRepository

	LLM        ports.LlmService
	Prototypes ports.CharacterPrototypeRepository
}

func NewAdapters() *Adapters {
	store := storage.NewInMemoryStore()
	return &Adapters{
		Users:          store,
		Sessions:       store,
		Characters:     store,
		Interrogations: store,
		Chat:           store,
		Evidence:       store,
		Reports:        store,
		Chronology:     store,

		ReadSessions: store,
		ReadChars:    store,
		ReadEvidence: store,
		ReadReports:  store,
		ReadChron:    store,
		ReadChat:     store,

		LLM:        llm.NewMockLlmService(),
		Prototypes: store,
	}
}
