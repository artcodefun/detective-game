package domain

type FinalReport struct {
	Who      string `json:"who" bson:"who"`
	Why      string `json:"why" bson:"why"`
	How      string `json:"how" bson:"how"`
	When     string `json:"when" bson:"when"`
	Evidence string `json:"evidence" bson:"evidence"`
}

type ScoreBreakdown struct {
	WhoCorrect      bool `json:"who_correct" bson:"who_correct"`
	WhyCorrect      bool `json:"why_correct" bson:"why_correct"`
	HowCorrect      bool `json:"how_correct" bson:"how_correct"`
	WhenCorrect     bool `json:"when_correct" bson:"when_correct"`
	EvidenceCorrect bool `json:"evidence_correct" bson:"evidence_correct"`
}

func (s ScoreBreakdown) CorrectCount() int {
	count := 0
	if s.WhoCorrect {
		count++
	}
	if s.WhyCorrect {
		count++
	}
	if s.HowCorrect {
		count++
	}
	if s.WhenCorrect {
		count++
	}
	if s.EvidenceCorrect {
		count++
	}
	return count
}

func (s ScoreBreakdown) Accuracy() float64 {
	return float64(s.CorrectCount()) / 5.0
}

type GameResult struct {
	PlayerReport      FinalReport       `json:"player_report" bson:"player_report"`
	Breakdown         ScoreBreakdown    `json:"breakdown" bson:"breakdown"`
	NarrativeFeedback string            `json:"narrative_feedback" bson:"narrative_feedback"`
	BreakdownDetails  map[string]string `json:"breakdown_details" bson:"breakdown_details"`
	MissedFacts       []string          `json:"missed_facts" bson:"missed_facts"`
}
