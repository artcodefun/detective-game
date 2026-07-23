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

func (t CrimeType) Label() string {
	switch t {
	case CrimeTypeMurder:
		return "убийство"
	case CrimeTypeTheft:
		return "кража"
	case CrimeTypeFraud:
		return "мошенничество"
	case CrimeTypeArson:
		return "поджог"
	case CrimeTypeKidnapping:
		return "похищение"
	case CrimeTypeBlackmail:
		return "шантаж"
	}
	return string(t)
}

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

func (t EvidenceType) Label() string {
	switch t {
	case EvidenceTypePhysical:
		return "вещественное"
	case EvidenceTypeDigital:
		return "цифровое"
	case EvidenceTypeDocument:
		return "документ"
	case EvidenceTypeTestimony:
		return "показания"
	}
	return string(t)
}

type Evidence struct {
	ID                  uuid.UUID    `json:"id"`
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	IconAsset           string       `json:"icon_asset"`
	DetailedDescription string       `json:"detailed_description"`
	Type                EvidenceType `json:"type"`
}
