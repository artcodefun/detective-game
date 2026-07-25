package readmodels

import (
	"github.com/artcodefun/detective-game/backend/internal/domain"
)

func SessionFromDomain(session *domain.Session) *Session {
	return &Session{
		ID:           session.ID,
		UserID:       session.UserID,
		Crime:        CrimeFromDomain(session.Crime),
		Timeline:     TimelineFromDomain(session.Timeline),
		CaseName:     session.CaseName,
		CaseBrief:    session.CaseBrief,
		ActionPoints: session.ActionPoints,
		Phase:        string(session.Phase),
		GameResult:   GameResultFromDomain(session.GameResult),
		CreatedAt:    session.CreatedAt,
	}
}

func CrimeFromDomain(crime domain.Crime) Crime {
	return Crime{
		CrimeType:     string(crime.Type),
		Victim:        crime.Victim,
		PerpetratorID: crime.PerpetratorID,
		Motive:        crime.Motive,
		Method:        crime.Method,
		TimeOfCrime:   crime.TimeOfCrime,
	}
}

func TimelineFromDomain(t domain.Timeline) Timeline {
	entries := make([]TimelineEntry, len(t.Entries))
	for i, e := range t.Entries {
		entries[i] = TimelineEntry{
			Time:        e.Time,
			Event:       e.Event,
			CharacterID: e.CharacterID,
		}
	}
	return Timeline{Entries: entries}
}

func GameResultFromDomain(r *domain.GameResult) *GameResult {
	if r == nil {
		return nil
	}
	return &GameResult{
		PlayerReport: FinalReport{
			Who:      r.PlayerReport.Who,
			Why:      r.PlayerReport.Why,
			How:      r.PlayerReport.How,
			When:     r.PlayerReport.When,
			Evidence: r.PlayerReport.Evidence,
		},
		Breakdown: ScoreBreakdown{
			WhoCorrect:      r.Breakdown.WhoCorrect,
			WhyCorrect:      r.Breakdown.WhyCorrect,
			HowCorrect:      r.Breakdown.HowCorrect,
			WhenCorrect:     r.Breakdown.WhenCorrect,
			EvidenceCorrect: r.Breakdown.EvidenceCorrect,
		},
		NarrativeFeedback: r.NarrativeFeedback,
		BreakdownDetails:  r.BreakdownDetails,
	}
}

func EvidenceFromDomain(ev *domain.Evidence) *Evidence {
	return &Evidence{
		ID:                  ev.ID,
		Name:                ev.Name,
		Description:         ev.Description,
		IconAsset:           ev.IconAsset,
		DetailedDescription: ev.DetailedDescription,
		Type:                string(ev.Type),
	}
}

func ActionReportFromDomain(r *domain.ActionReport) *ActionReport {
	title := actionTypeTitle(r.Type)
	desc := r.Body
	if len(desc) > 80 {
		desc = desc[:80] + "..."
	}
	return &ActionReport{
		ID:          r.ID,
		Type:        string(r.Type),
		Title:       title,
		Description: desc,
		Body:        r.Body,
		EvidenceID:  r.EvidenceID,
		CharacterID: r.CharacterID,
		Timestamp:   r.Timestamp,
	}
}

func actionTypeTitle(t domain.ActionType) string {
	switch t {
	case domain.ActionTypeDNAAnalysis:
		return "Анализ ДНК"
	case domain.ActionTypeFingerprints:
		return "Отпечатки пальцев"
	case domain.ActionTypeAlibiCheck:
		return "Проверка алиби"
	case domain.ActionTypeCameraReview:
		return "Записи с камер"
	case domain.ActionTypeCallHistory:
		return "История звонков"
	case domain.ActionTypeTransactionCheck:
		return "Банковские операции"
	}
	return string(t)
}

func ChronologyEntryFromDomain(c *domain.ChronologyEntry) *ChronologyEntry {
	details := make([]NotebookEntry, len(c.Details))
	for i, d := range c.Details {
		tags := make([]string, len(d.UserTags))
		for j, t := range d.UserTags {
			tags[j] = string(t)
		}
		details[i] = NotebookEntry{
			ID:          d.ID,
			Type:        string(d.Type),
			CharacterID: nil,
			Description: d.Description,
			UserTags:    tags,
			UserNote:    d.UserNote,
			Timestamp:   d.Timestamp,
		}
	}
	return &ChronologyEntry{
		ID:        c.ID,
		EventType: string(c.EventType),
		Title:     c.Title,
		Timestamp: c.Timestamp,
		Details:   details,
	}
}

func ChatMessageFromDomain(msg domain.ChatMessage) ChatMessage {
	return ChatMessage{
		ID:              msg.ID,
		SessionID:       msg.SessionID,
		InterrogationID: msg.InterrogationID,
		FromUser:        msg.FromUser,
		Text:            msg.Text,
		Statements:      msg.Statements,
		AttitudeDelta:   msg.AttitudeDelta,
		Timestamp:       msg.Timestamp,
	}
}

func InterrogationFromDomain(inter *domain.Interrogation) *Interrogation {
	return &Interrogation{
		ID:          inter.ID,
		SessionID:   inter.SessionID,
		CharacterID: inter.CharacterID,
		Phase:       string(inter.Phase),
		CreatedAt:   inter.CreatedAt,
		CompletedAt: inter.CompletedAt,
	}
}

func CharacterFromDomain(c domain.Character) Character {
	return Character{
		ID:                      c.ID,
		SessionID:               c.SessionID,
		Name:                    c.Name,
		Age:                     c.Age,
		Profession:              c.Profession,
		Personality:             c.Personality,
		Gender:                  string(c.Gender),
		Knowledge:               CharacterKnowledgeFromDomain(c.Knowledge),
		Secrets:                 c.Secrets,
		Relationships:           c.Relationships,
		Memories:                MemoriesFromDomain(c.Memories),
		Trust:                   c.Trust,
		InterrogationsRemaining: c.InterrogationsRemaining,
	}
}

func CharacterKnowledgeFromDomain(k domain.CharacterKnowledge) CharacterKnowledge {
	return CharacterKnowledge{
		KnownFacts:   k.KnownFacts,
		PartialFacts: k.PartialFacts,
		FalseBeliefs: k.FalseBeliefs,
	}
}

func MemoriesFromDomain(memories []domain.Memory) []Memory {
	out := make([]Memory, len(memories))
	for i, m := range memories {
		out[i] = Memory{
			ID:        m.ID,
			Content:   m.Content,
			IsTrue:    m.IsTrue,
			Timestamp: m.Timestamp,
		}
	}
	return out
}
