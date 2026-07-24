package storage

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/artcodefun/detective-game/backend/internal/application/readmodels"
	"github.com/artcodefun/detective-game/backend/internal/domain"
	"github.com/google/uuid"
)

type InMemoryStore struct {
	mu             sync.RWMutex
	sessions       map[string]*domain.Session
	characters     map[string]*domain.Character
	interrogations map[string]*domain.Interrogation
	evidence       map[string][]*domain.Evidence
	reports        map[string][]*domain.ActionReport
	chronology     map[string][]*domain.ChronologyEntry
	chatMessages   map[string][]*domain.ChatMessage
	prototypes     []domain.CharacterPrototype
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		sessions:       make(map[string]*domain.Session),
		characters:     make(map[string]*domain.Character),
		interrogations: make(map[string]*domain.Interrogation),
		evidence:       make(map[string][]*domain.Evidence),
		reports:        make(map[string][]*domain.ActionReport),
		chronology:     make(map[string][]*domain.ChronologyEntry),
		chatMessages:   make(map[string][]*domain.ChatMessage),
		prototypes: []domain.CharacterPrototype{
			{ID: 1, Name: "Иван Петров", Age: 55, Profession: "дворецкий",
				ImagePath:   "assets/characters/ivan_petrov.png",
				Personality: "Консервативный, преданный семье, скрытный. Говорит медленно, с расстановкой. Предпочитает отмалчиваться, но если задеть за живое — срывается.",
				AudioToneID: "tone_male_deep"},
			{ID: 2, Name: "Елена Соколова", Age: 42, Profession: "домохозяйка",
				ImagePath:   "assets/characters/elena_sokolova.png",
				Personality: "Эмоциональная, вспыльчивая, но ранимая. Говорит быстро, часто перебивает. Хочет казаться безразличной, но на деле очень переживает.",
				AudioToneID: "tone_female_high"},
			{ID: 3, Name: "Майкл Браун", Age: 48, Profession: "деловой партнёр",
				ImagePath:   "assets/characters/michael_brown.png",
				Personality: "Харизматичный, уверенный в себе, умело манипулирует. Говорит спокойно, с лёгкой усмешкой. Всегда контролирует эмоции.",
				AudioToneID: "tone_male_mid"},
			{ID: 4, Name: "Анна Коваль", Age: 29, Profession: "горничная",
				ImagePath:   "assets/characters/anna_koval.png",
				Personality: "Застенчивая, тревожная, боится потерять работу. Говорит тихо, запинается. Старается быть незаметной, но глаза выдают страх.",
				AudioToneID: "tone_female_soft"},
			{ID: 5, Name: "Дмитрий Орлов", Age: 61, Profession: "адвокат",
				ImagePath:   "assets/characters/dmitry_orlov.png",
				Personality: "Циничный, расчётливый, за словом в карман не лезет. Говорит чётко, рублеными фразами. Привык контролировать ситуацию.",
				AudioToneID: "tone_male_raspy"},
		},
	}
}

func charKey(sessionID uuid.UUID, prototypeID int) string {
	return sessionID.String() + ":" + fmt.Sprintf("%d", prototypeID)
}

func chatKey(sessionID uuid.UUID, characterID uuid.UUID) string {
	return sessionID.String() + ":" + characterID.String()
}

// CharacterPrototypeRepository

func (s *InMemoryStore) GetAll(_ context.Context) ([]domain.CharacterPrototype, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.CharacterPrototype, len(s.prototypes))
	copy(out, s.prototypes)
	return out, nil
}

func (s *InMemoryStore) GetRandom(_ context.Context, count int) ([]domain.CharacterPrototype, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if count >= len(s.prototypes) {
		out := make([]domain.CharacterPrototype, len(s.prototypes))
		copy(out, s.prototypes)
		return out, nil
	}
	shuffled := make([]domain.CharacterPrototype, len(s.prototypes))
	copy(shuffled, s.prototypes)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled[:count], nil
}

