import 'report.dart';

enum CrimeType { murder, theft, fraud, arson, kidnapping, blackmail }

enum EvidenceType { physical, digital, document, testimony }

class Crime {
  final String crimeType;
  final String victim;
  final String perpetratorId;
  final String motive;
  final String method;
  final String timeOfCrime;

  const Crime({
    required this.crimeType,
    required this.victim,
    required this.perpetratorId,
    required this.motive,
    required this.method,
    required this.timeOfCrime,
  });

  String get typeLabel {
    switch (CrimeType.values.firstWhere((e) => e.name == crimeType, orElse: () => CrimeType.murder)) {
      case CrimeType.murder:
        return 'убийство';
      case CrimeType.theft:
        return 'кража';
      case CrimeType.fraud:
        return 'мошенничество';
      case CrimeType.arson:
        return 'поджог';
      case CrimeType.kidnapping:
        return 'похищение';
      case CrimeType.blackmail:
        return 'шантаж';
    }
  }

  factory Crime.fromJson(Map<String, dynamic> json) {
    return Crime(
      crimeType: json['crime_type'] as String,
      victim: json['victim'] as String,
      perpetratorId: json['perpetrator_id'] as String,
      motive: json['motive'] as String,
      method: json['method'] as String,
      timeOfCrime: json['time_of_crime'] as String,
    );
  }

  Map<String, dynamic> toJson() => {
    'crime_type': crimeType,
    'victim': victim,
    'perpetrator_id': perpetratorId,
    'motive': motive,
    'method': method,
    'time_of_crime': timeOfCrime,
  };
}

class TimelineEntry {
  final String time;
  final String event;
  final String? characterId;

  const TimelineEntry({required this.time, required this.event, this.characterId});

  factory TimelineEntry.fromJson(Map<String, dynamic> json) {
    return TimelineEntry(
      time: json['time'] as String,
      event: json['event'] as String,
      characterId: json['character_id'] as String?,
    );
  }

  Map<String, dynamic> toJson() => {'time': time, 'event': event, if (characterId != null) 'character_id': characterId};
}

class Timeline {
  final List<TimelineEntry> entries;

  const Timeline({required this.entries});

  factory Timeline.fromJson(List<dynamic> json) {
    return Timeline(entries: json.map((e) => TimelineEntry.fromJson(e as Map<String, dynamic>)).toList());
  }

  List<Map<String, dynamic>> toJson() => entries.map((e) => e.toJson()).toList();
}

class Evidence {
  final String id;
  final String name;
  final String description;
  final String iconAsset;
  final String detailedDescription;
  final String type;

  const Evidence({
    required this.id,
    required this.name,
    required this.description,
    required this.iconAsset,
    required this.detailedDescription,
    required this.type,
  });

  String get typeLabel {
    final t = EvidenceType.values.firstWhere((e) => e.name == type, orElse: () => EvidenceType.physical);
    switch (t) {
      case EvidenceType.physical:
        return 'вещественное';
      case EvidenceType.digital:
        return 'цифровое';
      case EvidenceType.document:
        return 'документ';
      case EvidenceType.testimony:
        return 'показания';
    }
  }

  factory Evidence.fromJson(Map<String, dynamic> json) {
    return Evidence(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String,
      iconAsset: json['icon_asset'] as String,
      detailedDescription: json['detailed_description'] as String,
      type: json['type'] as String,
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'description': description,
    'icon_asset': iconAsset,
    'detailed_description': detailedDescription,
    'type': type,
  };
}

class Session {
  final String id;
  final String crimeType;
  final String victim;
  final String perpetratorId;
  final String motive;
  final String method;
  final String timeOfCrime;
  final Timeline? timeline;
  final String caseName;
  final String caseBrief;
  final int actionPoints;
  final String phase;
  final GameResult? gameResult;
  final DateTime createdAt;

  const Session({
    required this.id,
    this.crimeType = '',
    this.victim = '',
    this.perpetratorId = '',
    this.motive = '',
    this.method = '',
    this.timeOfCrime = '',
    this.timeline,
    this.caseName = '',
    this.caseBrief = '',
    this.actionPoints = 5,
    this.phase = 'idle',
    this.gameResult,
    required this.createdAt,
  });

  Crime get crime => Crime(
    crimeType: crimeType,
    victim: victim,
    perpetratorId: perpetratorId,
    motive: motive,
    method: method,
    timeOfCrime: timeOfCrime,
  );

  factory Session.fromJson(Map<String, dynamic> json) {
    final crimeJson = json['crime'] as Map<String, dynamic>?;
    return Session(
      id: json['id'] as String,
      crimeType: crimeJson?['crime_type'] as String? ?? '',
      victim: crimeJson?['victim'] as String? ?? '',
      perpetratorId: crimeJson?['perpetrator_id'] as String? ?? '',
      motive: crimeJson?['motive'] as String? ?? '',
      method: crimeJson?['method'] as String? ?? '',
      timeOfCrime: crimeJson?['time_of_crime'] as String? ?? '',
      timeline:
          json['timeline'] != null
              ? Timeline.fromJson((json['timeline'] as Map<String, dynamic>)['entries'] as List<dynamic>)
              : null,
      caseName: json['case_name'] as String? ?? '',
      caseBrief: json['case_brief'] as String? ?? '',
      actionPoints: json['action_points'] as int? ?? 5,
      phase: json['phase'] as String? ?? 'idle',
      gameResult: json['game_result'] != null ? GameResult.fromJson(json['game_result'] as Map<String, dynamic>) : null,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'crime': {
      'crime_type': crimeType,
      'victim': victim,
      'perpetrator_id': perpetratorId,
      'motive': motive,
      'method': method,
      'time_of_crime': timeOfCrime,
    },
    if (timeline != null) 'timeline': {'entries': timeline!.toJson()},
    'case_name': caseName,
    'case_brief': caseBrief,
    'action_points': actionPoints,
    'phase': phase,
    if (gameResult != null) 'game_result': gameResult!.toJson(),
    'created_at': createdAt.toIso8601String(),
  };
}
