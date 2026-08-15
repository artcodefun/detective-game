import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_markdown/flutter_markdown.dart';

import '../blocs/session_cubit.dart';
import '../models/action_report.dart';
import '../models/game_state.dart';
import '../models/scenario.dart';
import '../services/api_service.dart';

enum _ActionKind { evidenceAnalysis, suspectAction, camera, alibi }

enum _ActionRequest {
  dnaAnalysis,
  fingerprints,
  callHistory,
  cameraReview,
  transactionCheck,
  alibiCheck,
}

class _ActionData {
  final IconData icon;
  final String name;
  final String description;
  final int cost;
  final _ActionKind kind;
  final _ActionRequest request;

  const _ActionData({
    required this.icon,
    required this.name,
    required this.description,
    required this.cost,
    required this.kind,
    required this.request,
  });
}

const _actions = [
  _ActionData(
    icon: Icons.science,
    name: 'Анализ ДНК',
    description: 'Исследовать вещдоки на наличие ДНК',
    cost: 1,
    kind: _ActionKind.evidenceAnalysis,
    request: _ActionRequest.dnaAnalysis,
  ),
  _ActionData(
    icon: Icons.fingerprint,
    name: 'Отпечатки пальцев',
    description: 'Проверить отпечатки на вещдоках',
    cost: 1,
    kind: _ActionKind.evidenceAnalysis,
    request: _ActionRequest.fingerprints,
  ),
  _ActionData(
    icon: Icons.phone_in_talk,
    name: 'История звонков',
    description: 'Запросить детализацию звонков подозреваемого',
    cost: 2,
    kind: _ActionKind.suspectAction,
    request: _ActionRequest.callHistory,
  ),
  _ActionData(
    icon: Icons.videocam,
    name: 'Записи с камер',
    description: 'Просмотреть записи камер наблюдения',
    cost: 2,
    kind: _ActionKind.camera,
    request: _ActionRequest.cameraReview,
  ),
  _ActionData(
    icon: Icons.account_balance,
    name: 'Банковские операции',
    description: 'Проверить движение средств по счетам',
    cost: 2,
    kind: _ActionKind.suspectAction,
    request: _ActionRequest.transactionCheck,
  ),
  _ActionData(
    icon: Icons.access_time,
    name: 'Проверка алиби',
    description: 'Сверить показания с фактическим временем',
    cost: 1,
    kind: _ActionKind.alibi,
    request: _ActionRequest.alibiCheck,
  ),
];

class ActionsScreen extends StatelessWidget {
  const ActionsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final actionPoints = context.watch<SessionCubit>().state?.actionPoints ?? 0;
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
                  '$actionPoints',
                  style: theme.textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
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
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurface.withAlpha(140),
              ),
            ),
            const SizedBox(height: 16),
            Expanded(
              child: GridView.count(
                crossAxisCount: 2,
                mainAxisSpacing: 12,
                crossAxisSpacing: 12,
                childAspectRatio: 1.3,
                children:
                    _actions
                        .map((action) => _ActionCard(action: action))
                        .toList(),
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
                  style: theme.textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
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
                  style: TextStyle(
                    fontSize: 11,
                    color: colorScheme.primary,
                    fontWeight: FontWeight.w600,
                  ),
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
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Недостаточно очков действий')),
      );
      return;
    }
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (_) => _ActionSheet(action: action),
    );
  }
}

class _ActionSheet extends StatefulWidget {
  final _ActionData action;

  const _ActionSheet({required this.action});

  @override
  State<_ActionSheet> createState() => _ActionSheetState();
}

class _ActionSheetState extends State<_ActionSheet> {
  Future<List<Evidence>>? _evidenceFuture;
  Future<List<Character>>? _charactersFuture;
  final _alibiController = TextEditingController();
  String? _evidenceID;
  String? _characterID;
  bool _isSubmitting = false;
  ActionReport? _report;

  @override
  void initState() {
    super.initState();
    final api = context.read<ApiService>();
    switch (widget.action.kind) {
      case _ActionKind.evidenceAnalysis:
        _evidenceFuture = api.listEvidence();
      case _ActionKind.suspectAction:
      case _ActionKind.alibi:
        _charactersFuture = api.listCharacters();
      case _ActionKind.camera:
        break;
    }
  }

