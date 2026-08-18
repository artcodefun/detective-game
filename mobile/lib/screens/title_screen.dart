import 'dart:async';
import 'dart:math';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:url_launcher/url_launcher.dart';

import '../models/scenario.dart';
import '../services/api_service.dart';
import '../services/session_service.dart';
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
  bool _checkingRequest = false;
  bool _modalSheetOpen = false;
  late final ApiService _api;

  @override
  void initState() {
    super.initState();
    _api = context.read<ApiService>();
    _api.addListener(_onApiChanged);
    _onApiChanged();
  }

  @override
  void dispose() {
    _api.removeListener(_onApiChanged);
    super.dispose();
  }

  void _onApiChanged() {
    if (_api.initializationStatus == InitializationStatus.ready && _checking && !_checkingRequest) {
      _checkActiveSession();
    }
    if ((_api.initializationStatus == InitializationStatus.versionCheckFailed ||
            _api.initializationStatus == InitializationStatus.registrationFailed) &&
        !_modalSheetOpen &&
        mounted) {
      _modalSheetOpen = true;
      WidgetsBinding.instance.addPostFrameCallback((_) => _showInitializationError());
    }
    if (_api.initializationStatus == InitializationStatus.updateRequired && !_modalSheetOpen && mounted) {
      _modalSheetOpen = true;
      WidgetsBinding.instance.addPostFrameCallback((_) => _showUpdateRequired());
    }
    if (mounted) setState(() {});
  }

  Future<void> _showInitializationError() async {
    await showModalBottomSheet<void>(
      context: context,
      isDismissible: false,
      enableDrag: false,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(20))),
      builder:
          (sheetContext) => Padding(
            padding: const EdgeInsets.fromLTRB(20, 12, 20, 20),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Container(
                  width: 40,
                  height: 4,
                  decoration: BoxDecoration(
                    color: Theme.of(sheetContext).colorScheme.onSurface.withAlpha(60),
                    borderRadius: BorderRadius.circular(2),
                  ),
                ),
                const SizedBox(height: 20),
                const Icon(Icons.cloud_off, size: 36),
                const SizedBox(height: 12),
                Text('Не удалось подключиться к серверу', style: Theme.of(sheetContext).textTheme.titleMedium),
                const SizedBox(height: 8),
                const Text('Проверьте подключение к интернету и попробуйте ещё раз.', textAlign: TextAlign.center),
                const SizedBox(height: 20),
                SizedBox(
                  width: double.infinity,
                  child: FilledButton.icon(
                    onPressed: () {
                      Navigator.pop(sheetContext);
                      unawaited(_api.initialize());
                    },
                    icon: const Icon(Icons.refresh),
                    label: const Text('Повторить'),
                  ),
                ),
              ],
            ),
          ),
    );
    _modalSheetOpen = false;
  }

  Future<void> _showUpdateRequired() async {
    await showModalBottomSheet<void>(
      context: context,
      isDismissible: false,
      enableDrag: false,
      builder:
          (sheetContext) => Padding(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.system_update, size: 40),
                const SizedBox(height: 12),
                Text('Требуется обновление', style: Theme.of(sheetContext).textTheme.titleMedium),
                const SizedBox(height: 8),
                const Text('Обновите приложение для продолжения работы.', textAlign: TextAlign.center),
                const SizedBox(height: 20),
                if (_api.updateUrl case final url? when url.isNotEmpty)
                  FilledButton(
                    onPressed: () => launchUrl(Uri.parse(url), mode: LaunchMode.externalApplication),
                    child: const Text('Обновить'),
                  ),
              ],
            ),
          ),
    );
  }

  Future<void> _checkActiveSession() async {
    _checkingRequest = true;
    try {
      final session = await _api.getCurrentSession();
      if (mounted) setState(() => _activeSession = session);
    } catch (_) {
      // no active session
    }
    _checkingRequest = false;
    if (mounted) setState(() => _checking = false);
  }

  void _continueCase() {
    if (_activeSession == null) return;
    final s = _activeSession!;
    context.read<SessionService>().resume(
      SessionState(sessionId: s.id, caseName: s.caseName, actionPoints: s.actionPoints, phase: s.phase),
    );
    Navigator.push(context, MaterialPageRoute(builder: (_) => const DeskScreen()));
  }

  @override
  Widget build(BuildContext context) {
    final ready = _api.initializationStatus == InitializationStatus.ready;
    final loading = !ready || _checking;
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
              Text(
                'ДетектИИв',
                style: theme.textTheme.headlineLarge?.copyWith(fontWeight: FontWeight.bold, letterSpacing: 2),
              ),
              const SizedBox(height: 4),
              const _CaseNumberAnimation(),
              const SizedBox(height: 48),
              if (_activeSession != null || loading)
                _TitleButton(
                  label: loading ? 'Проверяем...' : 'Продолжить дело',
                  icon: loading ? Icons.hourglass_top : Icons.play_arrow,
                  onPressed: loading ? null : _continueCase,
                ),
              if (_activeSession != null || loading) const SizedBox(height: 12),
              _TitleButton(
                label: 'Новое дело',
                icon: Icons.folder_open,
                onPressed:
                    loading
                        ? null
                        : () {
                          Navigator.push(context, MaterialPageRoute(builder: (_) => const LoadingScreen()));
                        },
              ),
              const SizedBox(height: 12),
              _TitleButton(
                label: 'Предыдущие дела',
                icon: Icons.history,
                onPressed:
                    loading
                        ? null
                        : () {
                          Navigator.push(context, MaterialPageRoute(builder: (_) => const PreviousCasesScreen()));
                        },
              ),
              const SizedBox(height: 12),
              _TitleButton(
                label: 'Настройки',
                icon: Icons.settings,
                onPressed:
                    loading
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
      style: theme.textTheme.titleMedium?.copyWith(color: theme.colorScheme.onSurface.withAlpha(120), letterSpacing: 4),
    );
  }
}
