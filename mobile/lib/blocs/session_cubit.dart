import 'package:flutter_bloc/flutter_bloc.dart';

import '../models/game_state.dart';

class SessionCubit extends Cubit<GameSession?> {
  SessionCubit() : super(null);

  void newGame(GameSession session) => emit(session);

  void update(GameSession session) => emit(session);
}