func (s *InMemoryStore) ByID(_ context.Context, id int) (*domain.CharacterPrototype, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.prototypes {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, nil
}

// SessionRepository

func (s *InMemoryStore) Create(_ context.Context, session *domain.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID.String()] = session
	return nil
}

func (s *InMemoryStore) FindByID(_ context.Context, id uuid.UUID) (*domain.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[id.String()]
	if !ok {
		return nil, fmt.Errorf("session %s not found", id)
	}
	return session, nil
}

func (s *InMemoryStore) Update(_ context.Context, session *domain.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID.String()] = session
	return nil
}

// CharacterRepository

func (s *InMemoryStore) CreateCharacter(_ context.Context, character *domain.Character) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := charKey(character.SessionID, character.Prototype.ID)
	s.characters[key] = character
	return nil
}

func (s *InMemoryStore) FindCharacterBySessionAndID(_ context.Context, sessionID uuid.UUID, prototypeID int) (*domain.Character, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := charKey(sessionID, prototypeID)
	c, ok := s.characters[key]
	if !ok {
		return nil, fmt.Errorf("character %d not found in session %s", prototypeID, sessionID)
	}
	return c, nil
}

func (s *InMemoryStore) FindCharacterByID(_ context.Context, sessionID uuid.UUID, characterID uuid.UUID) (*domain.Character, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	prefix := sessionID.String() + ":"
	for k, v := range s.characters {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix && v.ID == characterID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("character %s not found in session %s", characterID, sessionID)
}

func (s *InMemoryStore) FindCharactersBySession(_ context.Context, sessionID uuid.UUID) ([]*domain.Character, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	prefix := sessionID.String() + ":"
	var result []*domain.Character
	for k, v := range s.characters {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			result = append(result, v)
		}
	}
	return result, nil
}

func (s *InMemoryStore) UpdateCharacter(_ context.Context, character *domain.Character) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := charKey(character.SessionID, character.Prototype.ID)
	s.characters[key] = character
	return nil
}

// InterrogationRepository

func (s *InMemoryStore) CreateInterrogation(_ context.Context, interrogation *domain.Interrogation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.interrogations[interrogation.ID.String()] = interrogation
	return nil
}

func (s *InMemoryStore) FindInterrogationByID(_ context.Context, id uuid.UUID) (*domain.Interrogation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	inter, ok := s.interrogations[id.String()]
	if !ok {
		return nil, fmt.Errorf("interrogation %s not found", id)
	}
	return inter, nil
}

func (s *InMemoryStore) FindActiveBySession(_ context.Context, sessionID uuid.UUID) (*domain.Interrogation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, inter := range s.interrogations {
		if inter.SessionID == sessionID && inter.IsActive() {
			return inter, nil
		}
	}
	return nil, nil
}

func (s *InMemoryStore) UpdateInterrogation(_ context.Context, interrogation *domain.Interrogation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.interrogations[interrogation.ID.String()] = interrogation
	return nil
}

// ChatMessageRepository

func (s *InMemoryStore) AppendChatMessage(_ context.Context, msg *domain.ChatMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := msg.InterrogationID.String()
	s.chatMessages[key] = append(s.chatMessages[key], msg)
	return nil
}

func (s *InMemoryStore) FindChatByInterrogation(_ context.Context, interrogationID uuid.UUID) ([]*domain.ChatMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*domain.ChatMessage
	for _, msgs := range s.chatMessages {
		for _, msg := range msgs {
			if msg.InterrogationID == interrogationID {
				result = append(result, msg)
			}
		}
	}
	if result == nil {
		return []*domain.ChatMessage{}, nil
	}
	return result, nil
}

func (s *InMemoryStore) FindChatByID(_ context.Context, id uuid.UUID) (*domain.ChatMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, msgs := range s.chatMessages {
		for _, msg := range msgs {
			if msg.ID == id {
				return msg, nil
			}
		}
	}
	return nil, fmt.Errorf("chat message %s not found", id)
}

// EvidenceRepository

