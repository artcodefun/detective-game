import 'package:flutter_bloc/flutter_bloc.dart';

import '../services/api_service.dart';

class SessionState {
  final String sessionId;
  final String caseName;
  final int actionPoints;
  final String phase;

  const SessionState({required this.sessionId, this.caseName = '', this.actionPoints = 5, this.phase = 'idle'});

  SessionState copyWith({String? sessionId, String? caseName, int? actionPoints, String? phase}) {
    return SessionState(
      sessionId: sessionId ?? this.sessionId,
      caseName: caseName ?? this.caseName,
      actionPoints: actionPoints ?? this.actionPoints,
      phase: phase ?? this.phase,
    );
  }
}

class SessionCubit extends Cubit<SessionState?> {
  final ApiService _api;

  SessionCubit(this._api) : super(null);

  ApiService get api => _api;

  Future<void> startNewGame() async {
    final sessionId = await _api.createSession();
    final session = await _api.getCurrentSession();
    emit(
      SessionState(
        sessionId: sessionId,
        caseName: session.caseName,
        actionPoints: session.actionPoints,
        phase: session.phase,
      ),
    );
  }

  Future<void> refreshSession() async {
    final state = this.state;
    if (state == null) return;
    final session = await _api.getCurrentSession();
    final newPoints = session.actionPoints;
    final newPhase = session.phase;
    if (newPoints != state.actionPoints || newPhase != state.phase) {
      emit(state.copyWith(actionPoints: newPoints, phase: newPhase));
    }
  }

  void resumeSession(String sessionId, String caseName, int actionPoints, String phase) {
    _api.setSessionId(sessionId);
    emit(SessionState(sessionId: sessionId, caseName: caseName, actionPoints: actionPoints, phase: phase));
  }

  void clear() {
    _api.clearSessionId();
    emit(null);
  }
}
