import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../models/chronology_entry.dart';
import '../models/notebook.dart';
import '../services/api_service.dart';

class NotebookScreen extends StatefulWidget {
  const NotebookScreen({super.key});

  @override
  State<NotebookScreen> createState() => _NotebookScreenState();
}

class _NotebookScreenState extends State<NotebookScreen> {
  late Future<List<ChronologyEntry>> _future;

  @override
  void initState() {
    super.initState();
    _future = context.read<ApiService>().getChronology();
  }

  Future<void> _reload() async {
    setState(() => _future = context.read<ApiService>().getChronology());
  }

  List<(ChronologyEntry, NotebookEntry)> _taggedNotes(List<ChronologyEntry> chronology) {
    final result = <(ChronologyEntry, NotebookEntry)>[];
    for (final chron in chronology) {
      for (final detail in chron.details) {
        if (detail.userTags.isNotEmpty || (detail.userNote != null && detail.userNote!.isNotEmpty)) {
          result.add((chron, detail));
        }
      }
    }
    return result;
  }

  void _openTagSheet(ChronologyEntry chron, NotebookEntry entry) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(20))),
      builder:
          (_) => _TagNoteSheet(
            entry: entry,
            onApply: (tags, note) async {
              try {
                await context.read<ApiService>().updateNotebookEntry(
                  chronId: chron.id,
                  noteId: entry.id,
                  tags: tags.map((t) => t.name).toList(),
                  note: note,
                );
                _reload();
              } catch (e) {
                if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Ошибка: $e')));
              }
            },
            onClear: () async {
              try {
                await context.read<ApiService>().updateNotebookEntry(chronId: chron.id, noteId: entry.id, tags: []);
                _reload();
              } catch (e) {
                if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Ошибка: $e')));
              }
            },
          ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return DefaultTabController(
      length: 2,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Блокнот'),
          bottom: const TabBar(tabs: [Tab(text: 'Хронология'), Tab(text: 'Заметки')]),
        ),
        body: FutureBuilder<List<ChronologyEntry>>(
          future: _future,
          builder: (_, snapshot) {
            if (snapshot.connectionState != ConnectionState.done) {
              return const Center(child: CircularProgressIndicator());
            }
            if (snapshot.hasError) {
              return Center(child: Text('Ошибка: ${snapshot.error}'));
            }
            final chronology = snapshot.data!;
            return TabBarView(children: [_buildTimeline(theme, chronology), _buildNotes(theme, chronology)]);
          },
        ),
      ),
    );
  }

  Widget _buildTimeline(ThemeData theme, List<ChronologyEntry> chronology) {
    if (chronology.isEmpty) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.history, size: 64, color: theme.colorScheme.onSurface.withAlpha(80)),
            const SizedBox(height: 16),
            Text(
              'Хронология пуста',
              style: theme.textTheme.bodyLarge?.copyWith(color: theme.colorScheme.onSurface.withAlpha(120)),
            ),
          ],
        ),
      );
    }
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: chronology.length,
      itemBuilder: (_, index) {
        final chron = chronology[index];
        return _ChronologyCard(entry: chron, onDetailTap: (detail) => _openTagSheet(chron, detail));
      },
    );
  }

  Widget _buildNotes(ThemeData theme, List<ChronologyEntry> chronology) {
    final tagged = _taggedNotes(chronology);
    if (tagged.isEmpty) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.flag_outlined, size: 64, color: theme.colorScheme.onSurface.withAlpha(80)),
            const SizedBox(height: 16),
            Text(
              'Нет помеченных записей',
              style: theme.textTheme.bodyLarge?.copyWith(color: theme.colorScheme.onSurface.withAlpha(120)),
            ),
          ],
        ),
      );
    }
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: tagged.length,
      itemBuilder: (_, index) {
        final (chron, entry) = tagged[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 8),
          child: InkWell(
            onTap: () => _openTagSheet(chron, entry),
            borderRadius: BorderRadius.circular(12),
            child: Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      _tagChips(entry.userTags),
                      const Spacer(),
                      Text(
                        _formatTime(entry.timestamp),
                        style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurface.withAlpha(120)),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Text(entry.description),
                  const SizedBox(height: 4),
                  Text(
                    '— ${chron.title}',
                    style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurface.withAlpha(100)),
                  ),
                  if (entry.userNote != null && entry.userNote!.isNotEmpty) ...[
                    const SizedBox(height: 8),
                    Container(
                      width: double.infinity,
                      padding: const EdgeInsets.all(8),
                      decoration: BoxDecoration(
                        color: theme.colorScheme.surfaceContainerHighest,
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Text(
                        entry.userNote!,
                        style: theme.textTheme.bodySmall?.copyWith(fontStyle: FontStyle.italic),
                      ),
                    ),
                  ],
                ],
              ),
            ),
          ),
        );
      },
    );
  }

  Widget _tagChips(List<NoteTag> tags) {
    return Wrap(
      spacing: 4,
      children:
          tags.map((t) {
            Color color;
            switch (t) {
              case NoteTag.strange:
                color = Colors.purple;
              case NoteTag.suspicious:
                color = Colors.orange;
              case NoteTag.lie:
                color = Colors.red;
              case NoteTag.key:
                color = Colors.blue;
            }
            return Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
              decoration: BoxDecoration(
                color: color.withAlpha(30),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: color.withAlpha(80)),
              ),
              child: Text(_tagLabel(t), style: TextStyle(fontSize: 11, color: color)),
            );
          }).toList(),
    );
  }

  static String _tagLabel(NoteTag tag) {
    switch (tag) {
      case NoteTag.strange:
        return 'странно';
      case NoteTag.suspicious:
        return 'подозрительно';
      case NoteTag.lie:
        return 'ложь';
      case NoteTag.key:
        return 'ключевое';
    }
  }

  String _formatTime(DateTime dt) => '${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
}

