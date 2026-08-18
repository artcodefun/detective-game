import 'package:flutter/foundation.dart';

import 'api_service.dart';

class SessionState {
  final String sessionId;
  final String caseName;
  final int actionPoints;
  final String phase;

  const SessionState({required this.sessionId, this.caseName = '', this.actionPoints = 5, this.phase = 'idle'});

  SessionState copyWith({int? actionPoints, String? phase}) => SessionState(
    sessionId: sessionId,
    caseName: caseName,
    actionPoints: actionPoints ?? this.actionPoints,
    phase: phase ?? this.phase,
  );
}

class SessionService extends ChangeNotifier {
  final ApiService _api;
  SessionState? _state;

  SessionService(this._api);

  SessionState? get state => _state;

  Future<void> startNewGame() async {
    final sessionId = await _api.createSession();
    final session = await _api.getCurrentSession();
    _setState(
      SessionState(
        sessionId: sessionId,
        caseName: session.caseName,
        actionPoints: session.actionPoints,
        phase: session.phase,
      ),
    );
  }

  Future<void> refresh() async {
    final currentState = _state;
    if (currentState == null) return;

    final session = await _api.getCurrentSession();
    if (session.actionPoints != currentState.actionPoints || session.phase != currentState.phase) {
      _setState(currentState.copyWith(actionPoints: session.actionPoints, phase: session.phase));
    }
  }

  void resume(SessionState session) {
    _api.setSessionId(session.sessionId);
    _setState(session);
  }

  void clear() {
    _api.clearSessionId();
    _setState(null);
  }

  void _setState(SessionState? value) {
    _state = value;
    notifyListeners();
  }
}
