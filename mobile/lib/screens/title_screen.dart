import 'dart:async';
import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../blocs/session_cubit.dart';
import '../models/scenario.dart';
import '../services/api_service.dart';
import 'desk_screen.dart';
import 'loading_screen.dart';
import 'previous_cases_screen.dart';
import 'settings_screen.dart';

class TitleScreen extends StatefulWidget {
  const TitleScreen({super.key});

  @override
  State<TitleScreen> createState() => _TitleScreenState();
}

class _TitleScreenState extends State<TitleScreen> {
  Session? _activeSession;
  bool _checking = true;

  @override
  void initState() {
    super.initState();
    _checkActiveSession();
  }

  Future<void> _checkActiveSession() async {
    try {
      final session = await context.read<ApiService>().getCurrentSession();
      if (mounted) setState(() => _activeSession = session);
    } catch (_) {
      // no active session
    }
    if (mounted) setState(() => _checking = false);
  }

  void _continueCase() {
    if (_activeSession == null) return;
    final s = _activeSession!;
    context.read<SessionCubit>().resumeSession(s.id, s.caseName, s.actionPoints, s.phase);
    Navigator.push(context, MaterialPageRoute(builder: (_) => const DeskScreen()));
  }

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
              Icon(Icons.search, size: 80, color: colorScheme.primary),
              const SizedBox(height: 16),
              Text('ДетектИИв', style: theme.textTheme.headlineLarge?.copyWith(fontWeight: FontWeight.bold, letterSpacing: 2)),
              const SizedBox(height: 4),
              const _CaseNumberAnimation(),
              const SizedBox(height: 48),
              if (_activeSession != null || _checking)
                _TitleButton(
                  label: _checking ? 'Проверяем...' : 'Продолжить дело',
                  icon: _checking ? Icons.hourglass_top : Icons.play_arrow,
                  onPressed: _checking ? null : _continueCase,
                ),
              if (_activeSession != null || _checking) const SizedBox(height: 12),
              _TitleButton(
                label: 'Новое дело',
                icon: Icons.folder_open,
                onPressed: _checking
                    ? null
                    : () {
                        Navigator.push(context, MaterialPageRoute(builder: (_) => const LoadingScreen()));
                      },
              ),
              const SizedBox(height: 12),
              _TitleButton(
                label: 'Предыдущие дела',
                icon: Icons.history,
                onPressed: _checking
                    ? null
                    : () {
                        Navigator.push(context, MaterialPageRoute(builder: (_) => const PreviousCasesScreen()));
                      },
              ),
              const SizedBox(height: 12),
              _TitleButton(
                label: 'Настройки',
                icon: Icons.settings,
                onPressed: _checking
                    ? null
                    : () {
                        Navigator.push(context, MaterialPageRoute(builder: (_) => const SettingsScreen()));
                      },
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _TitleButton extends StatelessWidget {
  final String label;
  final IconData icon;
  final VoidCallback? onPressed;

  const _TitleButton({required this.label, required this.icon, required this.onPressed});

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
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
          textStyle: const TextStyle(fontSize: 16),
        ),
      ),
    );
  }
}

const _scrambleChars = r'#$%&@!?*+=~<>/\![]{}^|';
const _scrambleLen = 4;

class _CaseNumberAnimation extends StatefulWidget {
  const _CaseNumberAnimation();

  @override
  State<_CaseNumberAnimation> createState() => _CaseNumberAnimationState();
}

class _CaseNumberAnimationState extends State<_CaseNumberAnimation> {
  final _random = Random();
  int _ticks = 0;
  int _locked = -1;
  Timer? _timer;
  late final String _number;

  static const _ticksPerChar = 8;

  @override
  void initState() {
    super.initState();
    _number = (1000 + _random.nextInt(9000)).toString();
    _timer = Timer.periodic(const Duration(milliseconds: 80), (_) => _scramble());
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  void _scramble() {
    _ticks++;
    _locked = _ticks ~/ _ticksPerChar;
    if (_locked >= _scrambleLen) {
      _timer?.cancel();
      _locked = _scrambleLen;
    }
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final chars = List.generate(_scrambleLen, (i) {
      if (i < _locked) return _number[i];
      return _scrambleChars[_random.nextInt(_scrambleChars.length)];
    });

    return Text(
      'Дело №${chars.join()}',
      style: theme.textTheme.titleMedium?.copyWith(
        color: theme.colorScheme.onSurface.withAlpha(120),
        letterSpacing: 4,
      ),
    );
  }
}
