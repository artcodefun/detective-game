import 'package:flutter/material.dart';

class ActionsScreen extends StatelessWidget {
  const ActionsScreen({super.key});

  static const _actions = [
    _ActionData(icon: Icons.science, name: 'Анализ ДНК', description: 'Исследовать вещдоки на наличие ДНК', cost: 1),
    _ActionData(
      icon: Icons.fingerprint,
      name: 'Отпечатки пальцев',
      description: 'Проверить отпечатки на вещдоках',
      cost: 1,
    ),
    _ActionData(
      icon: Icons.phone_in_talk,
      name: 'История звонков',
      description: 'Запросить детализацию звонков подозреваемого',
      cost: 2,
    ),
    _ActionData(
      icon: Icons.videocam,
      name: 'Записи с камер',
      description: 'Просмотреть записи камер наблюдения',
      cost: 2,
    ),
    _ActionData(
      icon: Icons.account_balance,
      name: 'Банковские операции',
      description: 'Проверить движение средств по счетам',
      cost: 2,
    ),
    _ActionData(
      icon: Icons.access_time,
      name: 'Проверка алиби',
      description: 'Сверить показания с фактическим временем',
      cost: 1,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Действия')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Потратьте очки действий, чтобы заказать анализ или запросить информацию',
              style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurface.withAlpha(140)),
            ),
            const SizedBox(height: 16),
            Expanded(
              child: GridView.count(
                crossAxisCount: 2,
                mainAxisSpacing: 12,
                crossAxisSpacing: 12,
                childAspectRatio: 1.3,
                children: _actions.map((a) => _ActionCard(action: a)).toList(),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ActionData {
  final IconData icon;
  final String name;
  final String description;
  final int cost;

  const _ActionData({required this.icon, required this.name, required this.description, required this.cost});
}

class _ActionCard extends StatelessWidget {
  final _ActionData action;

  const _ActionCard({required this.action});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Card(
      color: colorScheme.surfaceContainerHighest,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: InkWell(
        onTap: () => _openDetail(context),
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(action.icon, size: 32, color: colorScheme.primary),
              const SizedBox(height: 8),
              Flexible(
                child: Text(
                  action.name,
                  style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold),
                  textAlign: TextAlign.center,
                ),
              ),
              const SizedBox(height: 6),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                decoration: BoxDecoration(
                  color: colorScheme.primary.withAlpha(25),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Text(
                  '${action.cost} AP',
                  style: TextStyle(fontSize: 11, color: colorScheme.primary, fontWeight: FontWeight.w600),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _openDetail(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(20))),
      builder:
          (_) => Padding(
            padding: const EdgeInsets.fromLTRB(20, 12, 20, 20),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Center(
                  child: Container(
                    width: 40,
                    height: 4,
                    decoration: BoxDecoration(
                      color: colorScheme.onSurface.withAlpha(60),
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Icon(action.icon, size: 24, color: colorScheme.primary),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        action.name,
                        style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                Text(action.description, style: theme.textTheme.bodyMedium),
                const SizedBox(height: 8),
                Text(
                  'Стоимость: ${action.cost} AP',
                  style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.primary),
                ),
                const SizedBox(height: 16),
                SizedBox(
                  width: double.infinity,
                  child: FilledButton(onPressed: () => Navigator.pop(context), child: const Text('Заказать')),
                ),
              ],
            ),
          ),
    );
  }
}
