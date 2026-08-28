import 'dart:async';
import 'dart:developer' as developer;
import 'dart:math';

import 'package:audio_session/audio_session.dart';
import 'package:flutter_soloud/flutter_soloud.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../models/game_state.dart';

class AudioService {
  AudioService();

  static const _suspectReplyCount = 11;
  static const _maxMusicVolume = 0.3;
  static const _soundVolumeKey = 'sound_volume';
  static const _musicVolumeKey = 'music_volume';
  static final _suspectReplyAssets = List.generate(
    _suspectReplyCount,
    (index) => 'assets/audio/gibberish/gibberish-$index.mp3',
  );
  static const _musicAssets = [
    'assets/audio/music/alex-morgan-dark-suspense-thriller-528314.mp3',
    'assets/audio/music/jorisvermeer-crime-piano-tension-532082.mp3',
    'assets/audio/music/konstantinpazuzustudio-true-crime-podcast-background-police-procedural-518902.mp3',
    'assets/audio/music/the_mountain-detective-167604.mp3',
    'assets/audio/music/universfield-last-letter-531688.mp3',
  ];

  final _random = Random();
  SharedPreferences? _preferences;
  AudioSource? _suspectReplySource;
  Future<void>? _engineRestartTask;
  Future<void>? _musicStartTask;
  Timer? _nextTrackTimer;
  AudioSource? _musicSource;
  SoundHandle? _musicHandle;
  Duration? _musicDuration;
  Duration? _musicResumePosition;
  int? _currentMusicIndex;
  bool _pausedForAppLifecycle = false;
  bool _musicPaused = false;
  bool _resumeAfterAppLifecycle = false;
  double _soundVolume = 0.7;
  double _musicVolume = 0.5;

  double get soundVolume => _soundVolume;

  double get musicVolume => _musicVolume;

  Future<void> initialize() async {
    final soloud = SoLoud.instance;
    if (soloud.isInitialized) {
      soloud.deinit();
    }

    try {
      final preferences = await SharedPreferences.getInstance();
      _preferences = preferences;
      _soundVolume = _storedVolume(preferences, _soundVolumeKey, _soundVolume);
      _musicVolume = _storedVolume(preferences, _musicVolumeKey, _musicVolume);
    } catch (error, stackTrace) {
      developer.log('Failed to load audio settings', name: 'AudioService', error: error, stackTrace: stackTrace);
    }
  }

  void setSoundVolume(double volume) {
    _soundVolume = volume.clamp(0.0, 1.0).toDouble();
    unawaited(_saveVolume(_soundVolumeKey, _soundVolume));
  }

  void setMusicVolume(double volume) {
    _musicVolume = volume.clamp(0.0, 1.0).toDouble();
    if (_musicHandle case final handle?) {
      SoLoud.instance.setVolume(handle, _effectiveMusicVolume);
    }
    unawaited(_saveVolume(_musicVolumeKey, _musicVolume));
  }

  Future<void> playSuspectReply(Character character) async {
    try {
      await _engineRestartTask;
      final soloud = SoLoud.instance;
      if (!await _activatePlayback()) return;
      await _ensureEngineStarted();
      if (_suspectReplySource case final previousSource?) {
        await soloud.disposeSource(previousSource);
      }
      final source = await soloud.loadAsset(_suspectReplyAssets[_random.nextInt(_suspectReplyAssets.length)]);
      _suspectReplySource = source;
      // ignore: experimental_member_use
      source.filters.pitchShiftFilter.activate();

      final handle = soloud.play(source, volume: _soundVolume);
      // ignore: experimental_member_use
      source.filters.pitchShiftFilter.semitones(soundHandle: handle).value = _pitch(character);
      // ignore: experimental_member_use
      source.filters.pitchShiftFilter.timeStretch(handle, _speed(character));
    } catch (error, stackTrace) {
      developer.log('Failed to play suspect reply', name: 'AudioService', error: error, stackTrace: stackTrace);
    }
  }

  Future<void> resumeMusic() async {
    _musicPaused = false;
    await _engineRestartTask;
    if (_musicHandle case final handle?) {
      if (!await _activatePlayback()) return;
      SoLoud.instance.setPause(handle, false);
      _musicResumePosition = null;
      _scheduleNextTrack(handle, _musicSource!, _remainingMusicDuration(handle));
      return;
    }
    await (_musicStartTask ??= _loadAndStartMusic().whenComplete(() => _musicStartTask = null));
  }

  void pauseMusic() {
    _musicPaused = true;
    _nextTrackTimer?.cancel();
    if (_musicHandle case final handle?) {
      _musicResumePosition = SoLoud.instance.getPosition(handle);
      SoLoud.instance.setPause(handle, true);
    }
  }

  void pauseForAppLifecycle() {
    if (_pausedForAppLifecycle) return;

    _pausedForAppLifecycle = true;
    _resumeAfterAppLifecycle = !_musicPaused;
    pauseMusic();
  }

  Future<void> resumeAfterAppLifecycle() async {
    if (!_pausedForAppLifecycle) return;

    _pausedForAppLifecycle = false;
    if (!_resumeAfterAppLifecycle) return;

    _resumeAfterAppLifecycle = false;
    await resumeMusic();
  }

