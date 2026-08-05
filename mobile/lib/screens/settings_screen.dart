import 'package:flutter/material.dart';

class SettingsScreen extends StatefulWidget {
  const SettingsScreen({super.key});

  @override
  State<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends State<SettingsScreen> {
  double _soundVolume = 0.7;
  double _musicVolume = 0.5;

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
                Expanded(child: Slider(value: _soundVolume, onChanged: (v) => setState(() => _soundVolume = v))),
                Text('${(_soundVolume * 100).round()}%'),
              ],
            ),
            const SizedBox(height: 24),
            Text('Музыка', style: theme.textTheme.titleMedium),
            const SizedBox(height: 8),
            Row(
              children: [
                const Icon(Icons.music_note),
                Expanded(child: Slider(value: _musicVolume, onChanged: (v) => setState(() => _musicVolume = v))),
                Text('${(_musicVolume * 100).round()}%'),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