  @override
  void dispose() {
    _alibiController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    setState(() => _isSubmitting = true);
    try {
      final api = context.read<ApiService>();
      final sessionCubit = context.read<SessionCubit>();
      final report = await switch (widget.action.request) {
        _ActionRequest.dnaAnalysis => api.dnaAnalysis(_evidenceID!),
        _ActionRequest.fingerprints => api.fingerprintsCheck(_evidenceID!),
        _ActionRequest.callHistory => api.callHistory(_characterID!),
        _ActionRequest.cameraReview => api.cameraReview(),
        _ActionRequest.transactionCheck => api.transactionCheck(_characterID!),
        _ActionRequest.alibiCheck => api.alibiCheck(
          characterId: _characterID!,
          alibiText: _alibiController.text.trim(),
        ),
      };
      await sessionCubit.refreshSession();
      if (!mounted) return;
      setState(() {
        _report = report;
        _isSubmitting = false;
      });
    } catch (e) {
      if (!mounted) return;
      setState(() => _isSubmitting = false);
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('Ошибка: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    return PopScope(
      canPop: !_isSubmitting,
      child: SafeArea(
        top: false,
        child: Padding(
          padding: EdgeInsets.fromLTRB(
            20,
            12,
            20,
            20 + MediaQuery.viewInsetsOf(context).bottom,
          ),
          child: ConstrainedBox(
            constraints: BoxConstraints(
              maxHeight: MediaQuery.sizeOf(context).height * 0.85,
            ),
            child: SingleChildScrollView(
              keyboardDismissBehavior: ScrollViewKeyboardDismissBehavior.onDrag,
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
                  if (_report != null)
                    _buildReport(theme, colorScheme, _report!)
                  else ...[
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            widget.action.name,
                            style: theme.textTheme.titleMedium?.copyWith(
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                        if (!_isSubmitting)
                          IconButton(
                            onPressed: () => Navigator.pop(context),
                            icon: const Icon(Icons.close),
                            tooltip: 'Закрыть',
                          ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    _isSubmitting
                        ? _buildProgress(theme)
                        : _buildForm(theme, colorScheme),
                  ],
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildProgress(ThemeData theme) {
    return const Padding(
      padding: EdgeInsets.symmetric(vertical: 32),
      child: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            CircularProgressIndicator(),
            SizedBox(height: 16),
            Text('Выполняем действие…'),
          ],
        ),
      ),
    );
  }

  Widget _buildForm(ThemeData theme, ColorScheme colorScheme) {
    switch (widget.action.kind) {
      case _ActionKind.evidenceAnalysis:
        return _buildEvidenceForm(theme, colorScheme);
      case _ActionKind.suspectAction:
        return _buildCharacterForm(theme, colorScheme);
      case _ActionKind.camera:
        return _buildCameraForm(theme, colorScheme);
      case _ActionKind.alibi:
        return _buildAlibiForm(theme, colorScheme);
    }
  }

  Widget _buildEvidenceForm(ThemeData theme, ColorScheme colorScheme) {
    return FutureBuilder<List<Evidence>>(
      future: _evidenceFuture,
      builder: (_, snapshot) {
        if (!snapshot.hasData) {
          return const Center(
            child: Padding(
              padding: EdgeInsets.all(24),
              child: CircularProgressIndicator(),
            ),
          );
        }
        final evidence = snapshot.data!;
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              'Выберите улику для анализа',
              style: theme.textTheme.bodySmall?.copyWith(
                color: colorScheme.onSurface.withAlpha(140),
              ),
            ),
            const SizedBox(height: 12),
            _buildSelectableList(
              evidence
                  .map(
                    (item) => (
                      id: item.id,
                      title: item.name,
                      subtitle: item.description,
                    ),
                  )
                  .toList(),
              _evidenceID,
              (id) => setState(() => _evidenceID = id),
            ),
            const SizedBox(height: 12),
            _submitButton(_evidenceID != null),
          ],
        );
      },
    );
  }

  Widget _buildCharacterForm(ThemeData theme, ColorScheme colorScheme) {
    return FutureBuilder<List<Character>>(
      future: _charactersFuture,
      builder: (_, snapshot) {
        if (!snapshot.hasData) {
          return const Center(
            child: Padding(
              padding: EdgeInsets.all(24),
              child: CircularProgressIndicator(),
            ),
          );
        }
        final characters = snapshot.data!;
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              'Выберите подозреваемого',
              style: theme.textTheme.bodySmall?.copyWith(
                color: colorScheme.onSurface.withAlpha(140),
              ),
            ),
            const SizedBox(height: 12),
            _buildSelectableList(
              characters
                  .map(
                    (item) => (
                      id: item.id,
                      title: item.name,
                      subtitle: item.profession,
                    ),
                  )
                  .toList(),
              _characterID,
              (id) => setState(() => _characterID = id),
            ),
            const SizedBox(height: 12),
            _submitButton(_characterID != null),
          ],
        );
      },
    );
  }

