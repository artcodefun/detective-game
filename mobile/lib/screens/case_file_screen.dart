import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter_markdown/flutter_markdown.dart';

import '../models/action_report.dart';
import '../models/scenario.dart';
import '../services/api_service.dart';
import '../widgets/load_error_view.dart';
import 'document_screen.dart';

class CaseFileScreen extends StatelessWidget {
  const CaseFileScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 2,
      child: Scaffold(
        appBar: AppBar(title: const Text('Дело'), bottom: const TabBar(tabs: [Tab(text: 'Факты'), Tab(text: 'Улики')])),
        body: const TabBarView(children: [_FactsTab(), _EvidenceTab()]),
      ),
    );
  }
}

class _FactsTab extends StatefulWidget {
  const _FactsTab();

  @override
  State<_FactsTab> createState() => _FactsTabState();
}

class _FactsTabState extends State<_FactsTab> {
  late Future<String> _future;

  @override
  void initState() {
    super.initState();
    _future = _loadCaseBrief();
  }

  Future<String> _loadCaseBrief() {
    return context.read<ApiService>().getCurrentSession().then((s) => s.caseBrief);
  }

  void _retryLoadCaseBrief() {
    setState(() {
      _future = _loadCaseBrief();
    });
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<String>(
      future: _future,
      builder: (_, snapshot) {
        if (snapshot.connectionState != ConnectionState.done) {
          return const Center(child: CircularProgressIndicator());
        }
        if (snapshot.hasError) {
          return LoadErrorView(message: 'Не удалось загрузить материалы дела', onRetry: _retryLoadCaseBrief);
        }
        if (snapshot.data == null || snapshot.data!.isEmpty) {
          return Center(
            child: Text(
              'Документ ещё не готов',
              style: Theme.of(
                context,
              ).textTheme.bodyLarge?.copyWith(color: Theme.of(context).colorScheme.onSurface.withAlpha(120)),
            ),
          );
        }
        return Container(
          color: const Color(0xFFF5F0E8),
          child: Markdown(
            data: snapshot.data!,
            padding: const EdgeInsets.all(16),
            styleSheet: MarkdownStyleSheet(
              h1: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black87),
              h2: const TextStyle(fontSize: 17, fontWeight: FontWeight.bold, color: Colors.black87),
              h3: const TextStyle(fontSize: 15, fontWeight: FontWeight.w600, color: Colors.black87),
              p: const TextStyle(fontSize: 15, color: Colors.black87, height: 1.6),
              strong: const TextStyle(fontWeight: FontWeight.bold, color: Colors.black87),
              em: const TextStyle(fontStyle: FontStyle.italic, color: Colors.black87),
              listBullet: const TextStyle(fontSize: 15, color: Colors.black87),
              tableBody: const TextStyle(fontSize: 14, color: Colors.black87),
              tableHead: const TextStyle(fontSize: 14, fontWeight: FontWeight.bold, color: Colors.black87),
              tableBorder: TableBorder.all(color: Colors.grey.shade400, width: 0.5),
              tableColumnWidth: const FlexColumnWidth(),
              tableCellsPadding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              code: TextStyle(fontSize: 13, color: Colors.grey.shade800, backgroundColor: Colors.grey.shade200),
            ),
          ),
        );
      },
    );
  }
}

sealed class _Item {}

class _EvidenceItem extends _Item {
  final Evidence evidence;

  _EvidenceItem(this.evidence);
}

class _ReportItem extends _Item {
  final ActionReport report;

  _ReportItem(this.report);
}

class _EvidenceTab extends StatefulWidget {
  const _EvidenceTab();

  @override
  State<_EvidenceTab> createState() => _EvidenceTabState();
}

class _EvidenceTabState extends State<_EvidenceTab> {
  late Future<List<_Item>> _future;

  @override
  void initState() {
    super.initState();
    _future = _fetchEvidence();
  }

  Future<List<_Item>> _fetchEvidence() async {
    final api = context.read<ApiService>();
    final evidence = await api.listEvidence();
    final reports = await api.listReports();
    return [...reports.map(_ReportItem.new), ...evidence.map(_EvidenceItem.new)];
  }

  void _retryLoadEvidence() {
    setState(() {
      _future = _fetchEvidence();
    });
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<List<_Item>>(
      future: _future,
      builder: (_, snapshot) {
        if (snapshot.connectionState != ConnectionState.done) {
          return const Center(child: CircularProgressIndicator());
        }
        if (snapshot.hasError) {
          return LoadErrorView(message: 'Не удалось загрузить улики и отчёты', onRetry: _retryLoadEvidence);
        }
        final items = snapshot.data!;
        if (items.isEmpty) {
          return Center(
            child: Text(
              'Улик и отчётов пока нет',
              style: Theme.of(
                context,
              ).textTheme.bodyLarge?.copyWith(color: Theme.of(context).colorScheme.onSurface.withAlpha(120)),
            ),
          );
        }
        return ListView.builder(
          padding: const EdgeInsets.all(16),
          itemCount: items.length,
          itemBuilder: (_, index) {
            final item = items[index];
            if (item is _EvidenceItem) {
              return _EvidenceCard(evidence: item.evidence);
            }
            return _ReportCard(report: (item as _ReportItem).report);
          },
        );
      },
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
                    Text(evidence.name, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
                    const SizedBox(height: 2),
                    Text(
                      evidence.description,
                      style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(140)),
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
                Text(evidence.name, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                const SizedBox(height: 12),
                Text(evidence.detailedDescription, style: theme.textTheme.bodyMedium),
              ],
            ),
          ),
    );
  }
}

class _ReportCard extends StatelessWidget {
  final ActionReport report;

  const _ReportCard({required this.report});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: InkWell(
        onTap: () {
          Navigator.push(
            context,
            MaterialPageRoute(builder: (_) => DocumentScreen(title: report.title, body: report.body)),
          );
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Row(
            children: [
              Icon(Icons.description, size: 28, color: colorScheme.primary),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(report.title, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
                    const SizedBox(height: 2),
                    Text(
                      report.description,
                      style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(140)),
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
}
