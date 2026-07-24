package domain

import "github.com/google/uuid"

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
	Type          CrimeType `bson:"crime_type"`
	Victim        string    `bson:"victim"`
	PerpetratorID uuid.UUID `bson:"perpetrator_id"`
	Motive        string    `bson:"motive"`
	Method        string    `bson:"method"`
	TimeOfCrime   string    `bson:"time_of_crime"`
}

type TimelineEntry struct {
	Time        string     `bson:"time"`
	Event       string     `bson:"event"`
	CharacterID *uuid.UUID `bson:"character_id,omitempty"`
}

type Timeline struct {
	Entries []TimelineEntry `bson:"entries"`
}

type EvidenceType string

const (
	EvidenceTypePhysical  EvidenceType = "physical"
	EvidenceTypeDigital   EvidenceType = "digital"
	EvidenceTypeDocument  EvidenceType = "document"
	EvidenceTypeTestimony EvidenceType = "testimony"
)

type Evidence struct {
	ID                  uuid.UUID    `bson:"_id"`
	SessionID           uuid.UUID    `bson:"session_id"`
	Name                string       `bson:"name"`
	Description         string       `bson:"description"`
	IconAsset           string       `bson:"icon_asset"`
	DetailedDescription string       `bson:"detailed_description"`
	Type                EvidenceType `bson:"type"`
}
