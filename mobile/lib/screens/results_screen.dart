import 'package:flutter/material.dart';

import 'title_screen.dart';

class ResultsScreen extends StatelessWidget {
  const ResultsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Результаты'),
        automaticallyImplyLeading: false,
      ),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: [
            const Spacer(flex: 2),
            Icon(
              Icons.verified,
              size: 80,
              color: colorScheme.primary,
            ),
            const SizedBox(height: 16),
            Text(
              '80%',
              style: theme.textTheme.displaySmall?.copyWith(
                fontWeight: FontWeight.bold,
                color: colorScheme.primary,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'точность',
              style: theme.textTheme.bodyLarge?.copyWith(
                color: colorScheme.onSurface.withAlpha(150),
              ),
            ),
            const SizedBox(height: 32),
            _ResultRow(label: 'Преступник', correct: true),
            _ResultRow(label: 'Мотив', correct: true),
            _ResultRow(label: 'Способ', correct: true),
            _ResultRow(label: 'Время', correct: false),
            _ResultRow(label: 'Улики', correct: true),
            const Spacer(flex: 2),
            SizedBox(
              width: double.infinity,
              height: 48,
              child: ElevatedButton.icon(
                onPressed: () {
                  Navigator.pushAndRemoveUntil(
                    context,
                    MaterialPageRoute(builder: (_) => const TitleScreen()),
                    (_) => false,
                  );
                },
                icon: const Icon(Icons.home),
                label: const Text('На главную'),
              ),
            ),
            const Spacer(),
          ],
        ),
      ),
    );
  }
}

class _ResultRow extends StatelessWidget {
  final String label;
  final bool correct;

  const _ResultRow({required this.label, required this.correct});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = correct ? Colors.green : Colors.red;

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Row(
        children: [
          Icon(correct ? Icons.check_circle : Icons.cancel, color: color),
          const SizedBox(width: 12),
          Text(
            label,
            style: theme.textTheme.bodyLarge,
          ),
          const Spacer(),
          Text(
            correct ? 'Верно' : 'Неверно',
            style: theme.textTheme.bodySmall?.copyWith(
              color: color.withAlpha(200),
            ),
          ),
        ],
      ),
    );
  }
}
