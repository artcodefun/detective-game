import 'dart:io';
import 'dart:convert';

import 'package:device_info_plus/device_info_plus.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart' as http;
import 'package:package_info_plus/package_info_plus.dart';

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

enum InitializationStatus {
  initial,
  checkingVersion,
  updateRequired,
  registering,
  ready,
  versionCheckFailed,
  registrationFailed,
}

class ApiService extends ChangeNotifier {
  static const _accessTokenKey = 'detective_access_token';

  final String baseUrl;
  final http.Client _client;
  final FlutterSecureStorage _secureStorage;
  String _accessToken = '';
  String? _sessionId;
  String? _updateUrl;
  InitializationStatus _initializationStatus = InitializationStatus.initial;
  Object? _initializationError;

  String? get sessionId => _sessionId;

  InitializationStatus get initializationStatus => _initializationStatus;

  Object? get initializationError => _initializationError;

  String? get updateUrl => _updateUrl;

  ApiService({required this.baseUrl, http.Client? client, FlutterSecureStorage? secureStorage})
    : _client = client ?? http.Client(),
      _secureStorage = secureStorage ?? const FlutterSecureStorage();

  Future<void> initialize() async {
    if (_initializationStatus == InitializationStatus.ready ||
        _initializationStatus == InitializationStatus.updateRequired ||
        _initializationStatus == InitializationStatus.checkingVersion ||
        _initializationStatus == InitializationStatus.registering) {
      return;
    }

    _initializationStatus = InitializationStatus.checkingVersion;
    _initializationError = null;
    _updateUrl = null;
    notifyListeners();

    try {
      if (await _checkVersion()) {
        _initializationStatus = InitializationStatus.updateRequired;
        return;
      }
    } catch (error) {
      _initializationError = error;
      _initializationStatus = InitializationStatus.versionCheckFailed;
      notifyListeners();
      return;
    }

    try {
      final storedToken = await _secureStorage.read(key: _accessTokenKey);
      if (storedToken != null && storedToken.isNotEmpty) {
        _accessToken = storedToken;
        _initializationStatus = InitializationStatus.ready;
        return;
      }

      _initializationStatus = InitializationStatus.registering;
      notifyListeners();
      await _registerAnonymous();
      _initializationStatus = InitializationStatus.ready;
    } catch (error) {
      _initializationError = error;
      _initializationStatus = InitializationStatus.registrationFailed;
    } finally {
      notifyListeners();
    }
  }

  Future<bool> _checkVersion() async {
    final package = await PackageInfo.fromPlatform();
    final version = await _post(
      '/api/v1/app/version',
      body: {'platform': Platform.isIOS ? 'ios' : 'android', 'version': package.version},
    );
    if (version['update_required'] != true) return false;

    _updateUrl = version['update_url'] as String?;
    return true;
  }

  Future<void> _registerAnonymous() async {
    final device = await _deviceInfo();
    final response = await _post('/api/v1/auth/anonymous', body: device);
    _accessToken = response['access_token'] as String;
    await _secureStorage.write(key: _accessTokenKey, value: _accessToken);
  }

  Future<Map<String, String>> _deviceInfo() async {
    final deviceInfo = DeviceInfoPlugin();
    if (Platform.isIOS) {
      final info = await deviceInfo.iosInfo;
      return {
        'platform': 'ios',
        'manufacturer': 'Apple',
        'model': info.utsname.machine,
        'os_version': info.systemVersion,
      };
    }
    if (Platform.isAndroid) {
      final info = await deviceInfo.androidInfo;
      return {
        'platform': 'android',
        'manufacturer': info.manufacturer,
        'model': info.model,
        'os_version': info.version.release,
      };
    }
    throw UnsupportedError('unsupported platform');
  }

  void setSessionId(String id) => _sessionId = id;

  void clearSessionId() => _sessionId = null;

  Map<String, String> get _headers => {
    'Content-Type': 'application/json',
    'Accept-Language': 'ru',
    if (_accessToken.isNotEmpty) 'Authorization': 'Bearer $_accessToken',
    if (_sessionId != null) 'X-Session-ID': _sessionId!,
  };

  Future<Map<String, dynamic>> _get(String path) async {
    final uri = Uri.parse('$baseUrl$path');
    final res = await _client.get(uri, headers: _headers);
    return _handleResponse(res);
  }

