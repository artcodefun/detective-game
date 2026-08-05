import 'package:flutter/material.dart';

class ChatBubble extends StatelessWidget {
  final String text;
  final bool isPlayer;
  final String senderName;

  const ChatBubble({super.key, required this.text, required this.isPlayer, required this.senderName});

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: isPlayer ? MainAxisAlignment.end : MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          if (!isPlayer) ...[
            CircleAvatar(
              radius: 14,
              backgroundColor: colorScheme.primaryContainer,
              child: Text(
                senderName[0],
                style: TextStyle(fontSize: 12, color: colorScheme.onPrimaryContainer, fontWeight: FontWeight.bold),
              ),
            ),
            const SizedBox(width: 8),
          ],
          Flexible(
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
              decoration: BoxDecoration(
                color: isPlayer ? colorScheme.primary : colorScheme.surfaceContainerHighest,
                borderRadius: BorderRadius.only(
                  topLeft: const Radius.circular(16),
                  topRight: const Radius.circular(16),
                  bottomLeft: Radius.circular(isPlayer ? 16 : 4),
                  bottomRight: Radius.circular(isPlayer ? 4 : 16),
                ),
              ),
              child: Text(text, style: TextStyle(color: isPlayer ? colorScheme.onPrimary : colorScheme.onSurface)),
            ),
          ),
          if (isPlayer) const SizedBox(width: 8),
        ],
      ),
    );
  }
}
