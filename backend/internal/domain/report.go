package domain

type FinalReport struct {
	Who      string `json:"who"`
	Why      string `json:"why"`
	How      string `json:"how"`
	When     string `json:"when"`
	Evidence string `json:"evidence"`
}

type ScoreBreakdown struct {
	WhoCorrect      bool `json:"who_correct"`
	WhyCorrect      bool `json:"why_correct"`
	HowCorrect      bool `json:"how_correct"`
	WhenCorrect     bool `json:"when_correct"`
	EvidenceCorrect bool `json:"evidence_correct"`
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
	PlayerReport      FinalReport       `json:"player_report"`
	Breakdown         ScoreBreakdown    `json:"breakdown"`
	NarrativeFeedback string            `json:"narrative_feedback"`
	BreakdownDetails  map[string]string `json:"breakdown_details"`
	MissedFacts       []string          `json:"missed_facts"`
}
