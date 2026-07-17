enum CrimeType { murder, theft, fraud, arson, kidnapping, blackmail }

enum EvidenceType { physical, digital, document, testimony }

class Crime {
  final CrimeType type;
  final String victim;
  final String perpetratorId;
  final String motive;
  final String method;
  final String timeOfCrime;

  const Crime({
    required this.type,
    required this.victim,
    required this.perpetratorId,
    required this.motive,
    required this.method,
    required this.timeOfCrime,
  });

  String get typeLabel {
    switch (type) {
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

  static CrimeType crimeTypeFromString(String value) {
    switch (value) {
      case 'theft':
        return CrimeType.theft;
      case 'fraud':
        return CrimeType.fraud;
      case 'arson':
        return CrimeType.arson;
      case 'kidnapping':
        return CrimeType.kidnapping;
      case 'blackmail':
        return CrimeType.blackmail;
      default:
        return CrimeType.murder;
    }
  }

  factory Crime.fromJson(Map<String, dynamic> json) {
    return Crime(
      type: crimeTypeFromString(json['crime_type'] as String),
      victim: json['victim'] as String,
      perpetratorId: json['perpetrator_id'] as String,
      motive: json['motive'] as String,
      method: json['method'] as String,
      timeOfCrime: json['time_of_crime'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'crime_type': type.name,
      'victim': victim,
      'perpetrator_id': perpetratorId,
      'motive': motive,
      'method': method,
      'time_of_crime': timeOfCrime,
    };
  }
}

class TimelineEntry {
  final String time;
  final String event;
  final String? characterId;

  const TimelineEntry({
    required this.time,
    required this.event,
    this.characterId,
  });

  factory TimelineEntry.fromJson(Map<String, dynamic> json) {
    return TimelineEntry(
      time: json['time'] as String,
      event: json['event'] as String,
      characterId: json['character_id'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'time': time,
      'event': event,
      'character_id': characterId,
    };
  }
}

class Timeline {
  final List<TimelineEntry> entries;

  const Timeline({required this.entries});

  factory Timeline.fromJson(List<dynamic> json) {
    return Timeline(
      entries: json
          .map((e) => TimelineEntry.fromJson(e as Map<String, dynamic>))
          .toList(),
    );
  }

  List<Map<String, dynamic>> toJson() {
    return entries.map((e) => e.toJson()).toList();
  }
}

class Evidence {
  final String id;
  final String name;
  final String description;
  final String iconAsset;
  final String detailedDescription;
  final EvidenceType type;
  final bool analyzed;
  final String? analysisResult;

  const Evidence({
    required this.id,
    required this.name,
    required this.description,
    required this.iconAsset,
    required this.detailedDescription,
    required this.type,
    this.analyzed = false,
    this.analysisResult,
  });

  Evidence copyWith({
    String? id,
    String? name,
    String? description,
    String? iconAsset,
    String? detailedDescription,
    EvidenceType? type,
    bool? analyzed,
    String? analysisResult,
    bool clearAnalysisResult = false,
  }) {
    return Evidence(
      id: id ?? this.id,
      name: name ?? this.name,
      description: description ?? this.description,
      iconAsset: iconAsset ?? this.iconAsset,
      detailedDescription: detailedDescription ?? this.detailedDescription,
      type: type ?? this.type,
      analyzed: analyzed ?? this.analyzed,
      analysisResult:
          clearAnalysisResult ? null : analysisResult ?? this.analysisResult,
    );
  }

  String get typeLabel {
    switch (type) {
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

  static EvidenceType evidenceTypeFromString(String value) {
    switch (value) {
      case 'digital':
        return EvidenceType.digital;
      case 'document':
        return EvidenceType.document;
      case 'testimony':
        return EvidenceType.testimony;
      default:
        return EvidenceType.physical;
    }
  }

  factory Evidence.fromJson(Map<String, dynamic> json) {
    return Evidence(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String,
      iconAsset: json['icon_asset'] as String,
      detailedDescription: json['detailed_description'] as String,
      type: evidenceTypeFromString(json['type'] as String),
      analyzed: json['analyzed'] as bool? ?? false,
      analysisResult: json['analysis_result'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'description': description,
      'icon_asset': iconAsset,
      'detailed_description': detailedDescription,
      'type': type.name,
      'analyzed': analyzed,
      'analysis_result': analysisResult,
    };
  }
}
