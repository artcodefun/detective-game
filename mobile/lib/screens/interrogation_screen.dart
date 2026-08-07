import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:speech_to_text/speech_recognition_result.dart';
import 'package:speech_to_text/speech_to_text.dart' as stt;

import '../blocs/session_cubit.dart';
import '../models/game_state.dart';
import '../services/api_service.dart';
import '../widgets/chat_bubble.dart';
import '../widgets/mood_indicator.dart';

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
  bool _isWaiting = false;
  bool _isListening = false;
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
    _startInterrogation();
  }

  @override
  void dispose() {
    _speech.stop();
    _textController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _startInterrogation() async {
    try {
      final interId = widget.interrogationId ?? (await _api.createInterrogation(widget.characterId)).id;
      _interId = interId;

      if (widget.interrogationId != null) {
        final messages = await _api.getInterrogationMessages(interId);
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
      _isListening = false;
      await _speech.stop();
      setState(() => _isListening = false);
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
    _isListening = true;

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
    if (mounted) setState(() => _isListening = true);
  }

  void _onSpeechError(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(message)));
  }

  void _onSpeechResult(SpeechRecognitionResult result) {
    if (!_isListening) return;
    final words = result.recognizedWords;
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
    _isListening = false;
    setState(() => _isListening = false);
  }

  Future<void> _sendMessage() async {
    if (_isListening) {
      _isListening = false;
      await _speech.stop();
      if (mounted) setState(() => _isListening = false);
    }

    final text = _textController.text.trim();
    if (text.isEmpty || _isWaiting || _interId == null) return;

    _textController.clear();

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
      _isWaiting = true;
    });
    _scrollToBottom();

    try {
      final msg = await _api.addInterrogationMessage(interId: _interId!, message: text);
      if (mounted) {
        setState(() {
          _messages.add(msg);
          if (_character != null && msg.attitudeDelta != 0) {
            _character = _character!.copyWith(
              trust: (_character!.trust + msg.attitudeDelta).clamp(Character.minTrust, Character.maxTrust),
            );
          }
          _isWaiting = false;
        });
        _scrollToBottom();
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isWaiting = false);
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Ошибка: $e')));
      }
    }
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  Future<void> _closeInterrogation() async {
    if (_interId != null) {
      final cubit = context.read<SessionCubit>();
      try {
        await _api.completeInterrogation(_interId!);
        cubit.refreshSession();
      } catch (_) {}
    }
    if (mounted) Navigator.pop(context);
  }

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
            Text(char.name),
          ],
        ),
        actions: [
          IconButton(icon: const Icon(Icons.close), onPressed: _closeInterrogation, tooltip: 'Завершить допрос'),
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
    if (_messages.isEmpty) {
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

    return ListView.builder(
      controller: _scrollController,
      padding: const EdgeInsets.all(16),
      itemCount: _messages.length + (_isWaiting ? 1 : 0),
      itemBuilder: (_, index) {
        if (index == _messages.length && _isWaiting) {
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

        final msg = _messages[index];
        return ChatBubble(text: msg.text, isPlayer: msg.fromUser, senderName: msg.fromUser ? 'Вы' : characterName);
      },
    );
  }

  Widget _buildInputBar(ThemeData theme, ColorScheme colorScheme) {
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(12, 8, 12, 8),
        child: Row(
          children: [
            IconButton(
              onPressed: _isWaiting ? null : _toggleListening,
              icon: _isListening ? const Icon(Icons.mic, color: Colors.red) : const Icon(Icons.mic_none),
              tooltip: _isListening ? 'Остановить запись' : 'Голосовой ввод',
            ),
            const SizedBox(width: 4),
            Expanded(
              child: TextField(
                controller: _textController,
                enabled: !_isWaiting,
                maxLines: 3,
                minLines: 1,
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
              onPressed: _isWaiting ? null : _sendMessage,
              icon: const Icon(Icons.send),
              tooltip: 'Отправить',
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
