import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../blocs/session_cubit.dart';
import '../services/api_service.dart';
import 'desk_screen.dart';

class LoadingScreen extends StatefulWidget {
  const LoadingScreen({super.key});

  @override
  State<LoadingScreen> createState() => _LoadingScreenState();
}

class _LoadingScreenState extends State<LoadingScreen> {
  String _status = 'Готовим сценарий...';
  bool _error = false;

  @override
  void initState() {
    super.initState();
    _generate();
  }

  Future<void> _generate() async {
    try {
      setState(() => _status = 'Собираем улики, опрашиваем свидетелей...');
      await context.read<SessionCubit>().startNewGame();
      if (!mounted) return;
      Navigator.pushReplacement(context, MaterialPageRoute(builder: (_) => const DeskScreen()));
    } on ApiException catch (e) {
      if (!mounted) return;
      setState(() {
        _error = true;
        _status = 'Ошибка: ${e.message}';
      });
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _error = true;
        _status = 'Не удалось подключиться к серверу';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Новое дело'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => Navigator.pop(context)),
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 32),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                _error ? Icons.error_outline : Icons.auto_awesome,
                size: 64,
                color: _error ? colorScheme.error : colorScheme.primary,
              ),
              const SizedBox(height: 24),
              Text(_status, style: theme.textTheme.titleMedium, textAlign: TextAlign.center),
              const SizedBox(height: 24),
              if (!_error) ...[
                const LinearProgressIndicator(),
                const SizedBox(height: 16),
                Text(
                  'Пожалуйста, подождите...',
                  style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(120)),
                ),
              ],
              if (_error) ...[
                const SizedBox(height: 8),
                FilledButton.icon(
                  onPressed: () {
                    setState(() {
                      _error = false;
                      _status = 'Повторная попытка...';
                    });
                    _generate();
                  },
                  icon: const Icon(Icons.refresh),
                  label: const Text('Повторить'),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
