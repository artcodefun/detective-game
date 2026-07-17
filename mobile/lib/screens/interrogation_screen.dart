import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:speech_to_text/speech_recognition_result.dart';
import 'package:speech_to_text/speech_to_text.dart' as stt;

import '../blocs/session_cubit.dart';
import '../models/character.dart';
import '../models/chronology_entry.dart';
import '../models/game_state.dart';
import '../models/notebook.dart';
import '../services/mock_llm_service.dart';
import '../widgets/chat_bubble.dart';
import '../widgets/mood_indicator.dart';

class InterrogationScreen extends StatefulWidget {
  final CharacterState characterState;

  const InterrogationScreen({super.key, required this.characterState});

  @override
  State<InterrogationScreen> createState() => _InterrogationScreenState();
}

class _InterrogationScreenState extends State<InterrogationScreen> {
  final _llm = MockLlmService();
  final _messages = <InterrogationMessage>[];
  final _textController = TextEditingController();
  final _scrollController = ScrollController();
  final _speech = stt.SpeechToText();
  bool _isWaiting = false;
  bool _isListening = false;
  bool _speechInitialized = false;
  String _speechBaseText = '';
  String _speechFinalized = '';
  String _speechPartial = '';
  String? _activeChronologyId;
  late CharacterState _character;
  late GameSession _session;

  @override
  void initState() {
    super.initState();
    _character = widget.characterState;
    _session = context.read<SessionCubit>().state!;
    _initSpeech();
  }

  @override
  void dispose() {
    _speech.stop();
    _textController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _initSpeech() async {
    await _speech.initialize(onStatus: _onSpeechStatus);
    if (mounted) {
      setState(() => _speechInitialized = true);
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
    if (text.isEmpty || _isWaiting) return;

    _textController.clear();

    final playerMsg = InterrogationMessage(sender: 'Вы', text: text, timestamp: DateTime.now());

    setState(() {
      _messages.add(playerMsg);
      _isWaiting = true;
    });
    _scrollToBottom();

    final response = await _llm.respondInInterrogation(characterState: _character, playerMessage: text);

    final characterMsg = InterrogationMessage(
      sender: _character.base.name,
      text: response.answer,
      timestamp: DateTime.now(),
    );

    setState(() {
      _messages.add(characterMsg);
      _character = _character.applyAttitudeDelta(response.attitudeDelta).addMessage(playerMsg).addMessage(characterMsg);
      _isWaiting = false;
    });

    _saveStatements(response.statements);
    if (!mounted) return;
    context.read<SessionCubit>().update(_session);
    _scrollToBottom();
  }

  void _saveStatements(List<String> statements) {
    if (statements.isEmpty) return;
    final entries = statements.map((s) => NotebookEntry(
      id: 'note_${DateTime.now().millisecondsSinceEpoch}_${s.hashCode}',
      type: NotebookEntryType.statement,
      characterId: _character.base.id,
      description: s,
      timestamp: DateTime.now(),
    )).toList();

    if (_activeChronologyId == null) {
      final chron = ChronologyEntry(
        id: 'chron_${DateTime.now().millisecondsSinceEpoch}',
        eventType: ChronologyEventType.interrogation,
        title: 'Допрос (${_character.base.name})',
        timestamp: DateTime.now(),
        details: entries,
      );
      _session = _session.addChronologyEntry(chron);
      _activeChronologyId = chron.id;
    } else {
      _session = _session.addDetailsToChronology(_activeChronologyId!, entries);
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

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final char = _character.base;
    final mood = _character.trustLevel;

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
          IconButton(
            icon: const Icon(Icons.close),
            onPressed: () => Navigator.pop(context),
            tooltip: 'Завершить допрос',
          ),
        ],
      ),
      body: Column(
        children: [
          _buildCharacterHeader(theme, colorScheme, char, mood),
          const Divider(height: 1),
          Expanded(child: _buildChatList(theme)),
          const Divider(height: 1),
          _buildInputBar(theme, colorScheme),
        ],
      ),
    );
  }

  Widget _buildCharacterHeader(ThemeData theme, ColorScheme colorScheme, CharacterData char, TrustLevel mood) {
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

  Widget _buildChatList(ThemeData theme) {
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
        final isPlayer = msg.sender == 'Вы';

        return ChatBubble(text: msg.text, isPlayer: isPlayer, senderName: msg.sender);
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
