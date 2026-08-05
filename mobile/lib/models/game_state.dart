enum GamePhase { idle, generating, investigating, writingReport, finished }

enum TrustLevel { open, reserved, tense, closed }

class Memory {
  final String id;
  final String content;
  final bool isTrue;
  final String timestamp;

  const Memory({required this.id, required this.content, required this.isTrue, required this.timestamp});

  factory Memory.fromJson(Map<String, dynamic> json) {
    return Memory(
      id: json['id'] as String,
      content: json['content'] as String,
      isTrue: json['is_true'] as bool,
      timestamp: json['timestamp'] as String,
    );
  }

  Map<String, dynamic> toJson() => {'id': id, 'content': content, 'is_true': isTrue, 'timestamp': timestamp};
}

class CharacterKnowledge {
  final List<String> knownFacts;
  final List<String> partialFacts;
  final List<String> falseBeliefs;

  const CharacterKnowledge({this.knownFacts = const [], this.partialFacts = const [], this.falseBeliefs = const []});

  factory CharacterKnowledge.fromJson(Map<String, dynamic> json) {
    return CharacterKnowledge(
      knownFacts: List<String>.from(json['known_facts'] ?? []),
      partialFacts: List<String>.from(json['partial_facts'] ?? []),
      falseBeliefs: List<String>.from(json['false_beliefs'] ?? []),
    );
  }

  Map<String, dynamic> toJson() => {
    'known_facts': knownFacts,
    'partial_facts': partialFacts,
    'false_beliefs': falseBeliefs,
  };
}

class ChatMessage {
  final String id;
  final String sessionId;
  final String interrogationId;
  final bool fromUser;
  final String text;
  final List<String> statements;
  final int attitudeDelta;
  final DateTime timestamp;

  const ChatMessage({
    required this.id,
    required this.sessionId,
    required this.interrogationId,
    required this.fromUser,
    required this.text,
    this.statements = const [],
    this.attitudeDelta = 0,
    required this.timestamp,
  });

  String get sender => fromUser ? 'Вы' : '';

  factory ChatMessage.fromJson(Map<String, dynamic> json) {
    return ChatMessage(
      id: json['id'] as String,
      sessionId: json['session_id'] as String,
      interrogationId: json['interrogation_id'] as String,
      fromUser: json['from_user'] as bool,
      text: json['text'] as String,
      statements: List<String>.from(json['statements'] ?? []),
      attitudeDelta: json['attitude_delta'] as int? ?? 0,
      timestamp: DateTime.parse(json['timestamp'] as String),
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'session_id': sessionId,
    'interrogation_id': interrogationId,
    'from_user': fromUser,
    'text': text,
    'statements': statements,
    'attitude_delta': attitudeDelta,
    'timestamp': timestamp.toIso8601String(),
  };
}

class Character {
  final String id;
  final String sessionId;
  final String name;
  final int age;
  final String profession;
  final String personality;
  final String? gender;
  final CharacterKnowledge knowledge;
  final List<String> secrets;
  final Map<String, String> relationships;
  final List<Memory> memories;
  final int trust;
  final int interrogationsRemaining;

  static const int maxTrust = 100;
  static const int minTrust = 0;
  static const int maxInterrogations = 3;

  const Character({
    required this.id,
    required this.sessionId,
    required this.name,
    required this.age,
    required this.profession,
    required this.personality,
    this.gender,
    this.knowledge = const CharacterKnowledge(),
    this.secrets = const [],
    this.relationships = const {},
    this.memories = const [],
    this.trust = 50,
    this.interrogationsRemaining = maxInterrogations,
  });

  TrustLevel get trustLevel {
    if (trust >= 75) return TrustLevel.open;
    if (trust >= 50) return TrustLevel.reserved;
    if (trust >= 25) return TrustLevel.tense;
    return TrustLevel.closed;
  }

  bool get canInterrogate => interrogationsRemaining > 0;

  Character copyWith({
    String? id,
    String? sessionId,
    String? name,
    int? age,
    String? profession,
    String? personality,
    String? gender,
    CharacterKnowledge? knowledge,
    List<String>? secrets,
    Map<String, String>? relationships,
    List<Memory>? memories,
    int? trust,
    int? interrogationsRemaining,
  }) {
    return Character(
      id: id ?? this.id,
      sessionId: sessionId ?? this.sessionId,
      name: name ?? this.name,
      age: age ?? this.age,
      profession: profession ?? this.profession,
      personality: personality ?? this.personality,
      gender: gender ?? this.gender,
      knowledge: knowledge ?? this.knowledge,
      secrets: secrets ?? this.secrets,
      relationships: relationships ?? this.relationships,
      memories: memories ?? this.memories,
      trust: trust ?? this.trust,
      interrogationsRemaining: interrogationsRemaining ?? this.interrogationsRemaining,
    );
  }

  factory Character.fromJson(Map<String, dynamic> json) {
    return Character(
      id: json['id'] as String,
      sessionId: json['session_id'] as String,
      name: json['name'] as String,
      age: json['age'] as int,
      profession: json['profession'] as String,
      personality: json['personality'] as String,
      gender: json['gender'] as String?,
      knowledge:
          json['knowledge'] != null
              ? CharacterKnowledge.fromJson(json['knowledge'] as Map<String, dynamic>)
              : const CharacterKnowledge(),
      secrets: List<String>.from(json['secrets'] ?? []),
      relationships: Map<String, String>.from(json['relationships'] ?? {}),
      memories:
          (json['memories'] as List<dynamic>?)?.map((e) => Memory.fromJson(e as Map<String, dynamic>)).toList() ?? [],
      trust: json['trust'] as int? ?? 50,
      interrogationsRemaining: json['interrogations_remaining'] as int? ?? maxInterrogations,
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'session_id': sessionId,
    'name': name,
    'age': age,
    'profession': profession,
    'personality': personality,
    if (gender != null) 'gender': gender,
    'knowledge': knowledge.toJson(),
    'secrets': secrets,
    'relationships': relationships,
    'memories': memories.map((m) => m.toJson()).toList(),
    'trust': trust,
    'interrogations_remaining': interrogationsRemaining,
  };
}