class _ChronologyCard extends StatefulWidget {
  final ChronologyEntry entry;
  final void Function(NotebookEntry detail) onDetailTap;

  const _ChronologyCard({required this.entry, required this.onDetailTap});

  @override
  State<_ChronologyCard> createState() => _ChronologyCardState();
}

class _ChronologyCardState extends State<_ChronologyCard> {
  bool _expanded = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final entry = widget.entry;
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: Column(
        children: [
          InkWell(
            onTap: () => setState(() => _expanded = !_expanded),
            borderRadius: _expanded ? const BorderRadius.vertical(top: Radius.circular(12)) : BorderRadius.circular(12),
            child: Padding(
              padding: const EdgeInsets.all(12),
              child: Row(
                children: [
                  Container(
                    width: 48,
                    height: 48,
                    decoration: BoxDecoration(
                      color: colorScheme.primaryContainer,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Center(
                      child: Icon(_iconForEventType(entry.eventType), color: colorScheme.onPrimaryContainer, size: 20),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(entry.title, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
                        const SizedBox(height: 2),
                        Text(
                          '${entry.eventTypeLabel} • ${_formatTime(entry.timestamp)}',
                          style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(120)),
                        ),
                      ],
                    ),
                  ),
                  AnimatedRotation(
                    turns: _expanded ? 0.5 : 0,
                    duration: const Duration(milliseconds: 200),
                    child: Icon(Icons.expand_more, color: colorScheme.onSurface.withAlpha(120)),
                  ),
                ],
              ),
            ),
          ),
          AnimatedCrossFade(
            firstChild: const SizedBox.shrink(),
            secondChild: _buildDetails(theme, colorScheme),
            crossFadeState: _expanded ? CrossFadeState.showSecond : CrossFadeState.showFirst,
            duration: const Duration(milliseconds: 200),
          ),
        ],
      ),
    );
  }

  Widget _buildDetails(ThemeData theme, ColorScheme colorScheme) {
    final details = widget.entry.details;
    if (details.isEmpty) {
      return Padding(
        padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
        child: Text(
          'Нет записей',
          style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(100)),
        ),
      );
    }
    return Column(
      children: [
        const Divider(height: 1),
        ...details.map((detail) => _DetailTile(entry: detail, onTap: () => widget.onDetailTap(detail))),
      ],
    );
  }

  IconData _iconForEventType(ChronologyEventType type) {
    switch (type) {
      case ChronologyEventType.caseStarted:
        return Icons.auto_awesome;
      case ChronologyEventType.interrogation:
        return Icons.chat;
      case ChronologyEventType.labAnalysis:
        return Icons.science;
      case ChronologyEventType.alibiCheck:
        return Icons.access_time;
      case ChronologyEventType.cameraReview:
        return Icons.videocam;
      case ChronologyEventType.transactionCheck:
        return Icons.account_balance;
    }
  }

  String _formatTime(DateTime dt) => '${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
}

class _DetailTile extends StatelessWidget {
  final NotebookEntry entry;
  final VoidCallback onTap;

