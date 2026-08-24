import 'dart:async';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:speech_to_text/speech_recognition_result.dart';
import 'package:speech_to_text/speech_to_text.dart' as stt;

import '../models/game_state.dart';
import '../services/api_service.dart';
import '../services/audio_service.dart';
import '../services/session_service.dart';
import '../widgets/chat_bubble.dart';
import '../widgets/mood_indicator.dart';

const _maxQuestionsPerInterrogation = 10;
const _maxQuestionLength = 600;

class InterrogationScreen extends StatefulWidget {
  final String characterId;
  final String? interrogationId;

  const InterrogationScreen({super.key, required this.characterId, this.interrogationId});

  @override
  State<InterrogationScreen> createState() => _InterrogationScreenState();
}

class _InterrogationScreenState extends State<InterrogationScreen> {
  final _messages = <ChatMessage>[];
  final _textController = TextEditingController();
  final _scrollController = ScrollController();
  final _speech = stt.SpeechToText();
  late final AudioService _audio;
  bool _isWaiting = false;
  bool _isCompleted = false;
  String? _pendingMessageText;
  String? _failedMessageText;
  bool _isListening = false;
  bool _isStoppingSpeech = false;
  Completer<void>? _speechFinalResult;
  bool _speechInitialized = false;
  String _speechBaseText = '';
  String _speechFinalized = '';
  String _speechPartial = '';
  String? _interId;
  Character? _character;

  ApiService get _api => context.read<ApiService>();

  @override
  void initState() {
    super.initState();
    _audio = context.read<AudioService>();
    _audio.pauseMusic();
    _startInterrogation();
  }

  @override
  void dispose() {
    unawaited(_resumeMusic());
    _textController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _resumeMusic() async {
    if (_isListening) {
      await _stopListening();
    }
    await _audio.resumeMusic();
  }

  Future<void> _startInterrogation() async {
    try {
      final interrogation =
          widget.interrogationId == null
              ? await _api.createInterrogation(widget.characterId)
              : await _api.getInterrogation(widget.interrogationId!);
      _interId = interrogation.id;
      _isCompleted = !interrogation.isActive;

      if (widget.interrogationId != null) {
        final messages = await _api.getInterrogationMessages(interrogation.id);
        if (mounted) setState(() => _messages.addAll(messages));
      }

      final character = await _api.getCharacter(widget.characterId);
      if (mounted) setState(() => _character = character);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Ошибка: $e')));
        Navigator.pop(context);
      }
    }
  }

  Future<void> _toggleListening() async {
    if (_isListening) {
      await _stopListening();
      return;
    }

    if (!_speechInitialized) {
      final available = await _speech.initialize(
        onError: (_) => _onSpeechError('Ошибка инициализации микрофона'),
        onStatus: _onSpeechStatus,
      );
      if (!available) {
        _onSpeechError('Голосовой ввод недоступен на этом устройстве');
        return;
      }
      if (mounted) setState(() => _speechInitialized = true);
    }

    final hasPermission = await _speech.hasPermission;
    if (!hasPermission) {
      _onSpeechError('Разрешение на использование микрофона не получено');
      return;
    }

    _speechBaseText = _textController.text;
    _speechFinalized = '';
    _speechPartial = '';
    await _audio.beforeSpeechRecognition();
    _isListening = true;
    if (mounted) setState(() {});

    try {
      await _speech.listen(
        onResult: _onSpeechResult,
        listenOptions: stt.SpeechListenOptions(
          localeId: 'ru_RU',
          listenMode: stt.ListenMode.dictation,
          pauseFor: const Duration(seconds: 10),
          listenFor: const Duration(seconds: 120),
          autoPunctuation: true,
        ),
      );
    } catch (_) {
      await _finishListening();
      _onSpeechError('Не удалось начать голосовой ввод');
    }
  }

