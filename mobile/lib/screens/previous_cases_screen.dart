import 'package:flutter/material.dart';

import '../models/report.dart';
import 'results_screen.dart';

class PreviousCasesScreen extends StatelessWidget {
  const PreviousCasesScreen({super.key});

  static final _cases = [
    _CaseData(
      number: '45',
      status: 'Закрыто',
      trackRecord: '85%',
      date: '12.03.2026',
      description: 'Ограбление банка на Лиговском',
      report: FinalReport(
        who: 'Сергей Волков',
        why: 'Был должен крупную сумму криминальному авторитету',
        how: 'Взлом сейфа и подмена сигнализации',
        when: '22:30',
        evidence: 'Отпечатки на сейфе, следы взломщика, запись с камер',
      ),
      breakdown: ScoreBreakdown(whoCorrect: true, whyCorrect: true, howCorrect: true, whenCorrect: true, evidenceCorrect: false),
    ),
    _CaseData(
      number: '38',
      status: 'Закрыто',
      trackRecord: '72%',
      date: '28.01.2026',
      description: 'Исчезновение коллекционера',
      report: FinalReport(
        who: 'Артем Белов',
        why: 'Хотел завладеть коллекцией картин',
        how: 'Усыпил и вывез на машине',
        when: 'около 23:00',
        evidence: 'Волокна ткани из багажника, показания соседей',
      ),
      breakdown: ScoreBreakdown(whoCorrect: true, whyCorrect: true, howCorrect: true, whenCorrect: false, evidenceCorrect: false),
    ),
    _CaseData(
      number: '22',
      status: 'Закрыто',
      trackRecord: '91%',
      date: '15.11.2025',
      description: 'Поджог в особняке',
      report: FinalReport(
        who: 'Илья Морозов',
        why: 'Страховка и месть владельцу',
        how: 'Залил бензином и поджёг через фитиль',
        when: '03:15',
        evidence: 'Канистра с бензином, окурок, запись с камеры',
      ),
      breakdown: ScoreBreakdown(whoCorrect: true, whyCorrect: true, howCorrect: true, whenCorrect: true, evidenceCorrect: true),
    ),
    _CaseData(
      number: '17',
      status: 'Архив',
      trackRecord: '—',
      date: '03.09.2025',
      description: 'Серия квартирных краж',
      report: null,
      breakdown: null,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      appBar: AppBar(title: const Text('Предыдущие дела')),
      body: ListView.separated(
        padding: const EdgeInsets.all(16),
        itemCount: _cases.length,
        separatorBuilder: (_, __) => const Divider(height: 1),
        itemBuilder: (_, index) {
          final c = _cases[index];
          return ListTile(
            leading: CircleAvatar(
              backgroundColor: colorScheme.primaryContainer,
              child: Text(
                c.number,
                style: TextStyle(fontSize: 12, fontWeight: FontWeight.bold, color: colorScheme.onPrimaryContainer),
              ),
            ),
            title: Text('Дело №${c.number}', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
            subtitle: Text('${c.date} • ${c.description}', maxLines: 1, overflow: TextOverflow.ellipsis),
            trailing: Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
              decoration: BoxDecoration(
                color: c.trackRecord == '—'
                    ? colorScheme.surfaceContainerHighest
                    : int.parse(c.trackRecord.replaceAll('%', '')) >= 80
                        ? Colors.green.withAlpha(25)
                        : Colors.orange.withAlpha(25),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Text(
                c.trackRecord,
                style: TextStyle(
                  fontSize: 11,
                  color: c.trackRecord == '—'
                      ? colorScheme.onSurface.withAlpha(100)
                      : int.parse(c.trackRecord.replaceAll('%', '')) >= 80
                          ? Colors.green
                          : Colors.orange,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
            onTap: () {
              if (c.report == null || c.breakdown == null) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('Дело №${c.number} — ${c.status}. ${c.description}')),
                );
                return;
              }
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (_) => ResultsScreen(
                    result: GameResult(
                      playerReport: c.report!,
                      breakdown: c.breakdown!,
                      narrativeFeedback: 'Дело №${c.number} — ${c.status}. ${c.description}',
                      breakdownDetails: {},
                    ),
                    playerReport: c.report,
                    showHomeButton: false,
                  ),
                ),
              );
            },
          );
        },
      ),
    );
  }
}

class _CaseData {
  final String number;
  final String status;
  final String trackRecord;
  final String date;
  final String description;
  final FinalReport? report;
  final ScoreBreakdown? breakdown;

  _CaseData({
    required this.number,
    required this.status,
    required this.trackRecord,
    required this.date,
    required this.description,
    this.report,
    this.breakdown,
  });
}
