import 'package:flutter/material.dart';

import '../models/report.dart';
import 'title_screen.dart';

class ResultsScreen extends StatelessWidget {
  final GameResult result;

  const ResultsScreen({super.key, required this.result});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final breakdown = result.breakdown;
    final pct = (breakdown.accuracy * 100).toInt();

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
              pct >= 80 ? Icons.verified : pct >= 40 ? Icons.rate_review : Icons.sentiment_dissatisfied,
              size: 80,
              color: pct >= 80
                  ? Colors.green
                  : pct >= 40
                      ? Colors.orange
                      : Colors.red,
            ),
            const SizedBox(height: 16),
            Text(
              '$pct%',
              style: theme.textTheme.displaySmall?.copyWith(fontWeight: FontWeight.bold, color: colorScheme.primary),
            ),
            const SizedBox(height: 8),
            Text(
              'точность',
              style: theme.textTheme.bodyLarge?.copyWith(color: colorScheme.onSurface.withAlpha(150)),
            ),
            const SizedBox(height: 8),
            Text(
              result.narrativeFeedback,
              textAlign: TextAlign.center,
              style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(140)),
            ),
            const SizedBox(height: 24),
            _ResultRow(label: 'Преступник', correct: breakdown.whoCorrect,
                detail: result.breakdownDetails['who']),
            _ResultRow(label: 'Мотив', correct: breakdown.whyCorrect,
                detail: result.breakdownDetails['why']),
            _ResultRow(label: 'Способ', correct: breakdown.howCorrect,
                detail: result.breakdownDetails['how']),
            _ResultRow(label: 'Время', correct: breakdown.whenCorrect,
                detail: result.breakdownDetails['when']),
            _ResultRow(label: 'Улики', correct: breakdown.evidenceCorrect,
                detail: result.breakdownDetails['evidence']),
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
  final String? detail;

  const _ResultRow({required this.label, required this.correct, this.detail});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = correct ? Colors.green : Colors.red;

    return Padding(
      padding: const EdgeInsets.only(bottom: 6),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(correct ? Icons.check_circle : Icons.cancel, color: color, size: 20),
              const SizedBox(width: 8),
              Text(label, style: theme.textTheme.bodyMedium),
              const Spacer(),
              Text(
                correct ? 'Верно' : 'Неверно',
                style: theme.textTheme.bodySmall?.copyWith(color: color.withAlpha(200)),
              ),
            ],
          ),
          if (detail != null) ...[
            const SizedBox(height: 2),
            Padding(
              padding: const EdgeInsets.only(left: 28),
              child: Text(detail!, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurface.withAlpha(120))),
            ),
          ],
        ],
      ),
    );
  }
}