func (s *InMemoryStore) AppendEvidence(_ context.Context, sessionID uuid.UUID, evidence *domain.Evidence) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := sessionID.String()
	s.evidence[key] = append(s.evidence[key], evidence)
	return nil
}

func (s *InMemoryStore) FindEvidenceBySession(_ context.Context, sessionID uuid.UUID) ([]*domain.Evidence, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := sessionID.String()
	ev := s.evidence[key]
	if ev == nil {
		return []*domain.Evidence{}, nil
	}
	return ev, nil
}

func (s *InMemoryStore) FindEvidenceByID(_ context.Context, sessionID uuid.UUID, evidenceID uuid.UUID) (*domain.Evidence, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := sessionID.String()
	for _, ev := range s.evidence[key] {
		if ev.ID == evidenceID {
			return ev, nil
		}
	}
	return nil, fmt.Errorf("evidence %s not found in session %s", evidenceID, sessionID)
}

// ActionReportRepository

func (s *InMemoryStore) AppendReport(_ context.Context, sessionID uuid.UUID, report *domain.ActionReport) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := sessionID.String()
	s.reports[key] = append(s.reports[key], report)
	return nil
}

func (s *InMemoryStore) FindReportByID(_ context.Context, reportID uuid.UUID) (*domain.ActionReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, reports := range s.reports {
		for _, r := range reports {
			if r.ID == reportID {
				return r, nil
			}
		}
	}
	return nil, fmt.Errorf("report %s not found", reportID)
}

func (s *InMemoryStore) FindReportsBySession(_ context.Context, sessionID uuid.UUID) ([]*domain.ActionReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := sessionID.String()
	reports := s.reports[key]
	if reports == nil {
		return []*domain.ActionReport{}, nil
	}
	return reports, nil
}

// ChronologyRepository

func (s *InMemoryStore) AppendChronologyEntry(_ context.Context, sessionID uuid.UUID, entry *domain.ChronologyEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := sessionID.String()
	s.chronology[key] = append(s.chronology[key], entry)
	return nil
}

func (s *InMemoryStore) FindChronologyBySession(_ context.Context, sessionID uuid.UUID) ([]*domain.ChronologyEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := sessionID.String()
	entries := s.chronology[key]
	if entries == nil {
		return []*domain.ChronologyEntry{}, nil
	}
	return entries, nil
}

func (s *InMemoryStore) UpdateChronologyEntry(_ context.Context, sessionID uuid.UUID, chronologyID uuid.UUID, entryID uuid.UUID, tags []domain.NoteTag, note *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := sessionID.String()
	for i := range s.chronology[key] {
		if s.chronology[key][i].ID != chronologyID {
			continue
		}
		for j := range s.chronology[key][i].Details {
			if s.chronology[key][i].Details[j].ID != entryID {
				continue
			}
			s.chronology[key][i].Details[j].UserTags = tags
			s.chronology[key][i].Details[j].UserNote = note
			return nil
		}
	}
	return fmt.Errorf("entry %s not found in chronology %s", entryID, chronologyID)
}

// SessionReadRepository

func (s *InMemoryStore) GetSession(_ context.Context, sessionID uuid.UUID) (*readmodels.Session, error) {
	session, err := s.FindByID(context.Background(), sessionID)
	if err != nil {
		return nil, err
	}
	return readmodels.SessionFromDomain(session), nil
}

func (s *InMemoryStore) ListCharacters(_ context.Context, sessionID uuid.UUID) ([]*readmodels.Character, error) {
	chars, err := s.FindCharactersBySession(context.Background(), sessionID)
	if err != nil {
		return nil, err
	}
	result := make([]*readmodels.Character, len(chars))
	for i, c := range chars {
		char := readmodels.CharacterFromDomain(*c)
		result[i] = &char
	}
	return result, nil
}

func (s *InMemoryStore) GetCharacter(_ context.Context, sessionID uuid.UUID, characterID uuid.UUID) (*readmodels.Character, error) {
	c, err := s.FindCharacterByID(context.Background(), sessionID, characterID)
	if err != nil {
		return nil, err
	}
	char := readmodels.CharacterFromDomain(*c)
	return &char, nil
}

