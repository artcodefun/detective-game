import 'dart:async';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../services/api_service.dart';
import '../services/session_service.dart';
import 'desk_screen.dart';

const _messageTexts = [
  'Открываем дело...',
  'Осматриваем место преступления...',
  'Собираем вещественные доказательства...',
  'Опрашиваем свидетелей...',
  'Изучаем отпечатки пальцев...',
  'Проверяем алиби подозреваемых...',
  'Анализируем улики в лаборатории...',
  'Составляем фоторобот...',
  'Сверяем показания...',
  'Ищем несостыковки в хронологии...',
  'Заливаем кофе в термос...',
  'Протираем лупу...',
  'Надеваем перчатки...',
  'Достаём диктофон...',
  'Печатаем протокол...',
  'Проветриваем кабинет от сигарного дыма...',
  'Кормим служебного пса...',
  'Затачиваем карандаши...',
];

class LoadingScreen extends StatefulWidget {
  const LoadingScreen({super.key});

  @override
  State<LoadingScreen> createState() => _LoadingScreenState();
}

class _LoadingScreenState extends State<LoadingScreen> with SingleTickerProviderStateMixin {
  int _msgIndex = 0;
  bool _error = false;
  String _errorText = '';
  Timer? _timer;
  late AnimationController _fadeController;
  late Animation<double> _fadeAnim;
  late List<String> _messages;

  @override
  void initState() {
    super.initState();
    _fadeController = AnimationController(duration: const Duration(milliseconds: 400), vsync: this);
    _fadeAnim = CurvedAnimation(parent: _fadeController, curve: Curves.easeInOut);

    _messages = List.from(_messageTexts)..shuffle();
    _timer = Timer.periodic(const Duration(seconds: 5), (_) => _nextMessage());
    _fadeController.forward();

    _generate();
  }

  @override
  void dispose() {
    _timer?.cancel();
    _fadeController.dispose();
    super.dispose();
  }

  void _nextMessage() {
    _fadeController.reverse().then((_) {
      if (!mounted) return;
      setState(() => _msgIndex = (_msgIndex + 1) % _messages.length);
      _fadeController.forward();
    });
  }

  Future<void> _generate() async {
    try {
      await context.read<SessionService>().startNewGame();
      if (!mounted) return;
      _timer?.cancel();
      Navigator.pushReplacement(context, MaterialPageRoute(builder: (_) => const DeskScreen()));
    } on ApiException catch (e) {
      if (!mounted) return;
      _timer?.cancel();
      setState(() {
        _error = true;
        _errorText = e.message;
      });
    } catch (e) {
      if (!mounted) return;
      _timer?.cancel();
      setState(() {
        _error = true;
        _errorText = 'Не удалось подключиться к серверу';
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
                _error ? Icons.error_outline : Icons.search,
                size: 64,
                color: _error ? colorScheme.error : colorScheme.primary,
              ),
              const SizedBox(height: 32),
              if (!_error) ...[
                SizedBox(
                  height: 48,
                  child: Center(
                    child: FadeTransition(
                      opacity: _fadeAnim,
                      child: Text(
                        _messages[_msgIndex],
                        style: theme.textTheme.titleMedium,
                        textAlign: TextAlign.center,
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 24),
                const LinearProgressIndicator(),
                const SizedBox(height: 16),
                Text(
                  'Пожалуйста, подождите...',
                  style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(120)),
                ),
              ],
              if (_error) ...[
                Text('Ошибка: $_errorText', style: theme.textTheme.titleMedium, textAlign: TextAlign.center),
                const SizedBox(height: 24),
                FilledButton.icon(
                  onPressed: () {
                    setState(() {
                      _error = false;
                      _msgIndex = 0;
                      _messages = List.from(_messageTexts)..shuffle();
                    });
                    _timer = Timer.periodic(const Duration(seconds: 3), (_) => _nextMessage());
                    _fadeController.forward();
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
