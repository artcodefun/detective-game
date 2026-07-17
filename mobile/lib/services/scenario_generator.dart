import '../models/character.dart';
import '../models/game_state.dart';
import 'llm_service.dart';

class ScenarioGenerator {
  final LlmService _llm;

  ScenarioGenerator(this._llm);

  Future<GameSession> generate(List<CharacterData> selectedCharacters) async {
    return _llm.generateScenario(selectedCharacters);
  }
}
