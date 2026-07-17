import '../models/character.dart';
import '../models/game_state.dart';
import '../models/report.dart';
import '../models/scenario.dart';

class LlmInterrogationResponse {
  final String answer;
  final int attitudeDelta;
  final List<String> statements;

  const LlmInterrogationResponse({
    required this.answer,
    required this.attitudeDelta,
    required this.statements,
  });
}

class LlmFeedbackResponse {
  final String narrativeFeedback;
  final Map<String, String> breakdownDetails;

  const LlmFeedbackResponse({
    required this.narrativeFeedback,
    required this.breakdownDetails,
  });
}

abstract class LlmService {
  Future<GameSession> generateScenario(List<CharacterData> selectedCharacters);

  Future<LlmInterrogationResponse> respondInInterrogation({
    required CharacterState characterState,
    required String playerMessage,
  });

  Future<LlmFeedbackResponse> evaluateReport({
    required FinalReport playerReport,
    required Crime groundTruth,
  });
}