  Future<Map<String, dynamic>> _post(String path, {Map<String, dynamic>? body}) async {
    final uri = Uri.parse('$baseUrl$path');
    final res = await _client.post(uri, headers: _headers, body: body != null ? jsonEncode(body) : null);
    return _handleResponse(res);
  }

  Future<Map<String, dynamic>> _patch(String path, {Map<String, dynamic>? body}) async {
    final uri = Uri.parse('$baseUrl$path');
    final res = await _client.patch(uri, headers: _headers, body: body != null ? jsonEncode(body) : null);
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
    throw ApiException(res.statusCode, error['error'] as String? ?? 'unknown_error');
  }

  Map<String, dynamic> _handleResponse(http.Response res) {
    if (res.statusCode >= 200 && res.statusCode < 300) {
      if (res.body.isEmpty) return {};
      return jsonDecode(res.body) as Map<String, dynamic>;
    }
    final error = jsonDecode(res.body) as Map<String, dynamic>;
    throw ApiException(res.statusCode, error['error'] as String? ?? 'unknown_error');
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
    return list.map((e) => Session.fromJson(e as Map<String, dynamic>)).toList();
  }

  // ─── Characters ─────────────────────────────────────────

  Future<List<Character>> listCharacters() async {
    final list = await _getList('/api/v1/characters');
    return list.map((e) => Character.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<Character> getCharacter(String charId) async {
    final res = await _get('/api/v1/characters/$charId');
    return Character.fromJson(res);
  }

  // ─── Evidence ────────────────────────────────────────────

  Future<List<Evidence>> listEvidence() async {
    final list = await _getList('/api/v1/evidence');
    return list.map((e) => Evidence.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<Evidence> getEvidence(String evId) async {
    final res = await _get('/api/v1/evidence/$evId');
    return Evidence.fromJson(res);
  }

  // ─── Chronology ──────────────────────────────────────────

  Future<List<ChronologyEntry>> getChronology() async {
    final list = await _getList('/api/v1/chronology');
    return list.map((e) => ChronologyEntry.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<void> updateNotebookEntry({
    required String chronId,
    required String noteId,
    required List<String> tags,
    String? note,
  }) async {
    await _patch('/api/v1/chronology/$chronId/notes/$noteId', body: {'tags': tags, if (note != null) 'note': note});
  }

  // ─── Interrogations ──────────────────────────────────────

  Future<Interrogation> createInterrogation(String characterId) async {
    final res = await _post('/api/v1/interrogations', body: {'character_id': characterId});
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

  Future<ChatMessage> addInterrogationMessage({required String interId, required String message}) async {
    final res = await _post('/api/v1/interrogations/$interId/messages', body: {'message': message});
    return ChatMessage.fromJson(res);
  }

  Future<List<ChatMessage>> getInterrogationMessages(String interId) async {
    final list = await _getList('/api/v1/interrogations/$interId/messages');
    return list.map((e) => ChatMessage.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<void> completeInterrogation(String interId) async {
    await _patch('/api/v1/interrogations/$interId/complete');
  }

  // ─── Actions ─────────────────────────────────────────────

  Future<ActionReport> dnaAnalysis(String evidenceId) async {
    final res = await _post('/api/v1/actions/dna-analysis', body: {'evidence_id': evidenceId});
    return ActionReport.fromJson(res);
  }

  Future<ActionReport> fingerprintsCheck(String evidenceId) async {
    final res = await _post('/api/v1/actions/fingerprints', body: {'evidence_id': evidenceId});
    return ActionReport.fromJson(res);
  }

  Future<ActionReport> alibiCheck({required String characterId, required String alibiText}) async {
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
    final res = await _post('/api/v1/actions/call-history', body: {'character_id': characterId});
    return ActionReport.fromJson(res);
  }

  Future<ActionReport> transactionCheck(String characterId) async {
    final res = await _post('/api/v1/actions/transactions', body: {'character_id': characterId});
    return ActionReport.fromJson(res);
  }

  // ─── Reports ─────────────────────────────────────────────

  Future<List<ActionReport>> listReports() async {
    final list = await _getList('/api/v1/reports');
    return list.map((item) => ActionReport.fromJson(item as Map<String, dynamic>)).toList();
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
      body: {'who': who, 'why': why, 'how': how, 'when': when, 'evidence': evidence},
    );
    return GameResult.fromJson(res);
  }

  @override
  void dispose() {
    _client.close();
    super.dispose();
  }
}