  void _onSpeechError(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(message)));
  }

  void _onSpeechResult(SpeechRecognitionResult result) {
    if (!_isListening) return;
    final words = result.recognizedWords;
    if (result.finalResult) {
      final finalResult = _speechFinalResult;
      if (finalResult != null && !finalResult.isCompleted) {
        finalResult.complete();
      }
    }
    if (words.isEmpty) {
      if (_speechPartial.isNotEmpty) {
        final sep = _speechFinalized.isEmpty ? '' : ' ';
        final punctuation =
            _speechPartial.endsWith('.') || _speechPartial.endsWith('?') || _speechPartial.endsWith('!') ? '' : '.';
        _speechFinalized += sep + _speechPartial + punctuation;
        _speechPartial = '';
      }
      return;
    }
    if (result.finalResult) {
      final sep = _speechFinalized.isEmpty ? '' : ' ';
      _speechFinalized += sep + words;
      _speechPartial = '';
    } else {
      _speechPartial = words;
    }
    final parts = <String>[];
    if (_speechBaseText.isNotEmpty) parts.add(_speechBaseText);
    if (_speechFinalized.isNotEmpty) parts.add(_speechFinalized);
    if (_speechPartial.isNotEmpty) parts.add(_speechPartial);
    final fullText = parts.join(' ');
    _textController.text = fullText;
    _textController.selection = TextSelection.fromPosition(TextPosition(offset: fullText.length));
  }

  void _onSpeechStatus(String status) {
    if (!mounted) return;
    if (status == stt.SpeechToText.listeningStatus) return;
    if (_isStoppingSpeech) return;
    unawaited(_finishListening());
  }

  Future<void> _stopListening() async {
    if (!_isListening) return;
    _isStoppingSpeech = true;
    final finalResult = Completer<void>();
    _speechFinalResult = finalResult;
    try {
      await _speech.stop();
      await Future.any([finalResult.future, Future<void>.delayed(const Duration(seconds: 2))]);
      await Future<void>.delayed(const Duration(milliseconds: 100));
    } catch (_) {
      // Sending a message should not depend on stopping speech recognition.
    } finally {
      try {
        await _finishListening();
      } finally {
        _speechFinalResult = null;
        _isStoppingSpeech = false;
      }
    }
  }

  Future<void> _finishListening() async {
    final wasListening = _isListening;
    _isListening = false;
    if (mounted) setState(() {});
    if (wasListening) {
      await _audio.afterSpeechRecognition();
    }
  }

  Future<void> _sendMessage() async {
    if (_isListening) {
      await _stopListening();
    }

    final text = _textController.text.trim();
    if (text.isEmpty || _isWaiting || _interId == null || _isCompleted) return;

    if (_questionsAsked >= _maxQuestionsPerInterrogation) {
      setState(() => _isCompleted = true);
      return;
    }

    setState(() {
      _isWaiting = true;
      _pendingMessageText = text;
      _failedMessageText = null;
      _textController.clear();
    });
    _scrollToBottom();

    try {
      final msg = await _api.addInterrogationMessage(interId: _interId!, message: text);
      if (mounted) {
        setState(() {
          _messages.add(
            ChatMessage(
              id: '',
              sessionId: '',
              interrogationId: _interId!,
              fromUser: true,
              text: text,
              timestamp: DateTime.now(),
            ),
          );
          _messages.add(msg);
          _pendingMessageText = null;
          if (_character != null && msg.attitudeDelta != 0) {
            _character = _character!.copyWith(
              trust: (_character!.trust + msg.attitudeDelta).clamp(Character.minTrust, Character.maxTrust),
            );
          }
          _isCompleted = _questionsAsked >= _maxQuestionsPerInterrogation;
          _isWaiting = false;
        });
        if (_character != null) {
          unawaited(_audio.playSuspectReply(_character!));
        }
        _scrollToBottom();
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isWaiting = false;
          _pendingMessageText = null;
          _failedMessageText = text;
          _textController.value = TextEditingValue(text: text, selection: TextSelection.collapsed(offset: text.length));
        });
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Ошибка: $e')));
      }
    }
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(0, duration: const Duration(milliseconds: 300), curve: Curves.easeOut);
      }
    });
  }

  Future<void> _closeInterrogation() async {
    if (_isCompleted) {
      if (mounted) Navigator.pop(context);
      return;
    }
    if (_interId != null) {
      final sessionService = context.read<SessionService>();
      try {
        await _api.completeInterrogation(_interId!);
        unawaited(sessionService.refresh());
      } catch (_) {}
    }
    if (mounted) Navigator.pop(context);
  }

  Future<void> _confirmCloseInterrogation() async {
    if (_isWaiting) return;
    if (_isCompleted) {
      Navigator.pop(context);
      return;
    }

    final shouldComplete = await showModalBottomSheet<bool>(
      context: context,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(20))),
      builder:
          (sheetContext) => SafeArea(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(20, 20, 20, 24),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('Завершить допрос?', style: Theme.of(sheetContext).textTheme.titleMedium),
                  const SizedBox(height: 8),
                  const Text('Вы больше не сможете задавать вопросы этому персонажу в рамках этого допроса.'),
                  const SizedBox(height: 20),
                  SizedBox(
                    width: double.infinity,
                    child: FilledButton(
                      onPressed: () => Navigator.pop(sheetContext, true),
                      child: const Text('Завершить допрос'),
                    ),
                  ),
                  const SizedBox(height: 8),
                  SizedBox(
                    width: double.infinity,
                    child: TextButton(
                      onPressed: () => Navigator.pop(sheetContext, false),
                      child: const Text('Продолжить разговор'),
                    ),
                  ),
                ],
              ),
            ),
          ),
    );
    if (shouldComplete == true && mounted) {
      await _closeInterrogation();
    }
  }

  int get _questionsAsked =>
      _messages.where((message) => message.fromUser).length + (_pendingMessageText == null ? 0 : 1);

  int get _questionsRemaining =>
      (_maxQuestionsPerInterrogation - _questionsAsked).clamp(0, _maxQuestionsPerInterrogation);

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    if (_character == null) {
      return Scaffold(
        appBar: AppBar(title: const Text('Допрос')),
        body: const Center(child: CircularProgressIndicator()),
      );
    }

    final char = _character!;
    final mood = char.trustLevel;

    return Scaffold(
      resizeToAvoidBottomInset: true,
      appBar: AppBar(
        scrolledUnderElevation: 0,
        backgroundColor: colorScheme.surface,
        title: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            CircleAvatar(
              radius: 16,
              backgroundColor: colorScheme.primaryContainer,
              child: Text(char.name[0], style: TextStyle(fontSize: 14, color: colorScheme.onPrimaryContainer)),
            ),
            const SizedBox(width: 8),
            Flexible(child: Text(char.name)),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.close),
            onPressed: _isWaiting ? null : _confirmCloseInterrogation,
            tooltip: _isCompleted ? 'Закрыть' : 'Завершить допрос',
          ),
        ],
      ),
      body: Column(
        children: [
          _buildCharacterHeader(theme, colorScheme, char, mood),
          const Divider(height: 1),
          Expanded(child: _buildChatList(theme, char.name)),
          const Divider(height: 1),
          _buildInputBar(theme, colorScheme),
        ],
      ),
    );
  }

  Widget _buildCharacterHeader(ThemeData theme, ColorScheme colorScheme, Character char, TrustLevel mood) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(char.name, style: theme.textTheme.titleSmall),
                const SizedBox(height: 2),
                Text(
                  char.profession,
                  style: theme.textTheme.bodySmall?.copyWith(color: colorScheme.onSurface.withAlpha(140)),
                ),
              ],
            ),
          ),
          Icon(_getMoodIcon(mood), color: _getMoodColor(mood), size: 20),
          const SizedBox(width: 6),
          MoodIndicator(trustLevel: mood),
          const SizedBox(width: 6),
          Text(
            _getMoodLabel(mood),
            style: theme.textTheme.bodySmall?.copyWith(color: _getMoodColor(mood), fontWeight: FontWeight.w500),
          ),
        ],
      ),
    );
  }

  Widget _buildChatList(ThemeData theme, String characterName) {
    if (_messages.isEmpty && _pendingMessageText == null && _failedMessageText == null) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Text(
            'Задайте первый вопрос, чтобы начать допрос',
            textAlign: TextAlign.center,
            style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurface.withAlpha(100)),
          ),
        ),
      );
    }

    return Align(
      alignment: Alignment.topCenter,
      child: ListView.builder(
        controller: _scrollController,
        reverse: true,
        shrinkWrap: true,
        padding: const EdgeInsets.all(16),
        itemCount:
            _messages.length +
            (_pendingMessageText == null ? 0 : 1) +
            (_failedMessageText == null ? 0 : 1) +
            (_isWaiting ? 1 : 0),
        itemBuilder: (_, index) {
          var itemIndex = index;

          if (_isWaiting) {
            if (itemIndex == 0) {
              return const Padding(
                padding: EdgeInsets.symmetric(vertical: 8),
                child: Row(
                  children: [
                    SizedBox(width: 36),
                    SizedBox(width: 24, height: 24, child: CircularProgressIndicator(strokeWidth: 2)),
                  ],
                ),
              );
            }
            itemIndex--;
          }

          if (_pendingMessageText case final pendingText?) {
            if (itemIndex == 0) {
              return ChatBubble(text: pendingText, isPlayer: true, senderName: 'Вы');
            }
            itemIndex--;
          }

          if (_failedMessageText case final failedText?) {
            if (itemIndex == 0) {
              return ChatBubble(text: failedText, isPlayer: true, senderName: 'Вы', isFailed: true);
            }
            itemIndex--;
          }

          final message = _messages[_messages.length - 1 - itemIndex];
          return ChatBubble(
            text: message.text,
            isPlayer: message.fromUser,
            senderName: message.fromUser ? 'Вы' : characterName,
          );
        },
      ),
    );
  }

  Widget _buildInputBar(ThemeData theme, ColorScheme colorScheme) {
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(12, 8, 12, 8),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              _isCompleted
                  ? 'Допрос завершён'
                  : 'Вопросов осталось: $_questionsRemaining из $_maxQuestionsPerInterrogation',
              style: theme.textTheme.labelSmall?.copyWith(color: colorScheme.onSurface.withAlpha(140)),
            ),
            const SizedBox(height: 4),
            Row(
              children: [
                IconButton(
                  onPressed: _isWaiting || _isCompleted ? null : _toggleListening,
                  icon: _isListening ? const Icon(Icons.mic, color: Colors.red) : const Icon(Icons.mic_none),
                  tooltip: _isListening ? 'Остановить запись' : 'Голосовой ввод',
                ),
                const SizedBox(width: 4),
                Expanded(
                  child: TextField(
                    controller: _textController,
                    enabled: !_isWaiting && !_isCompleted,
                    maxLines: 3,
                    minLines: 1,
                    maxLength: _maxQuestionLength,
                    textCapitalization: TextCapitalization.sentences,
                    decoration: InputDecoration(
                      hintText: _isListening ? 'Говорите...' : 'Задайте вопрос...',
                      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                      border: OutlineInputBorder(borderRadius: BorderRadius.circular(24)),
                      filled: true,
                      fillColor: _isListening ? Colors.red.withAlpha(15) : colorScheme.surfaceContainerHighest,
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                IconButton.filled(
                  onPressed: _isWaiting || _isCompleted ? null : _sendMessage,
                  icon: const Icon(Icons.send),
                  tooltip: 'Отправить',
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  IconData _getMoodIcon(TrustLevel level) {
    switch (level) {
      case TrustLevel.open:
        return Icons.sentiment_satisfied_alt;
      case TrustLevel.reserved:
        return Icons.sentiment_neutral;
      case TrustLevel.tense:
        return Icons.sentiment_dissatisfied;
      case TrustLevel.closed:
        return Icons.mood_bad;
    }
  }

  Color _getMoodColor(TrustLevel level) {
    switch (level) {
      case TrustLevel.open:
        return Colors.green;
      case TrustLevel.reserved:
        return Colors.yellow.shade700;
      case TrustLevel.tense:
        return Colors.orange;
      case TrustLevel.closed:
        return Colors.red;
    }
  }

  String _getMoodLabel(TrustLevel level) {
    switch (level) {
      case TrustLevel.open:
        return 'Открыт';
      case TrustLevel.reserved:
        return 'Сдержан';
      case TrustLevel.tense:
        return 'Напряжён';
      case TrustLevel.closed:
        return 'Закрыт';
    }
  }
}
