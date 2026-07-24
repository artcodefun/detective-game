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

type NotebookEntryType string

const (
	NotebookEntryTypeStatement          NotebookEntryType = "statement"
	NotebookEntryTypeAnalysis           NotebookEntryType = "analysis"
	NotebookEntryTypeAlibiCheck         NotebookEntryType = "alibi_check"
	NotebookEntryTypeCameraRequest      NotebookEntryType = "camera_request"
	NotebookEntryTypeTransactionRequest NotebookEntryType = "transaction_request"
)

type NotebookEntry struct {
	ID          uuid.UUID         `bson:"id"`
	Type        NotebookEntryType `bson:"type"`
	CharacterID *uuid.UUID        `bson:"character_id,omitempty"`
	Description string            `bson:"description"`
	UserTags    []NoteTag         `bson:"user_tags"`
	UserNote    *string           `bson:"user_note,omitempty"`
	Timestamp   time.Time         `bson:"timestamp"`
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

type ChronologyEntry struct {
	ID        uuid.UUID           `bson:"id"`
	SessionID uuid.UUID           `bson:"session_id"`
	EventType ChronologyEventType `bson:"event_type"`
	Title     string              `bson:"title"`
	Timestamp time.Time           `bson:"timestamp"`
	Details   []NotebookEntry     `bson:"details"`
}

type ActionType string

const (
	ActionTypeDNAAnalysis      ActionType = "dna_analysis"
	ActionTypeFingerprints     ActionType = "fingerprints"
	ActionTypeAlibiCheck       ActionType = "alibi_check"
	ActionTypeCameraReview     ActionType = "camera_review"
	ActionTypeCallHistory      ActionType = "call_history"
	ActionTypeTransactionCheck ActionType = "transaction_check"
)

func (t ActionType) Cost() int {
	switch t {
	case ActionTypeDNAAnalysis, ActionTypeFingerprints, ActionTypeAlibiCheck:
		return 1
	default:
		return 2
	}
}

type ActionReport struct {
	ID          uuid.UUID  `bson:"_id"`
	SessionID   uuid.UUID  `bson:"session_id"`
	Type        ActionType `bson:"type"`
	Body        string     `bson:"body"`
	EvidenceID  *uuid.UUID `bson:"evidence_id,omitempty"`
	CharacterID *uuid.UUID `bson:"character_id,omitempty"`
	Timestamp   time.Time  `bson:"timestamp"`
}
