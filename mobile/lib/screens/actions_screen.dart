import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../blocs/session_cubit.dart';
import '../models/game_state.dart';
import '../models/scenario.dart';
import '../services/api_service.dart';

enum _ActionKind { evidenceAnalysis, suspectAction, camera, alibi }

class _ActionData {
  final IconData icon;
  final String name;
  final String description;
  final int cost;
  final _ActionKind kind;

  const _ActionData({
    required this.icon,
    required this.name,
    required this.description,
    required this.cost,
    required this.kind,
  });
}

final _actions = [
  _ActionData(
    icon: Icons.science,
    name: 'Анализ ДНК',
    description: 'Исследовать вещдоки на наличие ДНК',
    cost: 1,
    kind: _ActionKind.evidenceAnalysis,
  ),
  _ActionData(
    icon: Icons.fingerprint,
    name: 'Отпечатки пальцев',
    description: 'Проверить отпечатки на вещдоках',
    cost: 1,
    kind: _ActionKind.evidenceAnalysis,
  ),
  _ActionData(
    icon: Icons.phone_in_talk,
    name: 'История звонков',
    description: 'Запросить детализацию звонков подозреваемого',
    cost: 2,
    kind: _ActionKind.suspectAction,
  ),
  _ActionData(
    icon: Icons.videocam,
    name: 'Записи с камер',
    description: 'Просмотреть записи камер наблюдения',
    cost: 2,
    kind: _ActionKind.camera,
  ),
  _ActionData(
    icon: Icons.account_balance,
    name: 'Банковские операции',
    description: 'Проверить движение средств по счетам',
    cost: 2,
    kind: _ActionKind.suspectAction,
  ),
  _ActionData(
    icon: Icons.access_time,
    name: 'Проверка алиби',
    description: 'Сверить показания с фактическим временем',
    cost: 1,
    kind: _ActionKind.alibi,
  ),
];