  Future<void> afterSpeechRecognition() {
    return _engineRestartTask ??= _restartEngine().whenComplete(() => _engineRestartTask = null);
  }

  Future<void> beforeSpeechRecognition() async {
    await _engineRestartTask;
    await _musicStartTask;

    final session = await AudioSession.instance;
    await session.setActive(false);
  }

  Future<void> _restartEngine() async {
    final soloud = SoLoud.instance;
    if (soloud.isInitialized) {
      soloud.deinit();
    }
    _suspectReplySource = null;
    _musicSource = null;
    _musicHandle = null;
    _musicDuration = null;
    _nextTrackTimer?.cancel();

    if (!await _activatePlayback()) return;
    await soloud.init();
  }

  Future<void> _loadAndStartMusic() async {
    try {
      await _engineRestartTask;
      final soloud = SoLoud.instance;
      if (!await _activatePlayback()) return;
      await _ensureEngineStarted();

      final resumePosition = _musicResumePosition;
      final musicIndex = resumePosition != null && _currentMusicIndex != null ? _currentMusicIndex! : _nextMusicIndex();
      _musicSource = await soloud.loadAsset(_musicAssets[musicIndex]);
      _musicDuration = soloud.getLength(_musicSource!);
      _musicHandle = soloud.play(
        _musicSource!,
        volume: _effectiveMusicVolume,
        paused: _musicPaused || resumePosition != null,
      );
      if (resumePosition != null) {
        soloud.seek(_musicHandle!, resumePosition);
        _musicResumePosition = null;
        if (!_musicPaused) {
          soloud.setPause(_musicHandle!, false);
        }
      }
      if (!_musicPaused) {
        _scheduleNextTrack(_musicHandle!, _musicSource!, _remainingMusicDuration(_musicHandle!));
      }
    } catch (error, stackTrace) {
      developer.log('Failed to start background music', name: 'AudioService', error: error, stackTrace: stackTrace);
    }
  }

  int _nextMusicIndex() {
    if (_musicAssets.length == 1) return 0;

    final previousIndex = _currentMusicIndex;
    var index = _random.nextInt(_musicAssets.length);
    if (index == previousIndex) {
      index = (index + 1 + _random.nextInt(_musicAssets.length - 1)) % _musicAssets.length;
    }
    _currentMusicIndex = index;
    return index;
  }

  Duration _remainingMusicDuration(SoundHandle handle) {
    final elapsed = SoLoud.instance.getPosition(handle);
    final remaining = _musicDuration! - elapsed;
    return remaining.isNegative ? Duration.zero : remaining;
  }

  void _scheduleNextTrack(SoundHandle handle, AudioSource source, Duration delay) {
    _nextTrackTimer?.cancel();
    _nextTrackTimer = Timer(delay, () {
      _nextTrackTimer = null;
      if (_musicPaused || _musicHandle != handle) {
        return;
      }
      _musicHandle = null;
      _musicSource = null;
      _musicDuration = null;
      unawaited(_disposeMusicSourceAndContinue(source));
    });
  }

  Future<void> _disposeMusicSourceAndContinue(AudioSource source) async {
    try {
      await SoLoud.instance.disposeSource(source);
    } catch (error, stackTrace) {
      developer.log(
        'Failed to dispose background music source',
        name: 'AudioService',
        error: error,
        stackTrace: stackTrace,
      );
    }
    await resumeMusic();
  }

  Future<void> _ensureEngineStarted() async {
    if (!SoLoud.instance.isInitialized) {
      await SoLoud.instance.init();
    }
  }

  Future<bool> _activatePlayback() async {
    final session = await AudioSession.instance;
    await session.configure(AudioSessionConfiguration.music());
    return session.setActive(true);
  }

  double _storedVolume(SharedPreferences preferences, String key, double fallback) {
    return (preferences.getDouble(key) ?? fallback).clamp(0.0, 1.0).toDouble();
  }

  Future<void> _saveVolume(String key, double volume) async {
    try {
      final preferences = _preferences ??= await SharedPreferences.getInstance();
      await preferences.setDouble(key, volume);
    } catch (error, stackTrace) {
      developer.log('Failed to save audio setting', name: 'AudioService', error: error, stackTrace: stackTrace);
    }
  }

  double _pitch(Character character) {
    final base = switch (character.gender) {
      'male' => -2.0,
      'female' => 2.0,
      _ => 0.0,
    };
    return base + _variation(character.id, 0.6);
  }

  double _speed(Character character) {
    final age = character.age.clamp(18, 80).toDouble();
    final ageBasedSpeed = 1.10 - (age - 18) * 0.18 / 62;
    return (ageBasedSpeed + _variation(character.id, 0.02)).clamp(0.92, 1.08).toDouble();
  }

  double _variation(String value, double range) {
    var hash = 0;
    for (final codeUnit in value.codeUnits) {
      hash = 31 * hash + codeUnit;
    }
    return (hash.abs() % 1000 / 1000 - 0.5) * 2 * range;
  }

  double get _effectiveMusicVolume => _musicVolume * _maxMusicVolume;
}
