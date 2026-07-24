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
	ID          uuid.UUID         `json:"id" bson:"id"`
	Type        NotebookEntryType `json:"type" bson:"type"`
	CharacterID *uuid.UUID        `json:"character_id,omitempty" bson:"character_id,omitempty"`
	Description string            `json:"description" bson:"description"`
	UserTags    []NoteTag         `json:"user_tags" bson:"user_tags"`
	UserNote    *string           `json:"user_note,omitempty" bson:"user_note,omitempty"`
	Timestamp   time.Time         `json:"timestamp" bson:"timestamp"`
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
	ID        uuid.UUID           `json:"id" bson:"id"`
	EventType ChronologyEventType `json:"event_type" bson:"event_type"`
	Title     string              `json:"title" bson:"title"`
	Timestamp time.Time           `json:"timestamp" bson:"timestamp"`
	Details   []NotebookEntry     `json:"details" bson:"details"`
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
	ID          uuid.UUID  `json:"id" bson:"_id"`
	Type        ActionType `json:"type" bson:"type"`
	Body        string     `json:"body" bson:"body"`
	EvidenceID  *uuid.UUID `json:"evidence_id,omitempty" bson:"evidence_id,omitempty"`
	CharacterID *uuid.UUID `json:"character_id,omitempty" bson:"character_id,omitempty"`
	Timestamp   time.Time  `json:"timestamp" bson:"timestamp"`
}
