package domain

import (
	"time"

	"github.com/google/uuid"
)

type NoteTag string

const (
	NoteTagStrange    NoteTag = "strange"
	NoteTagSuspicious NoteTag = "suspicious"
	NoteTagLie        NoteTag = "lie"
	NoteTagKey        NoteTag = "key"
)

func (t NoteTag) Label() string {
	switch t {
	case NoteTagStrange:
		return "странно"
	case NoteTagSuspicious:
		return "подозрительно"
	case NoteTagLie:
		return "ложь"
	case NoteTagKey:
		return "ключевое"
	}
	return string(t)
}

type NotebookEntryType string

const (
	NotebookEntryTypeStatement          NotebookEntryType = "statement"
	NotebookEntryTypeAnalysis           NotebookEntryType = "analysis"
	NotebookEntryTypeAlibiCheck         NotebookEntryType = "alibi_check"
	NotebookEntryTypeCameraRequest      NotebookEntryType = "camera_request"
	NotebookEntryTypeTransactionRequest NotebookEntryType = "transaction_request"
)

func (t NotebookEntryType) Label() string {
	switch t {
	case NotebookEntryTypeStatement:
		return "показания"
	case NotebookEntryTypeAnalysis:
		return "анализ"
	case NotebookEntryTypeAlibiCheck:
		return "проверка алиби"
	case NotebookEntryTypeCameraRequest:
		return "запись с камер"
	case NotebookEntryTypeTransactionRequest:
		return "транзакция"
	}
	return string(t)
}

type NotebookEntry struct {
	ID          uuid.UUID         `json:"id"`
	Type        NotebookEntryType `json:"type"`
	CharacterID *uuid.UUID        `json:"character_id,omitempty"`
	Description string            `json:"description"`
	UserTags    []NoteTag         `json:"user_tags"`
	UserNote    *string           `json:"user_note,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
}

type ChronologyEventType string

const (
	ChronologyEventTypeCaseStarted      ChronologyEventType = "case_started"
	ChronologyEventTypeInterrogation    ChronologyEventType = "interrogation"
	ChronologyEventTypeLabAnalysis      ChronologyEventType = "lab_analysis"
	ChronologyEventTypeAlibiCheck       ChronologyEventType = "alibi_check"
	ChronologyEventTypeCameraReview     ChronologyEventType = "camera_review"
	ChronologyEventTypeTransactionCheck ChronologyEventType = "transaction_check"
)

func (t ChronologyEventType) Label() string {
	switch t {
	case ChronologyEventTypeCaseStarted:
		return "начало дела"
	case ChronologyEventTypeInterrogation:
		return "допрос"
	case ChronologyEventTypeLabAnalysis:
		return "экспертиза"
	case ChronologyEventTypeAlibiCheck:
		return "проверка алиби"
	case ChronologyEventTypeCameraReview:
		return "запись с камер"
	case ChronologyEventTypeTransactionCheck:
		return "транзакция"
	}
	return string(t)
}

type ChronologyEntry struct {
	ID        uuid.UUID           `json:"id"`
	EventType ChronologyEventType `json:"event_type"`
	Title     string              `json:"title"`
	Timestamp time.Time           `json:"timestamp"`
	Details   []NotebookEntry     `json:"details"`
}

type ActionReport struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Body        string     `json:"body"`
	EvidenceID  *uuid.UUID `json:"evidence_id,omitempty"`
	CharacterID *uuid.UUID `json:"character_id,omitempty"`
	Timestamp   time.Time  `json:"timestamp"`
}
