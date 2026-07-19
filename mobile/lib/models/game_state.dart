import 'character.dart';
import 'scenario.dart';
import 'chronology_entry.dart';
import 'notebook.dart';

enum GamePhase { idle, generating, investigating, writingReport, finished }

enum TrustLevel { open, reserved, tense, closed }

class Memory {
  final String id;
  final String content;
  final bool isTrue;
  final String timestamp;

  const Memory({
    required this.id,
    required this.content,
    required this.isTrue,
    required this.timestamp,
  });

  factory Memory.fromJson(Map<String, dynamic> json) {
    return Memory(
      id: json['id'] as String,
      content: json['content'] as String,
      isTrue: json['is_true'] as bool,
      timestamp: json['timestamp'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'content': content,
      'is_true': isTrue,
      'timestamp': timestamp,
    };
  }
}

class CharacterKnowledge {
  final List<String> knownFacts;
  final List<String> partialFacts;
  final List<String> falseBeliefs;

  const CharacterKnowledge({
    this.knownFacts = const [],
    this.partialFacts = const [],
    this.falseBeliefs = const [],
  });

  factory CharacterKnowledge.fromJson(Map<String, dynamic> json) {
    return CharacterKnowledge(
      knownFacts: List<String>.from(json['known_facts'] ?? []),
      partialFacts: List<String>.from(json['partial_facts'] ?? []),
      falseBeliefs: List<String>.from(json['false_beliefs'] ?? []),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'known_facts': knownFacts,
      'partial_facts': partialFacts,
      'false_beliefs': falseBeliefs,
    };
  }
}

class InterrogationMessage {
  final String sender;
  final String text;
  final DateTime timestamp;

  const InterrogationMessage({
    required this.sender,
    required this.text,
    required this.timestamp,
  });

