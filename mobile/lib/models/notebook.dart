enum NoteTag { strange, suspicious, lie, key }

enum NotebookEntryType { statement, analysis, alibiCheck, cameraRequest, transactionRequest }

class NotebookEntry {
  final String id;
  final NotebookEntryType type;
  final String? characterId;
  final String description;
  final List<NoteTag> userTags;
  final String? userNote;
  final DateTime timestamp;

  const NotebookEntry({
    required this.id,
    required this.type,
    this.characterId,
    required this.description,
    this.userTags = const [],
    this.userNote,
    required this.timestamp,
  });

  String get typeLabel {
    switch (type) {
      case NotebookEntryType.statement:
        return 'показания';
      case NotebookEntryType.analysis:
        return 'анализ';
      case NotebookEntryType.alibiCheck:
        return 'проверка алиби';
      case NotebookEntryType.cameraRequest:
        return 'запись с камер';
      case NotebookEntryType.transactionRequest:
        return 'транзакция';
    }
  }

  String get tagsLabel {
    if (userTags.isEmpty) return '';
    return userTags
        .map((t) {
          switch (t) {
            case NoteTag.strange:
              return 'странно';
            case NoteTag.suspicious:
              return 'подозрительно';
            case NoteTag.lie:
              return 'ложь';
            case NoteTag.key:
              return 'ключевое';
          }
        })
        .join(', ');
  }

  NotebookEntry copyWith({
    String? id,
    NotebookEntryType? type,
    String? characterId,
    String? description,
    List<NoteTag>? userTags,
    String? userNote,
    bool clearTags = false,
    bool clearNote = false,
    DateTime? timestamp,
  }) {
    return NotebookEntry(
      id: id ?? this.id,
      type: type ?? this.type,
      characterId: characterId ?? this.characterId,
      description: description ?? this.description,
      userTags: clearTags ? [] : userTags ?? this.userTags,
      userNote: clearNote ? null : userNote ?? this.userNote,
      timestamp: timestamp ?? this.timestamp,
    );
  }

  static NotebookEntryType entryTypeFromString(String value) {
    switch (value) {
      case 'analysis':
        return NotebookEntryType.analysis;
      case 'alibi_check':
        return NotebookEntryType.alibiCheck;
      case 'camera_request':
        return NotebookEntryType.cameraRequest;
      case 'transaction_request':
        return NotebookEntryType.transactionRequest;
      default:
        return NotebookEntryType.statement;
    }
  }

  factory NotebookEntry.fromJson(Map<String, dynamic> json) {
    return NotebookEntry(
      id: json['id'] as String,
      type: entryTypeFromString(json['type'] as String),
      characterId: json['character_id'] as String?,
      description: json['description'] as String,
      userTags:
          json['user_tags'] != null
              ? (json['user_tags'] as List<dynamic>).map((e) => NoteTag.values.byName(e as String)).toList()
              : [],
      userNote: json['user_note'] as String?,
      timestamp: DateTime.parse(json['timestamp'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'type': type.name,
      'character_id': characterId,
      'description': description,
      'user_tags': userTags.map((t) => t.name).toList(),
      'user_note': userNote,
      'timestamp': timestamp.toIso8601String(),
    };
  }
}
