import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../services/audio_service.dart';

class SettingsScreen extends StatefulWidget {
  const SettingsScreen({super.key});

  @override
  State<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends State<SettingsScreen> {
  late final AudioService _audio;
  late double _soundVolume;
  late double _musicVolume;

  @override
  void initState() {
    super.initState();
    _audio = context.read<AudioService>();
    _soundVolume = _audio.soundVolume;
    _musicVolume = _audio.musicVolume;
  }

  void _setSoundVolume(double volume) {
    _audio.setSoundVolume(volume);
    setState(() => _soundVolume = volume);
  }

  void _setMusicVolume(double volume) {
    _audio.setMusicVolume(volume);
    setState(() => _musicVolume = volume);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Настройки')),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Звук', style: theme.textTheme.titleMedium),
            const SizedBox(height: 8),
            Row(
              children: [
                const Icon(Icons.volume_up),
                Expanded(child: Slider(value: _soundVolume, onChanged: _setSoundVolume)),
                Text('${(_soundVolume * 100).round()}%'),
              ],
            ),
            const SizedBox(height: 24),
            Text('Музыка', style: theme.textTheme.titleMedium),
            const SizedBox(height: 8),
            Row(
              children: [
                const Icon(Icons.music_note),
                Expanded(child: Slider(value: _musicVolume, onChanged: _setMusicVolume)),
                Text('${(_musicVolume * 100).round()}%'),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
