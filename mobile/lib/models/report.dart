class FinalReport {
  final String who;
  final String why;
  final String how;
  final String when;
  final String evidence;

  const FinalReport({
    required this.who,
    required this.why,
    required this.how,
    required this.when,
    required this.evidence,
  });

  factory FinalReport.fromJson(Map<String, dynamic> json) {
    return FinalReport(
      who: json['who'] as String,
      why: json['why'] as String,
      how: json['how'] as String,
      when: json['when'] as String,
      evidence: json['evidence'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {'who': who, 'why': why, 'how': how, 'when': when, 'evidence': evidence};
  }
}

class ScoreBreakdown {
  final bool whoCorrect;
  final bool whyCorrect;
  final bool howCorrect;
  final bool whenCorrect;
  final bool evidenceCorrect;

  const ScoreBreakdown({
    required this.whoCorrect,
    required this.whyCorrect,
    required this.howCorrect,
    required this.whenCorrect,
    required this.evidenceCorrect,
  });

  int get correctCount {
    return [whoCorrect, whyCorrect, howCorrect, whenCorrect, evidenceCorrect].where((c) => c).length;
  }

  int get totalCount => 5;

  double get accuracy => correctCount / totalCount;

  factory ScoreBreakdown.fromJson(Map<String, dynamic> json) {
    return ScoreBreakdown(
      whoCorrect: json['who_correct'] as bool,
      whyCorrect: json['why_correct'] as bool,
      howCorrect: json['how_correct'] as bool,
      whenCorrect: json['when_correct'] as bool,
      evidenceCorrect: json['evidence_correct'] as bool,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'who_correct': whoCorrect,
      'why_correct': whyCorrect,
      'how_correct': howCorrect,
      'when_correct': whenCorrect,
      'evidence_correct': evidenceCorrect,
    };
  }
}

class ScoreBreakdownDetails {
  final String who;
  final String why;
  final String how;
  final String when;
  final String evidence;

  const ScoreBreakdownDetails({
    required this.who,
    required this.why,
    required this.how,
    required this.when,
    required this.evidence,
  });

  factory ScoreBreakdownDetails.fromJson(Map<String, dynamic> json) {
    return ScoreBreakdownDetails(
      who: json['who'] as String,
      why: json['why'] as String,
      how: json['how'] as String,
      when: json['when'] as String,
      evidence: json['evidence'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {'who': who, 'why': why, 'how': how, 'when': when, 'evidence': evidence};
  }
}

class GameResult {
  final FinalReport playerReport;
  final ScoreBreakdown breakdown;
  final String narrativeFeedback;
  final ScoreBreakdownDetails breakdownDetails;
  final List<String> missedFacts;

  const GameResult({
    required this.playerReport,
    required this.breakdown,
    required this.narrativeFeedback,
    required this.breakdownDetails,
    this.missedFacts = const [],
  });

  factory GameResult.fromJson(Map<String, dynamic> json) {
    return GameResult(
      playerReport: FinalReport.fromJson(json['player_report'] as Map<String, dynamic>),
      breakdown: ScoreBreakdown.fromJson(json['breakdown'] as Map<String, dynamic>),
      narrativeFeedback: json['narrative_feedback'] as String,
      breakdownDetails: ScoreBreakdownDetails.fromJson(json['breakdown_details'] as Map<String, dynamic>),
      missedFacts: List<String>.from(json['missed_facts'] ?? []),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'player_report': playerReport.toJson(),
      'breakdown': breakdown.toJson(),
      'narrative_feedback': narrativeFeedback,
      'breakdown_details': breakdownDetails,
      'missed_facts': missedFacts,
    };
  }
}