func (s *InMemoryStore) ListChatByInterrogation(_ context.Context, interrogationID uuid.UUID) ([]*readmodels.ChatMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*readmodels.ChatMessage
	for _, msgs := range s.chatMessages {
		for _, m := range msgs {
			if m.InterrogationID == interrogationID {
				msg := readmodels.ChatMessageFromDomain(*m)
				result = append(result, &msg)
			}
		}
	}
	if result == nil {
		return []*readmodels.ChatMessage{}, nil
	}
	return result, nil
}

func (s *InMemoryStore) GetInterrogation(_ context.Context, interrogationID uuid.UUID) (*readmodels.Interrogation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	inter, ok := s.interrogations[interrogationID.String()]
	if !ok {
		return nil, fmt.Errorf("interrogation %s not found", interrogationID)
	}
	return readmodels.InterrogationFromDomain(inter), nil
}

func (s *InMemoryStore) GetChatMessage(_ context.Context, messageID uuid.UUID) (*readmodels.ChatMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, msgs := range s.chatMessages {
		for _, msg := range msgs {
			if msg.ID == messageID {
				m := readmodels.ChatMessageFromDomain(*msg)
				return &m, nil
			}
		}
	}
	return nil, fmt.Errorf("chat message %s not found", messageID)
}

func (s *InMemoryStore) ListEvidence(_ context.Context, sessionID uuid.UUID) ([]*readmodels.Evidence, error) {
	ev, err := s.FindEvidenceBySession(context.Background(), sessionID)
	if err != nil {
		return nil, err
	}
	result := make([]*readmodels.Evidence, len(ev))
	for i, e := range ev {
		result[i] = readmodels.EvidenceFromDomain(e)
	}
	return result, nil
}

func (s *InMemoryStore) GetEvidence(_ context.Context, sessionID uuid.UUID, evidenceID uuid.UUID) (*readmodels.Evidence, error) {
	ev, err := s.FindEvidenceByID(context.Background(), sessionID, evidenceID)
	if err != nil {
		return nil, err
	}
	return readmodels.EvidenceFromDomain(ev), nil
}

func (s *InMemoryStore) GetReport(_ context.Context, reportID uuid.UUID) (*readmodels.ActionReport, error) {
	r, err := s.FindReportByID(context.Background(), reportID)
	if err != nil {
		return nil, err
	}
	return readmodels.ActionReportFromDomain(r), nil
}

func (s *InMemoryStore) ListReports(_ context.Context, sessionID uuid.UUID) ([]*readmodels.ActionReport, error) {
	reports, err := s.FindReportsBySession(context.Background(), sessionID)
	if err != nil {
		return nil, err
	}
	result := make([]*readmodels.ActionReport, len(reports))
	for i, r := range reports {
		result[i] = readmodels.ActionReportFromDomain(r)
	}
	return result, nil
}

func (s *InMemoryStore) GetChronology(_ context.Context, sessionID uuid.UUID) ([]*readmodels.ChronologyEntry, error) {
	entries, err := s.FindChronologyBySession(context.Background(), sessionID)
	if err != nil {
		return nil, err
	}
	items := make([]*readmodels.ChronologyEntry, len(entries))
	for i, c := range entries {
		items[i] = readmodels.ChronologyEntryFromDomain(c)
	}
	return items, nil
}

func (s *InMemoryStore) GetGameResult(_ context.Context, sessionID uuid.UUID) (*readmodels.GameResult, error) {
	session, err := s.FindByID(context.Background(), sessionID)
	if err != nil {
		return nil, err
	}
	if session.GameResult == nil {
		return nil, fmt.Errorf("session %s has no result yet", sessionID)
	}
	return readmodels.GameResultFromDomain(session.GameResult), nil
}

func (s *InMemoryStore) ListHistory(_ context.Context) ([]*readmodels.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]*readmodels.Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		items = append(items, readmodels.SessionFromDomain(session))
	}
	return items, nil
}


