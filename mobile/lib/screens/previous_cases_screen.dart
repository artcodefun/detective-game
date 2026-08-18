import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/scenario.dart';
import '../services/api_service.dart';
import 'results_screen.dart';

class PreviousCasesScreen extends StatefulWidget {
  const PreviousCasesScreen({super.key});

  @override
  State<PreviousCasesScreen> createState() => _PreviousCasesScreenState();
}

class _PreviousCasesScreenState extends State<PreviousCasesScreen> {
  late Future<List<Session>> _future;

  @override
  void initState() {
    super.initState();
    _future = context.read<ApiService>().listHistory();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      appBar: AppBar(title: const Text('Предыдущие дела')),
      body: FutureBuilder<List<Session>>(
        future: _future,
        builder: (_, snapshot) {
          if (snapshot.connectionState != ConnectionState.done) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return Center(child: Text('Ошибка: ${snapshot.error}'));
          }
          final sessions = snapshot.data!;

          if (sessions.isEmpty) {
            return Center(
              child: Text(
                'Пока нет завершённых дел',
                style: theme.textTheme.bodyLarge?.copyWith(color: colorScheme.onSurface.withAlpha(120)),
              ),
            );
          }

          return ListView.separated(
            padding: const EdgeInsets.all(16),
            itemCount: sessions.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (_, index) {
              final s = sessions[index];

              int? accuracy;
              String trackRecord = '—';

              if (s.gameResult != null) {
                accuracy = (s.gameResult!.breakdown.accuracy * 100).toInt();
                trackRecord = '$accuracy%';
              }

              return ListTile(
                leading: CircleAvatar(
                  backgroundColor: colorScheme.primaryContainer,
                  child: Text(
                    s.caseName.isNotEmpty ? s.caseName[0].toUpperCase() : '?',
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: colorScheme.onPrimaryContainer),
                  ),
                ),
                title: Text(s.caseName, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
                subtitle: Text(_formatDate(s.createdAt), maxLines: 1, overflow: TextOverflow.ellipsis),
                trailing: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color:
                        accuracy == null
                            ? colorScheme.surfaceContainerHighest
                            : accuracy >= 80
                            ? Colors.green.withAlpha(25)
                            : Colors.orange.withAlpha(25),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    trackRecord,
                    style: TextStyle(
                      fontSize: 11,
                      color:
                          accuracy == null
                              ? colorScheme.onSurface.withAlpha(100)
                              : accuracy >= 80
                              ? Colors.green
                              : Colors.orange,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
                onTap: () {
                  if (s.gameResult == null) return;
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder:
                          (_) => ResultsScreen(
                            result: s.gameResult!,
                            playerReport: s.gameResult!.playerReport,
                            showHomeButton: false,
                          ),
                    ),
                  );
                },
              );
            },
          );
        },
      ),
    );
  }

  String _formatDate(DateTime dt) {
    return '${dt.day.toString().padLeft(2, '0')}.${dt.month.toString().padLeft(2, '0')}.${dt.year}';
  }
}
