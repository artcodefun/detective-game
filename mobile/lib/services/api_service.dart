import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:uuid/uuid.dart';

import '../models/action_report.dart';
import '../models/chronology_entry.dart';
import '../models/game_state.dart';
import '../models/interrogation.dart';
import '../models/report.dart';
import '../models/scenario.dart';

class ApiException implements Exception {
  final int statusCode;
  final String message;

  const ApiException(this.statusCode, this.message);

  @override
  String toString() => 'ApiException($statusCode): $message';
}

class ApiService {
  static const _userIdKey = 'detective_user_id';

  final String baseUrl;
  final http.Client _client;
  String _userId = '';
  String? _sessionId;

  String get userId => _userId;

  String? get sessionId => _sessionId;

  ApiService({required this.baseUrl, http.Client? client})
    : _client = client ?? http.Client();

  static Future<String> loadOrCreateUserId() async {
    final prefs = await SharedPreferences.getInstance();
    final existing = prefs.getString(_userIdKey);
    if (existing != null) return existing;
    final newId = const Uuid().v4();
    await prefs.setString(_userIdKey, newId);
    return newId;
  }

  Future<void> init() async {
    _userId = await loadOrCreateUserId();
  }

  void setSessionId(String id) => _sessionId = id;

  void clearSessionId() => _sessionId = null;

  Map<String, String> get _headers => {
    'Content-Type': 'application/json',
    'Accept-Language': 'ru',
    'X-User-ID': _userId,
    if (_sessionId != null) 'X-Session-ID': _sessionId!,
  };

  Future<Map<String, dynamic>> _get(String path) async {
    final uri = Uri.parse('$baseUrl$path');
    final res = await _client.get(uri, headers: _headers);
    return _handleResponse(res);
  }

  Future<Map<String, dynamic>> _post(
    String path, {
    Map<String, dynamic>? body,
  }) async {
    final uri = Uri.parse('$baseUrl$path');
    final res = await _client.post(
      uri,
      headers: _headers,
      body: body != null ? jsonEncode(body) : null,
    );
    return _handleResponse(res);
  }

  Future<Map<String, dynamic>> _patch(
    String path, {
    Map<String, dynamic>? body,
  }) async {
    final uri = Uri.parse('$baseUrl$path');
    final res = await _client.patch(
      uri,
      headers: _headers,
      body: body != null ? jsonEncode(body) : null,
    );
    return _handleResponse(res);
  }

  Future<List<dynamic>> _getList(String path) async {
    final uri = Uri.parse('$baseUrl$path');
    final res = await _client.get(uri, headers: _headers);
    if (res.statusCode >= 200 && res.statusCode < 300) {
      if (res.body.isEmpty) return [];
      return jsonDecode(res.body) as List<dynamic>;
    }
    final error = jsonDecode(res.body) as Map<String, dynamic>;
    throw ApiException(
      res.statusCode,
      error['error'] as String? ?? 'unknown_error',
    );
  }

  Map<String, dynamic> _handleResponse(http.Response res) {
    if (res.statusCode >= 200 && res.statusCode < 300) {
      if (res.body.isEmpty) return {};
      return jsonDecode(res.body) as Map<String, dynamic>;
    }
    final error = jsonDecode(res.body) as Map<String, dynamic>;
    throw ApiException(
      res.statusCode,
      error['error'] as String? ?? 'unknown_error',
    );
  }

  // ─── Sessions ────────────────────────────────────────────

  Future<String> createSession() async {
    final res = await _post('/api/v1/sessions');
    final id = res['session_id'] as String;
    _sessionId = id;
    return id;
  }

  Future<Session> getCurrentSession() async {
    final res = await _get('/api/v1/sessions/current');
    return Session.fromJson(res);
  }

