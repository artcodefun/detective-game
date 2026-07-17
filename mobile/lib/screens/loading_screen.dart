import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../blocs/session_cubit.dart';
import '../models/chronology_entry.dart';
import '../services/mock_llm_service.dart';
import '../services/scenario_generator.dart';
import 'desk_screen.dart';

class LoadingScreen extends StatefulWidget {
  const LoadingScreen({super.key});

  @override
  State<LoadingScreen> createState() => _LoadingScreenState();
}

class _LoadingScreenState extends State<LoadingScreen> {
  final _generator = ScenarioGenerator(MockLlmService());
  String _status = 'Готовим сценарий...';

  @override
  void initState() {
    super.initState();
    _generate();
  }

  Future<void> _generate() async {
    setState(() => _status = 'Собираем улики, опрашиваем свидетелей...');

    final characters = MockLlmService.characterPool;
    final session = await _generator.generate(characters);

    if (!mounted) return;

    final started = ChronologyEntry(
      id: 'chron_start_${DateTime.now().millisecondsSinceEpoch}',
      eventType: ChronologyEventType.caseStarted,
      title: 'Дело №${session.id} открыто',
      timestamp: DateTime.now(),
    );
    final enriched = session.addChronologyEntry(started);

    context.read<SessionCubit>().newGame(enriched);
    Navigator.pushReplacement(
      context,
      MaterialPageRoute(builder: (_) => const DeskScreen()),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Новое дело'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 32),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                Icons.auto_awesome,
                size: 64,
                color: colorScheme.primary,
              ),
              const SizedBox(height: 24),
              Text(
                _status,
                style: theme.textTheme.titleMedium,
              ),
              const SizedBox(height: 24),
              const LinearProgressIndicator(),
              const SizedBox(height: 16),
              Text(
                'Пожалуйста, подождите...',
                style: theme.textTheme.bodySmall?.copyWith(
                  color: colorScheme.onSurface.withAlpha(120),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