class ActionsScreen extends StatelessWidget {
  const ActionsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final ap = context.watch<SessionCubit>().state?.actionPoints ?? 0;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Действия'),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 12),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.search, size: 18, color: theme.colorScheme.primary),
                const SizedBox(width: 4),
                Text(
                  '$ap',
                  style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold),
                ),
              ],
            ),
          ),
        ],
      ),
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
        onTap: () => _openSheet(context),
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

  void _openSheet(BuildContext context) {
    final session = context.read<SessionCubit>().state;
    if (session == null) return;
    if (session.actionPoints < action.cost) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Недостаточно очков действий')));
      return;
    }
    _showActionSheet(context);
  }

  Future<void> _executeAction(
    BuildContext context, {
    String? evidenceId,
    String? characterId,
    String? alibiText,
  }) async {
    final api = context.read<ApiService>();
    final cubit = context.read<SessionCubit>();
    try {
      switch (action.name) {
        case 'Анализ ДНК':
          await api.dnaAnalysis(evidenceId!);
        case 'Отпечатки пальцев':
          await api.fingerprintsCheck(evidenceId!);
        case 'История звонков':
          await api.callHistory(characterId!);
        case 'Записи с камер':
          await api.cameraReview();
        case 'Банковские операции':
          await api.transactionCheck(characterId!);
        case 'Проверка алиби':
          await api.alibiCheck(characterId: characterId!, alibiText: alibiText!);
        default:
          return;
      }
      await cubit.refreshSession();
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('${action.name} — выполнен')));
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Ошибка: $e')));
      }
    }
  }

  Future<void> _showActionSheet(BuildContext context) async {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final api = context.read<ApiService>();

    Widget Function(BuildContext ctx) buildContent;
    switch (action.kind) {
      case _ActionKind.evidenceAnalysis:
        String? selected;
        buildContent =
            (ctx) => StatefulBuilder(
              builder: (ctx, setSheetState) {
                return FutureBuilder<List<Evidence>>(
                  future: api.listEvidence().then(
                    (list) =>
                        list
                            .where(
                              (e) =>
                                  ![
                                    'dna_analysis',
                                    'fingerprints',
                                    'alibi_check',
                                    'camera_review',
                                    'call_history',
                                    'transaction_check',
                                  ].contains(e.type),
                            )
                            .toList(),
                  ),
                  builder: (_, snapshot) {
                    final evidenceList = snapshot.data ?? [];
                    return Column(
                      mainAxisSize: MainAxisSize.min,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Выберите улику для анализа',
                          style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(140)),
                        ),
                        const SizedBox(height: 12),
                        ConstrainedBox(
                          constraints: const BoxConstraints(maxHeight: 240),
                          child: ListView(
                            shrinkWrap: true,
                            children:
                                evidenceList.map((e) {
                                  final isSelected = selected == e.id;
                                  return ListTile(
                                    title: Text(e.name),
                                    subtitle: Text(e.description, maxLines: 1, overflow: TextOverflow.ellipsis),
                                    leading: Icon(
                                      isSelected ? Icons.radio_button_checked : Icons.radio_button_unchecked,
                                    ),
                                    onTap: () => setSheetState(() => selected = e.id),
                                    dense: true,
                                  );
                                }).toList(),
                          ),
                        ),
                        const SizedBox(height: 12),
                        SizedBox(
                          width: double.infinity,
                          child: FilledButton(
                            onPressed:
                                selected != null
                                    ? () {
                                      Navigator.pop(ctx);
                                      _executeAction(ctx, evidenceId: selected);
                                    }
                                    : null,
                            child: Text('Заказать (${action.cost} AP)'),
                          ),
                        ),
                      ],
                    );
                  },
                );
              },
            );
      case _ActionKind.suspectAction:
        String? selected;
        buildContent =
            (ctx) => StatefulBuilder(
              builder: (ctx, setSheetState) {
                return FutureBuilder<List<Character>>(
                  future: api.listCharacters(),
                  builder: (_, snapshot) {
                    final chars = snapshot.data ?? [];
                    return Column(
                      mainAxisSize: MainAxisSize.min,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Выберите подозреваемого',
                          style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(140)),
                        ),
                        const SizedBox(height: 12),
                        ConstrainedBox(
                          constraints: const BoxConstraints(maxHeight: 240),
                          child: ListView(
                            shrinkWrap: true,
                            children:
                                chars.map((c) {
                                  final isSelected = selected == c.id;
                                  return ListTile(
                                    title: Text(c.name),
                                    subtitle: Text(c.profession, maxLines: 1, overflow: TextOverflow.ellipsis),
                                    leading: Icon(
                                      isSelected ? Icons.radio_button_checked : Icons.radio_button_unchecked,
                                    ),
                                    onTap: () => setSheetState(() => selected = c.id),
                                    dense: true,
                                  );
                                }).toList(),
                          ),
                        ),
                        const SizedBox(height: 12),
                        SizedBox(
                          width: double.infinity,
                          child: FilledButton(
                            onPressed:
                                selected != null
                                    ? () {
                                      Navigator.pop(ctx);
                                      _executeAction(ctx, characterId: selected);
                                    }
                                    : null,
                            child: Text('Заказать (${action.cost} AP)'),
                          ),
                        ),
                      ],
                    );
                  },
                );
              },
            );
      case _ActionKind.camera:
        buildContent =
            (ctx) => Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Будут запрошены записи с камер наблюдения за вечер предполагаемого преступления.',
                  style: theme.textTheme.bodyMedium,
                ),
                const SizedBox(height: 8),
                Text(
                  'Стоимость: ${action.cost} AP',
                  style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.primary),
                ),
                const SizedBox(height: 16),
                SizedBox(
                  width: double.infinity,
                  child: FilledButton(
                    onPressed: () {
                      Navigator.pop(ctx);
                      _executeAction(ctx);
                    },
                    child: const Text('Заказать'),
                  ),
                ),
              ],
            );
      case _ActionKind.alibi:
        buildContent = (ctx) {
          final controller = TextEditingController();
          return StatefulBuilder(
            builder: (ctx, setSheetState) {
              final hasText = controller.text.trim().isNotEmpty;
              return Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('Проверка алиби — опрос свидетелей и проверка места.', style: theme.textTheme.bodyMedium),
                  const SizedBox(height: 12),
                  TextField(
                    controller: controller,
                    maxLines: 4,
                    onChanged: (_) => setSheetState(() {}),
                    decoration: InputDecoration(
                      hintText: 'Опишите алиби, которое нужно проверить...',
                      border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                    ),
                  ),
                  const SizedBox(height: 16),
                  SizedBox(
                    width: double.infinity,
                    child: FilledButton(
                      onPressed:
                          hasText
                              ? () {
                                Navigator.pop(ctx);
                                _executeAction(ctx, alibiText: controller.text.trim());
                              }
                              : null,
                      child: Text('Заказать (${action.cost} AP)'),
                    ),
                  ),
                ],
              );
            },
          );
        };
    }

    if (!context.mounted) return;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(20))),
      builder:
          (ctx) => Padding(
            padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom),
            child: SingleChildScrollView(
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
                  Text(action.name, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 12),
                  buildContent(ctx),
                ],
              ),
            ),
          ),
    );
  }
}