  Future<List<Session>> listHistory() async {
    final list = await _getList('/api/v1/sessions/history');
    return list
        .map((e) => Session.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  // ─── Characters ─────────────────────────────────────────

  Future<List<Character>> listCharacters() async {
    final list = await _getList('/api/v1/characters');
    return list
        .map((e) => Character.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<Character> getCharacter(String charId) async {
    final res = await _get('/api/v1/characters/$charId');
    return Character.fromJson(res);
  }

  // ─── Evidence ────────────────────────────────────────────

  Future<List<Evidence>> listEvidence() async {
    final list = await _getList('/api/v1/evidence');
    return list
        .map((e) => Evidence.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<Evidence> getEvidence(String evId) async {
    final res = await _get('/api/v1/evidence/$evId');
    return Evidence.fromJson(res);
  }

  // ─── Chronology ──────────────────────────────────────────

  Future<List<ChronologyEntry>> getChronology() async {
    final list = await _getList('/api/v1/chronology');
    return list
        .map((e) => ChronologyEntry.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<void> updateNotebookEntry({
    required String chronId,
    required String noteId,
    required List<String> tags,
    String? note,
  }) async {
    await _patch(
      '/api/v1/chronology/$chronId/notes/$noteId',
      body: {'tags': tags, if (note != null) 'note': note},
    );
  }

  // ─── Interrogations ──────────────────────────────────────

  Future<Interrogation> createInterrogation(String characterId) async {
    final res = await _post(
      '/api/v1/interrogations',
      body: {'character_id': characterId},
    );
    return Interrogation.fromJson(res);
  }

  Future<Interrogation?> getActiveInterrogation() async {
    try {
      final res = await _get('/api/v1/interrogations/active');
      return Interrogation.fromJson(res);
    } on ApiException catch (e) {
      if (e.statusCode == 404) return null;
      rethrow;
    }
  }

  Future<Interrogation> getInterrogation(String interId) async {
    final res = await _get('/api/v1/interrogations/$interId');
    return Interrogation.fromJson(res);
  }

  Future<ChatMessage> addInterrogationMessage({
    required String interId,
    required String message,
  }) async {
    final res = await _post(
      '/api/v1/interrogations/$interId/messages',
      body: {'message': message},
    );
    return ChatMessage.fromJson(res);
  }

  Future<List<ChatMessage>> getInterrogationMessages(String interId) async {
    final list = await _getList('/api/v1/interrogations/$interId/messages');
    return list
        .map((e) => ChatMessage.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<void> completeInterrogation(String interId) async {
    await _patch('/api/v1/interrogations/$interId/complete');
  }

  // ─── Actions ─────────────────────────────────────────────

  Future<ActionReport> dnaAnalysis(String evidenceId) async {
    final res = await _post(
      '/api/v1/actions/dna-analysis',
      body: {'evidence_id': evidenceId},
    );
    return ActionReport.fromJson(res);
  }

  Future<ActionReport> fingerprintsCheck(String evidenceId) async {
    final res = await _post(
      '/api/v1/actions/fingerprints',
      body: {'evidence_id': evidenceId},
    );
    return ActionReport.fromJson(res);
  }

  Future<ActionReport> alibiCheck({
    required String characterId,
    required String alibiText,
  }) async {
    final res = await _post(
      '/api/v1/actions/alibi-check',
      body: {'character_id': characterId, 'alibi_text': alibiText},
    );
    return ActionReport.fromJson(res);
  }

  Future<ActionReport> cameraReview() async {
    final res = await _post('/api/v1/actions/camera-review');
    return ActionReport.fromJson(res);
  }

  Future<ActionReport> callHistory(String characterId) async {
    final res = await _post(
      '/api/v1/actions/call-history',
      body: {'character_id': characterId},
    );
    return ActionReport.fromJson(res);
  }

  Future<ActionReport> transactionCheck(String characterId) async {
    final res = await _post(
      '/api/v1/actions/transactions',
      body: {'character_id': characterId},
    );
    return ActionReport.fromJson(res);
  }

  // ─── Reports ─────────────────────────────────────────────

  Future<List<ActionReport>> listReports() async {
    final list = await _getList('/api/v1/reports');
    return list
        .map((item) => ActionReport.fromJson(item as Map<String, dynamic>))
        .toList();
  }

  Future<GameResult> submitReport({
    required String who,
    required String why,
    required String how,
    required String when,
    required String evidence,
  }) async {
    final res = await _post(
      '/api/v1/reports',
      body: {
        'who': who,
        'why': why,
        'how': how,
        'when': when,
        'evidence': evidence,
      },
    );
    return GameResult.fromJson(res);
  }

  void dispose() => _client.close();
}
