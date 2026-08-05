import 'notebook.dart';

enum ChronologyEventType { caseStarted, interrogation, labAnalysis, alibiCheck, cameraReview, transactionCheck }

class ChronologyEntry {
  final String id;
  final ChronologyEventType eventType;
  final String title;
  final DateTime timestamp;
  final List<NotebookEntry> details;

  const ChronologyEntry({
    required this.id,
    required this.eventType,
    required this.title,
    required this.timestamp,
    this.details = const [],
  });

  String get eventTypeLabel {
    switch (eventType) {
      case ChronologyEventType.caseStarted:
        return 'начало дела';
      case ChronologyEventType.interrogation:
        return 'допрос';
      case ChronologyEventType.labAnalysis:
        return 'экспертиза';
      case ChronologyEventType.alibiCheck:
        return 'проверка алиби';
      case ChronologyEventType.cameraReview:
        return 'запись с камер';
      case ChronologyEventType.transactionCheck:
        return 'транзакция';
    }
  }

  ChronologyEntry copyWith({
    String? id,
    ChronologyEventType? eventType,
    String? title,
    DateTime? timestamp,
    List<NotebookEntry>? details,
  }) {
    return ChronologyEntry(
      id: id ?? this.id,
      eventType: eventType ?? this.eventType,
      title: title ?? this.title,
      timestamp: timestamp ?? this.timestamp,
      details: details ?? this.details,
    );
  }

  static ChronologyEventType _eventTypeFromString(String value) {
    switch (value) {
      case 'case_started':
        return ChronologyEventType.caseStarted;
      case 'lab_analysis':
        return ChronologyEventType.labAnalysis;
      case 'alibi_check':
        return ChronologyEventType.alibiCheck;
      case 'camera_review':
        return ChronologyEventType.cameraReview;
      case 'transaction_check':
        return ChronologyEventType.transactionCheck;
      default:
        return ChronologyEventType.interrogation;
    }
  }

  factory ChronologyEntry.fromJson(Map<String, dynamic> json) {
    return ChronologyEntry(
      id: json['id'] as String,
      eventType: _eventTypeFromString(json['event_type'] as String),
      title: json['title'] as String,
      timestamp: DateTime.parse(json['timestamp'] as String),
      details:
          (json['details'] as List<dynamic>?)?.map((e) => NotebookEntry.fromJson(e as Map<String, dynamic>)).toList() ??
          [],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'event_type': eventType.name,
      'title': title,
      'timestamp': timestamp.toIso8601String(),
      'details': details.map((e) => e.toJson()).toList(),
    };
  }
}
