import 'package:flutter/material.dart';

import 'loading_screen.dart';
import 'previous_cases_screen.dart';
import 'settings_screen.dart';

class TitleScreen extends StatelessWidget {
  const TitleScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      body: Center(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 32),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                Icons.search,
                size: 80,
                color: colorScheme.primary,
              ),
              const SizedBox(height: 16),
              Text(
                'Детектив',
                style: theme.textTheme.headlineLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                  letterSpacing: 2,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                'Дело №...',
                style: theme.textTheme.titleMedium?.copyWith(
                  color: colorScheme.onSurface.withAlpha(120),
                  letterSpacing: 4,
                ),
              ),
              const SizedBox(height: 48),
              _TitleButton(
                label: 'Новое дело',
                icon: Icons.folder_open,
                onPressed: () => _startNewCase(context),
              ),
              const SizedBox(height: 12),
              _TitleButton(
                label: 'Предыдущие дела',
                icon: Icons.history,
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (_) => const PreviousCasesScreen(),
                    ),
                  );
                },
              ),
              const SizedBox(height: 12),
              _TitleButton(
                label: 'Настройки',
                icon: Icons.settings,
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (_) => const SettingsScreen(),
                    ),
                  );
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _startNewCase(BuildContext context) {
    Navigator.push(
      context,
      MaterialPageRoute(builder: (_) => const LoadingScreen()),
    );
  }
}

class _TitleButton extends StatelessWidget {
  final String label;
  final IconData icon;
  final VoidCallback onPressed;

  const _TitleButton({
    required this.label,
    required this.icon,
    required this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return SizedBox(
      width: 260,
      height: 52,
      child: ElevatedButton.icon(
        onPressed: onPressed,
        icon: Icon(icon),
        label: Text(label),
        style: ElevatedButton.styleFrom(
          foregroundColor: colorScheme.onPrimaryContainer,
          backgroundColor: colorScheme.primaryContainer,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          textStyle: const TextStyle(fontSize: 16),
        ),
      ),
    );
  }
}
