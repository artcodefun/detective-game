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
	Type          CrimeType `json:"crime_type" bson:"crime_type"`
	Victim        string    `json:"victim" bson:"victim"`
	PerpetratorID uuid.UUID `json:"perpetrator_id" bson:"perpetrator_id"`
	Motive        string    `json:"motive" bson:"motive"`
	Method        string    `json:"method" bson:"method"`
	TimeOfCrime   string    `json:"time_of_crime" bson:"time_of_crime"`
}

type TimelineEntry struct {
	Time        string     `json:"time" bson:"time"`
	Event       string     `json:"event" bson:"event"`
	CharacterID *uuid.UUID `json:"character_id,omitempty" bson:"character_id,omitempty"`
}

type Timeline struct {
	Entries []TimelineEntry `json:"entries" bson:"entries"`
}

type EvidenceType string

const (
	EvidenceTypePhysical  EvidenceType = "physical"
	EvidenceTypeDigital   EvidenceType = "digital"
	EvidenceTypeDocument  EvidenceType = "document"
	EvidenceTypeTestimony EvidenceType = "testimony"
)

type Evidence struct {
	ID                  uuid.UUID    `json:"id" bson:"_id"`
	Name                string       `json:"name" bson:"name"`
	Description         string       `json:"description" bson:"description"`
	IconAsset           string       `json:"icon_asset" bson:"icon_asset"`
	DetailedDescription string       `json:"detailed_description" bson:"detailed_description"`
	Type                EvidenceType `json:"type" bson:"type"`
}
