class Interrogation {
  final String id;
  final String sessionId;
  final String characterId;
  final String phase;
  final DateTime createdAt;
  final DateTime? completedAt;

  const Interrogation({
    required this.id,
    required this.sessionId,
    required this.characterId,
    required this.phase,
    required this.createdAt,
    this.completedAt,
  });

  bool get isActive => phase == 'active';

  factory Interrogation.fromJson(Map<String, dynamic> json) {
    return Interrogation(
      id: json['id'] as String,
      sessionId: json['session_id'] as String,
      characterId: json['character_id'] as String,
      phase: json['phase'] as String,
      createdAt: DateTime.parse(json['created_at'] as String),
      completedAt: json['completed_at'] != null ? DateTime.parse(json['completed_at'] as String) : null,
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'session_id': sessionId,
    'character_id': characterId,
    'phase': phase,
    'created_at': createdAt.toIso8601String(),
    if (completedAt != null) 'completed_at': completedAt!.toIso8601String(),
  };
}
