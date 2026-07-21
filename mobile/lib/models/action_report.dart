class ActionReport {
  final String id;
  final String title;
  final String description;
  final String body;
  final String? evidenceId;
  final String? characterId;
  final DateTime timestamp;

  const ActionReport({
    required this.id,
    required this.title,
    required this.description,
    required this.body,
    this.evidenceId,
    this.characterId,
    required this.timestamp,
  });

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'body': body,
      'evidence_id': evidenceId,
      'character_id': characterId,
      'timestamp': timestamp.toIso8601String(),
    };
  }
}
