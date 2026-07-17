import 'package:flutter/material.dart';

class EvidenceLockerScreen extends StatelessWidget {
  const EvidenceLockerScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(
        title: const Text('Улики'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Хранилище вещественных доказательств',
              style: theme.textTheme.titleMedium,
            ),
            const SizedBox(height: 16),
            Expanded(
              child: GridView.count(
                crossAxisCount: 2,
                mainAxisSpacing: 12,
                crossAxisSpacing: 12,
                childAspectRatio: 1.2,
                children: [
                  _EvidenceCard(
                    icon: Icons.science,
                    name: 'Нож',
                    description: 'Кухонный нож со следами',
                  ),
                  _EvidenceCard(
                    icon: Icons.medication,
                    name: 'Пузырёк',
                    description: 'Остатки неизвестного вещества',
                  ),
                  _EvidenceCard(
                    icon: Icons.description,
                    name: 'Записка',
                    description: 'Анонимное письмо',
                  ),
                  _EvidenceCard(
                    icon: Icons.fingerprint,
                    name: 'Отпечатки',
                    description: 'С места преступления',
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _EvidenceCard extends StatelessWidget {
  final IconData icon;
  final String name;
  final String description;

  const _EvidenceCard({
    required this.icon,
    required this.name,
    required this.description,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Card(
      color: colorScheme.surfaceContainerHighest,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, size: 32, color: colorScheme.primary),
            const SizedBox(height: 8),
            Text(
              name,
              style: theme.textTheme.titleSmall?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              description,
              textAlign: TextAlign.center,
              style: theme.textTheme.bodySmall?.copyWith(
                color: colorScheme.onSurface.withAlpha(140),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
