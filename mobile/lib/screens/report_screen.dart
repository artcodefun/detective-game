import 'package:flutter/material.dart';

import 'results_screen.dart';

class ReportScreen extends StatelessWidget {
  const ReportScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Финальный отчёт'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Изложите вашу версию событий',
              style: theme.textTheme.titleMedium,
            ),
            const SizedBox(height: 4),
            Text(
              'Осветите все обязательные пункты',
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurface.withAlpha(120),
              ),
            ),
            const SizedBox(height: 20),
            Expanded(
              child: TextField(
                maxLines: null,
                expands: true,
                textAlignVertical: TextAlignVertical.top,
                decoration: InputDecoration(
                  hintText:
                      'Кто совершил преступление?\nПочему?\nКаким способом?\nВ какое время?\nКакие улики это подтверждают?',
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  contentPadding: const EdgeInsets.all(16),
                ),
              ),
            ),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              height: 48,
              child: ElevatedButton.icon(
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (_) => const ResultsScreen(),
                    ),
                  );
                },
                icon: const Icon(Icons.send),
                label: const Text('Отправить отчёт'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
