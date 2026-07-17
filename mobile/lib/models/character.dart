class CharacterData {
  final String id;
  final String name;
  final int age;
  final String profession;
  final String imagePath;
  final String personality;
  final String audioToneId;

  const CharacterData({
    required this.id,
    required this.name,
    required this.age,
    required this.profession,
    required this.imagePath,
    required this.personality,
    required this.audioToneId,
  });

  CharacterData copyWith({
    String? id,
    String? name,
    int? age,
    String? profession,
    String? imagePath,
    String? personality,
    String? audioToneId,
  }) {
    return CharacterData(
      id: id ?? this.id,
      name: name ?? this.name,
      age: age ?? this.age,
      profession: profession ?? this.profession,
      imagePath: imagePath ?? this.imagePath,
      personality: personality ?? this.personality,
      audioToneId: audioToneId ?? this.audioToneId,
    );
  }

  factory CharacterData.fromJson(Map<String, dynamic> json) {
    return CharacterData(
      id: json['id'] as String,
      name: json['name'] as String,
      age: json['age'] as int,
      profession: json['profession'] as String,
      imagePath: json['image_path'] as String,
      personality: json['personality'] as String,
      audioToneId: json['audio_tone_id'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'age': age,
      'profession': profession,
      'image_path': imagePath,
      'personality': personality,
      'audio_tone_id': audioToneId,
    };
  }
}

class CharacterPool {
  final List<CharacterData> characters;

  const CharacterPool({required this.characters});

  List<CharacterData> getRandom(int count) {
    if (count >= characters.length) return List.from(characters);
    final shuffled = List<CharacterData>.from(characters)..shuffle();
    return shuffled.take(count).toList();
  }

  CharacterData? byId(String id) {
    for (final c in characters) {
      if (c.id == id) return c;
    }
    return null;
  }

  factory CharacterPool.fromJson(List<dynamic> json) {
    return CharacterPool(characters: json.map((e) => CharacterData.fromJson(e as Map<String, dynamic>)).toList());
  }

  List<Map<String, dynamic>> toJson() {
    return characters.map((c) => c.toJson()).toList();
  }
}
