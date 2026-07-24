package domain

import "github.com/google/uuid"

// CrimeType represents type of crime.
type CrimeType string

const (
	CrimeTypeMurder     CrimeType = "murder"
	CrimeTypeTheft      CrimeType = "theft"
	CrimeTypeFraud      CrimeType = "fraud"
	CrimeTypeArson      CrimeType = "arson"
	CrimeTypeKidnapping CrimeType = "kidnapping"
	CrimeTypeBlackmail  CrimeType = "blackmail"
)

type Crime struct {
	Type          CrimeType `json:"crime_type"`
	Victim        string    `json:"victim"`
	PerpetratorID uuid.UUID `json:"perpetrator_id"`
	Motive        string    `json:"motive"`
	Method        string    `json:"method"`
	TimeOfCrime   string    `json:"time_of_crime"`
}

type TimelineEntry struct {
	Time        string     `json:"time"`
	Event       string     `json:"event"`
	CharacterID *uuid.UUID `json:"character_id,omitempty"`
}

type Timeline struct {
	Entries []TimelineEntry `json:"entries"`
}

type EvidenceType string

const (
	EvidenceTypePhysical  EvidenceType = "physical"
	EvidenceTypeDigital   EvidenceType = "digital"
	EvidenceTypeDocument  EvidenceType = "document"
	EvidenceTypeTestimony EvidenceType = "testimony"
)

type Evidence struct {
	ID                  uuid.UUID    `json:"id"`
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	IconAsset           string       `json:"icon_asset"`
	DetailedDescription string       `json:"detailed_description"`
	Type                EvidenceType `json:"type"`
}
