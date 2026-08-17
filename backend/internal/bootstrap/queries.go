package bootstrap

import (
	"github.com/artcodefun/detective-game/backend/internal/application/queries"
)

type Queries struct {
	User       *queries.UserQueries
	Session    *queries.SessionQueries
	Character  *queries.CharacterQueries
	Evidence   *queries.EvidenceQueries
	Chronology *queries.ChronologyQueries
	Chat       *queries.ChatQueries
}

func NewQueries(a *Adapters) *Queries {
	return &Queries{
		User:       queries.NewUserQueries(a.Users),
		Session:    queries.NewSessionQueries(a.ReadSessions),
		Character:  queries.NewCharacterQueries(a.ReadChars),
		Evidence:   queries.NewEvidenceQueries(a.ReadEvidence, a.ReadReports),
		Chronology: queries.NewChronologyQueries(a.ReadChron),
		Chat:       queries.NewChatQueries(a.ReadChat, a.ReadChars),
	}
}