  Widget _buildCameraForm(ThemeData theme, ColorScheme colorScheme) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisSize: MainAxisSize.min,
      children: [
        Text(
          'Будут запрошены записи с камер наблюдения за вечер предполагаемого преступления.',
          style: theme.textTheme.bodyMedium,
        ),
        const SizedBox(height: 8),
        Text(
          'Стоимость: ${widget.action.cost} AP',
          style: theme.textTheme.bodySmall?.copyWith(
            color: colorScheme.primary,
          ),
        ),
        const SizedBox(height: 16),
        _submitButton(true),
      ],
    );
  }

  Widget _buildAlibiForm(ThemeData theme, ColorScheme colorScheme) {
    return FutureBuilder<List<Character>>(
      future: _charactersFuture,
      builder: (_, snapshot) {
        if (!snapshot.hasData) {
          return const Center(
            child: Padding(
              padding: EdgeInsets.all(24),
              child: CircularProgressIndicator(),
            ),
          );
        }
        final characters = snapshot.data!;
        final canSubmit =
            _characterID != null && _alibiController.text.trim().isNotEmpty;
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              'Выберите подозреваемого и опишите алиби для проверки.',
              style: theme.textTheme.bodyMedium,
            ),
            const SizedBox(height: 12),
            _buildSelectableList(
              characters
                  .map(
                    (item) => (
                      id: item.id,
                      title: item.name,
                      subtitle: item.profession,
                    ),
                  )
                  .toList(),
              _characterID,
              (id) => setState(() => _characterID = id),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _alibiController,
              maxLines: 4,
              onEditingComplete: () => FocusScope.of(context).unfocus(),
              onChanged: (_) => setState(() {}),
              scrollPadding: EdgeInsets.symmetric(vertical: 100),
              decoration: InputDecoration(
                hintText: 'Опишите алиби, которое нужно проверить...',
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
            ),
            const SizedBox(height: 16),
            _submitButton(canSubmit),
          ],
        );
      },
    );
  }

  Widget _buildSelectableList(
    List<({String id, String title, String subtitle})> items,
    String? selectedID,
    ValueChanged<String> onSelected,
  ) {
    return Column(
      children:
          items
              .map(
                (item) => ListTile(
                  title: Text(item.title),
                  subtitle: Text(
                    item.subtitle,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  leading: Icon(
                    selectedID == item.id
                        ? Icons.radio_button_checked
                        : Icons.radio_button_unchecked,
                  ),
                  onTap: () => onSelected(item.id),
                  dense: true,
                ),
              )
              .toList(),
    );
  }

  Widget _submitButton(bool enabled) {
    return SizedBox(
      width: double.infinity,
      child: FilledButton(
        onPressed: enabled ? _submit : null,
        child: Text('Заказать (${widget.action.cost} AP)'),
      ),
    );
  }

  Widget _buildReport(
    ThemeData theme,
    ColorScheme colorScheme,
    ActionReport report,
  ) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisSize: MainAxisSize.min,
      children: [
        Row(
          children: [
            Icon(Icons.check_circle, color: colorScheme.primary),
            const SizedBox(width: 8),
            Expanded(
              child: Text(
                report.title,
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Text(
          report.description,
          style: theme.textTheme.bodySmall?.copyWith(
            color: colorScheme.onSurface.withAlpha(140),
          ),
        ),
        const Divider(height: 32),
        MarkdownBody(data: report.body),
        const SizedBox(height: 16),
        SizedBox(
          width: double.infinity,
          child: FilledButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Готово'),
          ),
        ),
      ],
    );
  }
}
