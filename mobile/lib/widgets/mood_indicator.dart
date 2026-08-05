import 'package:flutter/material.dart';

import '../models/game_state.dart';

class MoodIndicator extends StatelessWidget {
  final TrustLevel trustLevel;

  const MoodIndicator({super.key, required this.trustLevel});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        _MoodDot(filled: trustLevel == TrustLevel.open, color: Colors.green),
        const SizedBox(width: 4),
        _MoodDot(filled: trustLevel == TrustLevel.reserved, color: Colors.yellow),
        const SizedBox(width: 4),
        _MoodDot(filled: trustLevel == TrustLevel.tense, color: Colors.orange),
        const SizedBox(width: 4),
        _MoodDot(filled: trustLevel == TrustLevel.closed, color: Colors.red),
      ],
    );
  }

  String get label {
    switch (trustLevel) {
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

  IconData get icon {
    switch (trustLevel) {
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
}

class _MoodDot extends StatelessWidget {
  final bool filled;
  final Color color;

  const _MoodDot({required this.filled, required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 10,
      height: 10,
      decoration: BoxDecoration(shape: BoxShape.circle, color: filled ? color : color.withAlpha(60)),
    );
  }
}
