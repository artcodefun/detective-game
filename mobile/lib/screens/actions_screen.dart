import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../blocs/session_cubit.dart';
import '../models/chronology_entry.dart';
import '../models/game_state.dart';
import '../models/notebook.dart';
import '../models/scenario.dart';

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
    description: 'Просмотреть записи камер наблюдения в районе преступления',
    cost: 2,
    kind: _ActionKind.camera,
  ),
  _ActionData(
    icon: Icons.account_balance,
    name: 'Банковские операции',
    description: 'Проверить движение средств по счетам подозреваемого',
    cost: 2,
    kind: _ActionKind.suspectAction,
  ),
  _ActionData(
    icon: Icons.access_time,
    name: 'Проверка алиби',
    description: 'Сверить показания подозреваемого с фактическим временем',
    cost: 1,
    kind: _ActionKind.alibi,
  ),
];

class ActionsScreen extends StatelessWidget {
  const ActionsScreen({super.key});

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
                child: Text(action.name, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold), textAlign: TextAlign.center),
              ),
              const SizedBox(height: 6),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                decoration: BoxDecoration(
                  color: colorScheme.primary.withAlpha(25),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Text('${action.cost} AP', style: TextStyle(fontSize: 11, color: colorScheme.primary, fontWeight: FontWeight.w600)),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _openSheet(BuildContext context) {
    final session = context.read<SessionCubit>().state!;
    if (session.actionPoints < action.cost) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Недостаточно очков действий')),
      );
      return;
    }
    _showActionSheet(context);
  }

  void _executeAction(BuildContext context, {String? evidenceId, String? characterId, String? alibiText}) {
    final session = context.read<SessionCubit>().state!;
    final spent = session.spendActionPoints(action.cost);
    var updated = spent;

    if (evidenceId != null) {
      final idx = updated.evidence.indexWhere((e) => e.id == evidenceId);
      if (idx >= 0) {
        final ev = updated.evidence[idx];
        final evUpdated = ev.copyWith(analyzed: true, analysisResult: '${action.name} завершён. Следов не обнаружено.');
        final list = List<Evidence>.from(updated.evidence)..[idx] = evUpdated;
        updated = updated.copyWith(evidence: list);
      }
    }

    final entry = ChronologyEntry(
      id: 'chron_action_${DateTime.now().millisecondsSinceEpoch}',
      eventType: _eventTypeForAction(),
      title: action.name,
      timestamp: DateTime.now(),
      details: [
        NotebookEntry(
          id: 'note_action_${DateTime.now().millisecondsSinceEpoch}',
          type: NotebookEntryType.analysis,
          description: characterId != null
              ? 'Запрошено для: ${_characterName(session, characterId)}'
              : evidenceId != null
                  ? 'Исследована улика: ${_evidenceName(session, evidenceId)}'
                  : alibiText != null
                      ? 'Проверка алиби: $alibiText'
                      : 'Запрос выполнен',
          timestamp: DateTime.now(),
        ),
      ],
    );
    updated = updated.addChronologyEntry(entry);
    context.read<SessionCubit>().update(updated);

    if (context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('${action.name} — заказ отправлен')),
      );
    }
  }

  ChronologyEventType _eventTypeForAction() {
    switch (action.kind) {
      case _ActionKind.evidenceAnalysis:
        return ChronologyEventType.labAnalysis;
      case _ActionKind.suspectAction:
      case _ActionKind.camera:
        return ChronologyEventType.cameraReview;
      case _ActionKind.alibi:
        return ChronologyEventType.alibiCheck;
    }
  }

  String _characterName(GameSession session, String id) {
    return session.characterById(id)?.base.name ?? id;
  }

  String _evidenceName(GameSession session, String id) {
    final ev = session.evidence.where((e) => e.id == id).firstOrNull;
    return ev?.name ?? id;
  }

  Future<void> _showActionSheet(BuildContext context) async {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final session = context.read<SessionCubit>().state!;

    Widget Function(BuildContext ctx) buildContent;
    switch (action.kind) {
      case _ActionKind.evidenceAnalysis: {
        String? selected;
        buildContent = (ctx) => StatefulBuilder(
          builder: (ctx, setSheetState) => Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('Выберите улику для анализа', style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(140))),
              const SizedBox(height: 12),
              ConstrainedBox(
                constraints: const BoxConstraints(maxHeight: 240),
                child: ListView(
                  shrinkWrap: true,
                  children: session.evidence.map((e) {
                    final isSelected = selected == e.id;
                    return ListTile(
                      title: Text(e.name),
                      subtitle: Text(e.description, maxLines: 1, overflow: TextOverflow.ellipsis),
                      leading: Icon(isSelected ? Icons.radio_button_checked : Icons.radio_button_unchecked),
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
                  onPressed: selected != null
                      ? () { Navigator.pop(ctx); _executeAction(ctx, evidenceId: selected); }
                      : null,
                  child: Text('Заказать (${action.cost} AP)'),
                ),
              ),
            ],
          ),
        );
        break;
      }
      case _ActionKind.suspectAction: {
        String? selected;
        buildContent = (ctx) => StatefulBuilder(
          builder: (ctx, setSheetState) => Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('Выберите подозреваемого', style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(140))),
              const SizedBox(height: 12),
              ConstrainedBox(
                constraints: const BoxConstraints(maxHeight: 240),
                child: ListView(
                  shrinkWrap: true,
                  children: session.characters.map((c) {
                    final isSelected = selected == c.base.id;
                    return ListTile(
                      title: Text(c.base.name),
                      subtitle: Text(c.base.profession, maxLines: 1, overflow: TextOverflow.ellipsis),
                      leading: Icon(isSelected ? Icons.radio_button_checked : Icons.radio_button_unchecked),
                      onTap: () => setSheetState(() => selected = c.base.id),
                      dense: true,
                    );
                  }).toList(),
                ),
              ),
              const SizedBox(height: 12),
              SizedBox(
                width: double.infinity,
                child: FilledButton(
                  onPressed: selected != null
                      ? () { Navigator.pop(ctx); _executeAction(ctx, characterId: selected); }
                      : null,
                  child: Text('Заказать (${action.cost} AP)'),
                ),
              ),
            ],
          ),
        );
        break;
      }
      case _ActionKind.camera:
        buildContent = (ctx) => Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Будут запрошены записи с камер наблюдения, установленных в районе места преступления и на прилегающих улицах за вечер предполагаемого преступления.', style: theme.textTheme.bodyMedium),
            const SizedBox(height: 8),
            Text('Стоимость: ${action.cost} AP', style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.primary)),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: FilledButton(
                onPressed: () { Navigator.pop(ctx); _executeAction(ctx); },
                child: const Text('Заказать'),
              ),
            ),
          ],
        );
        break;
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
                  Text('Мы направим агентов, чтобы проверить алиби подозреваемого. Они смогут опросить персонал, проверить записи в заведениях или найти свидетелей. Проверить можно только те алиби, которые привязаны к конкретному месту — например, «был в ресторане» или «встречался с адвокатом в офисе».', style: theme.textTheme.bodyMedium),
                  const SizedBox(height: 12),
                  TextField(
                    controller: controller,
                    maxLines: 4,
                    onChanged: (_) => setSheetState(() {}),
                    decoration: InputDecoration(
                      hintText: 'Майкл Браун утверждает, что с 20:00 до 21:30 был в ресторане «У Ланга»...',
                      border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                    ),
                  ),
                  const SizedBox(height: 16),
                  SizedBox(
                    width: double.infinity,
                    child: FilledButton(
                      onPressed: hasText
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
        break;
    }

    if (!context.mounted) return;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(20))),
      builder: (ctx) => Padding(
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
