class ActionReport {
  final String id;
  final String? type;
  final String title;
  final String description;
  final String body;
  final String? evidenceId;
  final String? characterId;
  final DateTime timestamp;

  const ActionReport({
    required this.id,
    this.type,
    required this.title,
    required this.description,
    required this.body,
    this.evidenceId,
    this.characterId,
    required this.timestamp,
  });

  factory ActionReport.fromJson(Map<String, dynamic> json) {
    return ActionReport(
      id: json['id'] as String,
      type: json['type'] as String?,
      title: json['title'] as String,
      description: json['description'] as String,
      body: json['body'] as String,
      evidenceId: json['evidence_id'] as String?,
      characterId: json['character_id'] as String?,
      timestamp: DateTime.parse(json['timestamp'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      if (type != null) 'type': type,
      'title': title,
      'description': description,
      'body': body,
      if (evidenceId != null) 'evidence_id': evidenceId,
      if (characterId != null) 'character_id': characterId,
      'timestamp': timestamp.toIso8601String(),
    };
  }
}
