import 'package:flutter_test/flutter_test.dart';

import 'package:detective_game/main.dart';

void main() {
  testWidgets('Title screen shows three main buttons', (WidgetTester tester) async {
    await tester.pumpWidget(const DetectiveGameApp());
    expect(find.text('Детектив'), findsOneWidget);
    expect(find.text('Новое дело'), findsOneWidget);
    expect(find.text('Предыдущие дела'), findsOneWidget);
    expect(find.text('Настройки'), findsOneWidget);
  });
}
