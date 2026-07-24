package domain

type FinalReport struct {
	Who      string `bson:"who"`
	Why      string `bson:"why"`
	How      string `bson:"how"`
	When     string `bson:"when"`
	Evidence string `bson:"evidence"`
}

type ScoreBreakdown struct {
	WhoCorrect      bool `bson:"who_correct"`
	WhyCorrect      bool `bson:"why_correct"`
	HowCorrect      bool `bson:"how_correct"`
	WhenCorrect     bool `bson:"when_correct"`
	EvidenceCorrect bool `bson:"evidence_correct"`
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
	PlayerReport      FinalReport       `bson:"player_report"`
	Breakdown         ScoreBreakdown    `bson:"breakdown"`
	NarrativeFeedback string            `bson:"narrative_feedback"`
	BreakdownDetails  map[string]string `bson:"breakdown_details"`
	MissedFacts       []string          `bson:"missed_facts"`
}
