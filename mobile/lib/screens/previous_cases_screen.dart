import 'package:flutter/material.dart';

class PreviousCasesScreen extends StatelessWidget {
  const PreviousCasesScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Предыдущие дела'),
      ),
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.inbox,
              size: 64,
              color: theme.colorScheme.onSurface.withAlpha(80),
            ),
            const SizedBox(height: 16),
            Text(
              'Пока нет завершённых дел',
              style: theme.textTheme.bodyLarge?.copyWith(
                color: theme.colorScheme.onSurface.withAlpha(120),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