  const _DetailTile({required this.entry, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final hasTags = entry.userTags.isNotEmpty;
    final hasNote = entry.userNote != null && entry.userNote!.isNotEmpty;
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (hasTags)
              Padding(
                padding: const EdgeInsets.only(top: 6, right: 8),
                child: Icon(Icons.flag, size: 14, color: theme.colorScheme.primary),
              ),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(entry.description, style: theme.textTheme.bodyMedium),
                  if (hasTags) ...[const SizedBox(height: 4), _TagLabel(entry: entry)],
                  if (hasNote) ...[
                    const SizedBox(height: 4),
                    Text(
                      entry.userNote!,
                      style: theme.textTheme.bodySmall?.copyWith(
                        fontStyle: FontStyle.italic,
                        color: theme.colorScheme.onSurface.withAlpha(160),
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _TagLabel extends StatelessWidget {
  final NotebookEntry entry;

  const _TagLabel({required this.entry});

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 4,
      children:
          entry.userTags.map((t) {
            Color color;
            switch (t) {
              case NoteTag.strange:
                color = Colors.purple;
              case NoteTag.suspicious:
                color = Colors.orange;
              case NoteTag.lie:
                color = Colors.red;
              case NoteTag.key:
                color = Colors.blue;
            }
            return Container(
              padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
              decoration: BoxDecoration(
                color: color.withAlpha(25),
                borderRadius: BorderRadius.circular(10),
                border: Border.all(color: color.withAlpha(60)),
              ),
              child: Text(entry.tagsLabel, style: TextStyle(fontSize: 10, color: color, fontWeight: FontWeight.w500)),
            );
          }).toList(),
    );
  }
}

class _TagNoteSheet extends StatefulWidget {
  final NotebookEntry entry;
  final void Function(List<NoteTag> tags, String? note) onApply;
  final VoidCallback onClear;

  const _TagNoteSheet({required this.entry, required this.onApply, required this.onClear});

  @override
  State<_TagNoteSheet> createState() => _TagNoteSheetState();
}

class _TagNoteSheetState extends State<_TagNoteSheet> {
  late List<NoteTag> _selectedTags;
  late TextEditingController _noteController;

  @override
  void initState() {
    super.initState();
    _selectedTags = List.from(widget.entry.userTags);
    _noteController = TextEditingController(text: widget.entry.userNote ?? '');
  }

  @override
  void dispose() {
    _noteController.dispose();
    super.dispose();
  }

  void _toggleTag(NoteTag tag) {
    setState(() {
      if (_selectedTags.contains(tag)) {
        _selectedTags.remove(tag);
      } else {
        _selectedTags.add(tag);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final bottom = MediaQuery.of(context).viewInsets.bottom;
    return Padding(
      padding: EdgeInsets.fromLTRB(20, 12, 20, 12 + bottom),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Center(
            child: Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: theme.colorScheme.onSurface.withAlpha(60),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
          const SizedBox(height: 16),
          Text('Пометить запись', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 4),
          Text(
            widget.entry.description,
            style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurface.withAlpha(140)),
          ),
          const SizedBox(height: 16),
          Text('Метки', style: theme.textTheme.labelLarge),
          const SizedBox(height: 8),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children:
                NoteTag.values.map((tag) {
                  final selected = _selectedTags.contains(tag);
                  Color color;
                  String label;
                  switch (tag) {
                    case NoteTag.strange:
                      color = Colors.purple;
                      label = 'Странно';
                    case NoteTag.suspicious:
                      color = Colors.orange;
                      label = 'Подозрительно';
                    case NoteTag.lie:
                      color = Colors.red;
                      label = 'Ложь';
                    case NoteTag.key:
                      color = Colors.blue;
                      label = 'Ключевое';
                  }
                  return FilterChip(
                    label: Text(label),
                    selected: selected,
                    selectedColor: color.withAlpha(40),
                    checkmarkColor: color,
                    onSelected: (_) => _toggleTag(tag),
                  );
                }).toList(),
          ),
          const SizedBox(height: 16),
          Text('Заметка', style: theme.textTheme.labelLarge),
          const SizedBox(height: 8),
          TextField(
            controller: _noteController,
            maxLines: 3,
            minLines: 1,
            decoration: InputDecoration(
              hintText: 'Ваш комментарий...',
              border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
              contentPadding: const EdgeInsets.all(12),
            ),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              OutlinedButton(
                onPressed: () {
                  widget.onClear();
                  Navigator.pop(context);
                },
                child: const Text('Очистить'),
              ),
              const Spacer(),
              OutlinedButton(onPressed: () => Navigator.pop(context), child: const Text('Отмена')),
              const SizedBox(width: 8),
              FilledButton(
                onPressed: () {
                  widget.onApply(
                    _selectedTags,
                    _noteController.text.trim().isEmpty ? null : _noteController.text.trim(),
                  );
                  Navigator.pop(context);
                },
                child: const Text('Применить'),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
