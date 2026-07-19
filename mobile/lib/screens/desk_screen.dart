import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../blocs/session_cubit.dart';
import '../models/game_state.dart';
import 'actions_screen.dart';
import 'case_file_screen.dart';
import 'interrogation_screen.dart';
import 'notebook_screen.dart';
import 'report_screen.dart';

class DeskScreen extends StatelessWidget {
  const DeskScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final session = context.watch<SessionCubit>().state!;
    return Scaffold(
      appBar: AppBar(
        title: Text('Дело №${session.id}'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.pop(context),
        ),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 12),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(
                  Icons.search,
                  size: 18,
                  color: Theme.of(context).colorScheme.primary,
                ),
                const SizedBox(width: 4),
                Text(
                  '${session.actionPoints}',
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(width: 8),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Expanded(
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Expanded(
                    child:                     _DeskItem(
                      icon: Icons.folder,
                      label: 'Дело',
                      subtitle: 'факты и улики',
                      onTap: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (_) => const CaseFileScreen(),
                          ),
                        );
                      },
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child:                     _DeskItem(
                      icon: Icons.phone,
                      label: 'Телефон',
                      subtitle: 'подозреваемые',
                      onTap: () => _openPhoneBook(context),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 12),
            Expanded(
                  child: _DeskItem(
                    icon: Icons.book,
                    label: 'Блокнот',
                    subtitle: 'хронология и заметки',
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (_) => const NotebookScreen(),
                        ),
                      );
                    },
                  ),
            ),
            const SizedBox(height: 12),
            Expanded(
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Expanded(
                    child: _DeskItem(
                      icon: Icons.build,
                      label: 'Действия',
                      subtitle: 'анализы и запросы',
                      onTap: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (_) => const ActionsScreen(),
                          ),
                        );
                      },
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _DeskItem(
                      icon: Icons.description,
                      label: 'Отчёт',
                      subtitle: 'финальная версия',
                      onTap: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (_) => const ReportScreen(),
                          ),
                        );
                      },
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _openPhoneBook(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (_) => const _PhoneBookSheet(),
    );
  }
}

class _DeskItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final String subtitle;
  final VoidCallback onTap;

  const _DeskItem({
    required this.icon,
    required this.label,
    required this.subtitle,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Card(
      elevation: 2,
      color: colorScheme.surfaceContainerHighest,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(icon, size: 36, color: colorScheme.primary),
              const SizedBox(height: 8),
              Text(label, style: theme.textTheme.titleSmall?.copyWith(
                fontWeight: FontWeight.bold,
              )),
              const SizedBox(height: 2),
              Text(subtitle, style: theme.textTheme.bodySmall?.copyWith(
                color: colorScheme.onSurface.withAlpha(120),
              )),
            ],
          ),
        ),
      ),
    );
  }
}

class _PhoneBookSheet extends StatelessWidget {
  const _PhoneBookSheet();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final session = context.read<SessionCubit>().state!;

    return DraggableScrollableSheet(
      initialChildSize: 0.6,
      minChildSize: 0.4,
      maxChildSize: 0.85,
      expand: false,
      builder: (_, scrollController) {
        return Padding(
          padding: const EdgeInsets.fromLTRB(20, 12, 20, 20),
          child: Column(
            children: [
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: colorScheme.onSurface.withAlpha(60),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 16),
              Text(
                'Телефонная книга',
                style: theme.textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                'Вызовите подозреваемого на допрос',
                style: theme.textTheme.bodySmall?.copyWith(
                  color: colorScheme.onSurface.withAlpha(120),
                ),
              ),
              const SizedBox(height: 16),
              Expanded(
                child: ListView.separated(
                  controller: scrollController,
                  itemCount: session.characters.length,
                  separatorBuilder: (_, __) => const Divider(height: 1),
                  itemBuilder: (_, index) {
                    final char = session.characters[index];
                    return _SuspectTile(characterState: char);
                  },
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}

class _SuspectTile extends StatelessWidget {
  final CharacterState characterState;

  const _SuspectTile({
    required this.characterState,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final char = characterState.base;
    final remaining = characterState.interrogationsRemaining;

    return ListTile(
      leading: CircleAvatar(
        backgroundColor: colorScheme.primaryContainer,
        child: Text(
          char.name[0],
          style: TextStyle(
            color: colorScheme.onPrimaryContainer,
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
      title: Text(char.name),
      subtitle: Text('Допросов: $remaining/3'),
      trailing: IconButton(
        icon: Icon(Icons.call, color: colorScheme.primary),
        onPressed: remaining > 0
            ? () {
                Navigator.pop(context);
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (_) => InterrogationScreen(
                      characterState: characterState,
                    ),
                  ),
                );
              }
            : null,
      ),
    );
  }
}