  factory InterrogationMessage.fromJson(Map<String, dynamic> json) {
    return InterrogationMessage(
      sender: json['sender'] as String,
      text: json['text'] as String,
      timestamp: DateTime.parse(json['timestamp'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'sender': sender,
      'text': text,
      'timestamp': timestamp.toIso8601String(),
    };
  }
}

class CharacterState {
  final CharacterData base;
  final CharacterKnowledge knowledge;
  final List<String> secrets;
  final Map<String, String> relationships;
  final List<Memory> memories;
  final int trust;
  final int interrogationsRemaining;
  final List<InterrogationMessage> chatHistory;

  static const int maxTrust = 100;
  static const int minTrust = 0;
  static const int maxInterrogations = 3;

  const CharacterState({
    required this.base,
    required this.knowledge,
    this.secrets = const [],
    this.relationships = const {},
    this.memories = const [],
    this.trust = 50,
    this.interrogationsRemaining = maxInterrogations,
    this.chatHistory = const [],
  });

  TrustLevel get trustLevel {
    if (trust >= 75) return TrustLevel.open;
    if (trust >= 50) return TrustLevel.reserved;
    if (trust >= 25) return TrustLevel.tense;
    return TrustLevel.closed;
  }

  bool get canInterrogate => interrogationsRemaining > 0;

  CharacterState applyAttitudeDelta(int delta) {
    final newTrust = (trust + delta).clamp(minTrust, maxTrust);
    return copyWith(trust: newTrust);
  }

  CharacterState decrementInterrogation() {
    return copyWith(
      interrogationsRemaining: interrogationsRemaining - 1,
    );
  }

  CharacterState addMessage(InterrogationMessage message) {
    final updatedHistory = List<InterrogationMessage>.from(chatHistory)
      ..add(message);
    return copyWith(chatHistory: updatedHistory);
  }

  CharacterState copyWith({
    CharacterData? base,
    CharacterKnowledge? knowledge,
    List<String>? secrets,
    Map<String, String>? relationships,
    List<Memory>? memories,
    int? trust,
    int? interrogationsRemaining,
    List<InterrogationMessage>? chatHistory,
  }) {
    return CharacterState(
      base: base ?? this.base,
      knowledge: knowledge ?? this.knowledge,
      secrets: secrets ?? this.secrets,
      relationships: relationships ?? this.relationships,
      memories: memories ?? this.memories,
      trust: trust ?? this.trust,
      interrogationsRemaining:
          interrogationsRemaining ?? this.interrogationsRemaining,
      chatHistory: chatHistory ?? this.chatHistory,
    );
  }

  factory CharacterState.fromJson(Map<String, dynamic> json) {
    return CharacterState(
      base: CharacterData.fromJson(json['base'] as Map<String, dynamic>),
      knowledge: CharacterKnowledge.fromJson(
          json['knowledge'] as Map<String, dynamic>),
      secrets: List<String>.from(json['secrets'] ?? []),
      relationships:
          Map<String, String>.from(json['relationships'] ?? {}),
      memories: (json['memories'] as List<dynamic>?)
              ?.map((e) => Memory.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      trust: json['trust'] as int? ?? 50,
      interrogationsRemaining:
          json['interrogations_remaining'] as int? ?? maxInterrogations,
      chatHistory: (json['chat_history'] as List<dynamic>?)
              ?.map((e) =>
                  InterrogationMessage.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'base': base.toJson(),
      'knowledge': knowledge.toJson(),
      'secrets': secrets,
      'relationships': relationships,
      'memories': memories.map((m) => m.toJson()).toList(),
      'trust': trust,
      'interrogations_remaining': interrogationsRemaining,
      'chat_history': chatHistory.map((m) => m.toJson()).toList(),
    };
  }
}

class GameSession {
  final String id;
  final Crime crime;
  final Timeline timeline;
  final List<Evidence> evidence;
  final List<CharacterState> characters;
  final List<ChronologyEntry> chronology;
  final int actionPoints;
  final GamePhase phase;
  final DateTime createdAt;

  static const int maxActionPoints = 5;

  GameSession({
    required this.id,
    required this.crime,
    required this.timeline,
    required this.evidence,
    required this.characters,
    this.chronology = const [],
    this.actionPoints = maxActionPoints,
    this.phase = GamePhase.idle,
    DateTime? createdAt,
  }) : createdAt = createdAt ?? DateTime.now();

  bool get canSpendActionPoint => actionPoints > 0;

  CharacterState? characterById(String id) {
    for (final c in characters) {
      if (c.base.id == id) return c;
    }
    return null;
  }

  GameSession updateCharacter(String id, CharacterState updated) {
    final idx = characters.indexWhere((c) => c.base.id == id);
    if (idx == -1) return this;
    final newChars = List<CharacterState>.from(characters);
    newChars[idx] = updated;
    return copyWith(characters: newChars);
  }

  GameSession spendActionPoint() {
    return copyWith(actionPoints: actionPoints - 1);
  }

  GameSession spendActionPoints(int amount) {
    return copyWith(actionPoints: actionPoints - amount);
  }

  GameSession addChronologyEntry(ChronologyEntry entry) {
    final updated = List<ChronologyEntry>.from(chronology)..add(entry);
    return copyWith(chronology: updated);
  }

  GameSession addDetailsToChronology(String chronologyId, List<NotebookEntry> details) {
    if (details.isEmpty) return this;
    final idx = chronology.indexWhere((c) => c.id == chronologyId);
    if (idx == -1) return this;
    final updated = List<ChronologyEntry>.from(chronology);
    updated[idx] = updated[idx].copyWith(
      details: [...updated[idx].details, ...details],
    );
    return copyWith(chronology: updated);
  }

  GameSession updateNotebookEntry({
    required String chronologyId,
    required String entryId,
    List<NoteTag>? userTags,
    String? userNote,
    bool clearTags = false,
    bool clearNote = false,
  }) {
    final chronIdx = chronology.indexWhere((c) => c.id == chronologyId);
    if (chronIdx == -1) return this;
    final chron = chronology[chronIdx];
    final detailIdx = chron.details.indexWhere((d) => d.id == entryId);
    if (detailIdx == -1) return this;
    final updatedDetails = List<NotebookEntry>.from(chron.details);
    updatedDetails[detailIdx] = updatedDetails[detailIdx].copyWith(
      userTags: userTags,
      clearTags: clearTags,
      userNote: userNote,
      clearNote: clearNote,
    );
    final updatedChronology = List<ChronologyEntry>.from(chronology);
    updatedChronology[chronIdx] = chron.copyWith(details: updatedDetails);
    return copyWith(chronology: updatedChronology);
  }

  GameSession copyWith({
    String? id,
    Crime? crime,
    Timeline? timeline,
    List<Evidence>? evidence,
    List<CharacterState>? characters,
    List<ChronologyEntry>? chronology,
    int? actionPoints,
    GamePhase? phase,
    DateTime? createdAt,
  }) {
    return GameSession(
      id: id ?? this.id,
      crime: crime ?? this.crime,
      timeline: timeline ?? this.timeline,
      evidence: evidence ?? this.evidence,
      characters: characters ?? this.characters,
      chronology: chronology ?? this.chronology,
      actionPoints: actionPoints ?? this.actionPoints,
      phase: phase ?? this.phase,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  factory GameSession.fromJson(Map<String, dynamic> json) {
    return GameSession(
      id: json['id'] as String,
      crime: Crime.fromJson(json['crime'] as Map<String, dynamic>),
      timeline: Timeline.fromJson(json['timeline'] as List<dynamic>),
      evidence: (json['evidence'] as List<dynamic>)
          .map((e) => Evidence.fromJson(e as Map<String, dynamic>))
          .toList(),
      characters: (json['characters'] as List<dynamic>)
          .map((e) => CharacterState.fromJson(e as Map<String, dynamic>))
          .toList(),
      chronology: (json['chronology'] as List<dynamic>?)
              ?.map(
                  (e) => ChronologyEntry.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      actionPoints: json['action_points'] as int? ?? maxActionPoints,
      phase: GamePhase.values.byName(json['phase'] as String),
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'crime': crime.toJson(),
      'timeline': timeline.toJson(),
      'evidence': evidence.map((e) => e.toJson()).toList(),
      'characters': characters.map((c) => c.toJson()).toList(),
      'chronology': chronology.map((e) => e.toJson()).toList(),
      'action_points': actionPoints,
      'phase': phase.name,
      'created_at': createdAt.toIso8601String(),
    };
  }
}
