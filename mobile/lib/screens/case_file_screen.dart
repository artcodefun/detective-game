import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../blocs/session_cubit.dart';
import '../models/scenario.dart';

class CaseFileScreen extends StatelessWidget {
  const CaseFileScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 2,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Дело'),
          bottom: const TabBar(
            tabs: [
              Tab(text: 'Факты'),
              Tab(text: 'Улики'),
            ],
          ),
        ),
        body: TabBarView(
          children: [
            _FactsTab(),
            _EvidenceTab(),
          ],
        ),
      ),
    );
  }
}

class _FactsTab extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final session = context.watch<SessionCubit>().state!;
    final crime = session.crime;

    final buffer = StringBuffer()
      ..writeln('Тип преступления: ${crime.typeLabel}')
      ..writeln()
      ..writeln('Жертва: ${crime.victim}')
      ..writeln()
      ..writeln('Время: ${crime.timeOfCrime}')
      ..writeln()
      ..writeln('Обстоятельства:')
      ..writeln('Мотив: ${crime.motive}')
      ..writeln('Способ: ${crime.method}')
      ..writeln()
      ..writeln('Подозреваемые:');
    for (final c in session.characters) {
      buffer.writeln('  • ${c.base.name} — ${c.base.profession}');
    }

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Container(
        decoration: BoxDecoration(
      color: Colors.white,
      borderRadius: BorderRadius.circular(8),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withAlpha(20),
              blurRadius: 8,
              offset: const Offset(0, 2),
            ),
          ],
        ),
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Дело №${session.id}',
              style: const TextStyle(
                fontSize: 22,
                fontWeight: FontWeight.bold,
                color: Colors.black87,
              ),
            ),
            const SizedBox(height: 16),
            Text(
              buffer.toString().trim(),
              style: const TextStyle(
                fontSize: 16,
                color: Colors.black87,
                height: 1.6,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _EvidenceTab extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final session = context.watch<SessionCubit>().state!;
    final evidence = session.evidence;

    if (evidence.isEmpty) {
      return Center(
        child: Text(
          'Улик пока нет',
          style: Theme.of(context).textTheme.bodyLarge?.copyWith(
            color: Theme.of(context).colorScheme.onSurface.withAlpha(120),
          ),
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: evidence.length,
      itemBuilder: (_, index) => _EvidenceCard(evidence: evidence[index]),
    );
  }
}

class _EvidenceCard extends StatelessWidget {
  final Evidence evidence;

  const _EvidenceCard({required this.evidence});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: InkWell(
        onTap: () => _openDetail(context),
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Row(
            children: [
              Icon(Icons.inventory_2, size: 28, color: colorScheme.primary),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            evidence.name,
                            style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold),
                          ),
                        ),
                        if (evidence.analyzed)
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
                            decoration: BoxDecoration(
                              color: colorScheme.primary.withAlpha(30),
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Text(
                              'анализ готов',
                              style: TextStyle(fontSize: 10, color: colorScheme.primary),
                            ),
                          ),
                      ],
                    ),
                    const SizedBox(height: 2),
                    Text(
                      evidence.description,
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: colorScheme.onSurface.withAlpha(140),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 4),
              Icon(Icons.chevron_right, size: 20, color: colorScheme.onSurface.withAlpha(80)),
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
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (_) => Padding(
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
            Text(evidence.name, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 12),
            Text(evidence.detailedDescription, style: theme.textTheme.bodyMedium),
            if (evidence.analyzed && evidence.analysisResult != null) ...[
              const SizedBox(height: 16),
              Text('Результат анализа', style: theme.textTheme.labelLarge?.copyWith(color: colorScheme.primary)),
              const SizedBox(height: 4),
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: colorScheme.surfaceContainerHighest,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Text(evidence.analysisResult!, style: theme.textTheme.bodyMedium),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
