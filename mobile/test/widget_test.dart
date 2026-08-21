import 'package:flutter_test/flutter_test.dart';

import 'package:detective_game/main.dart';
import 'package:detective_game/services/api_service.dart';

void main() {
  testWidgets('Title screen shows three main buttons', (WidgetTester tester) async {
    final api = ApiService(baseUrl: 'http://localhost:8080');
    await tester.pumpWidget(DetectiveGameApp(api: api));
    expect(find.text('ДетектИИв'), findsOneWidget);
    expect(find.text('Новое дело'), findsOneWidget);
    expect(find.text('Предыдущие дела'), findsOneWidget);
    expect(find.text('Настройки'), findsOneWidget);
  });
}
