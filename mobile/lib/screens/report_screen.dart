import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../services/api_service.dart';
import 'results_screen.dart';

class ReportScreen extends StatefulWidget {
  const ReportScreen({super.key});

  @override
  State<ReportScreen> createState() => _ReportScreenState();
}

class _ReportScreenState extends State<ReportScreen> {
  final _whoController = TextEditingController();
  final _whyController = TextEditingController();
  final _howController = TextEditingController();
  final _whenController = TextEditingController();
  final _evidenceController = TextEditingController();
  bool _submitting = false;

  @override
  void dispose() {
    _whoController.dispose();
    _whyController.dispose();
    _howController.dispose();
    _whenController.dispose();
    _evidenceController.dispose();
    super.dispose();
  }

  bool get _allFilled =>
      _whoController.text.trim().isNotEmpty &&
      _whyController.text.trim().isNotEmpty &&
      _howController.text.trim().isNotEmpty &&
      _whenController.text.trim().isNotEmpty &&
      _evidenceController.text.trim().isNotEmpty;

  Future<void> _submit() async {
    if (!_allFilled || _submitting) return;
    setState(() => _submitting = true);

    try {
      final result = await context.read<ApiService>().submitReport(
        who: _whoController.text.trim(),
        why: _whyController.text.trim(),
        how: _howController.text.trim(),
        when: _whenController.text.trim(),
        evidence: _evidenceController.text.trim(),
      );

      if (!mounted) return;

      Navigator.push(
        context,
        MaterialPageRoute(builder: (_) => ResultsScreen(result: result, playerReport: result.playerReport)),
      );
    } catch (e) {
      if (!mounted) return;
      setState(() => _submitting = false);
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Ошибка: $e')));
    }
  }

  Future<void> _confirmSubmit() async {
    if (!_allFilled || _submitting) return;
    final shouldSubmit = await showModalBottomSheet<bool>(
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
                  Text('Отправить отчёт?', style: Theme.of(sheetContext).textTheme.titleMedium),
                  const SizedBox(height: 8),
                  const Text('Расследование завершится. Изменить отчёт или вернуться к действиям уже не получится.'),
                  const SizedBox(height: 20),
                  SizedBox(
                    width: double.infinity,
                    child: FilledButton(
                      onPressed: () => Navigator.pop(sheetContext, true),
                      child: const Text('Завершить расследование'),
                    ),
                  ),
                  const SizedBox(height: 8),
                  SizedBox(
                    width: double.infinity,
                    child: TextButton(
                      onPressed: () => Navigator.pop(sheetContext, false),
                      child: const Text('Проверить ещё раз'),
                    ),
                  ),
                ],
              ),
            ),
          ),
    );
    if (shouldSubmit == true && mounted) {
      FocusScope.of(context).unfocus();
      await _submit();
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Финальный отчёт')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Изложите вашу версию событий', style: theme.textTheme.titleMedium),
            const SizedBox(height: 4),
            Text(
              'Заполните все разделы',
              style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurface.withAlpha(120)),
            ),
            const SizedBox(height: 20),
            _Field(label: 'Кто преступник?', hint: 'Майкл Браун...', controller: _whoController, onChanged: _onChanged),
            const SizedBox(height: 12),
            _Field(label: 'Почему?', hint: 'Хищение средств...', controller: _whyController, onChanged: _onChanged),
            const SizedBox(height: 12),
            _Field(
              label: 'Каким способом?',
              hint: 'Отравление цианидом...',
              controller: _howController,
              onChanged: _onChanged,
            ),
            const SizedBox(height: 12),
            _Field(label: 'В какое время?', hint: 'Около 22:15...', controller: _whenController, onChanged: _onChanged),
            const SizedBox(height: 12),
            _Field(
              label: 'Какие улики это подтверждают?',
              hint: 'Бокал с цианидом, фин. документы...',
              controller: _evidenceController,
              onChanged: _onChanged,
            ),
            const SizedBox(height: 24),
            SizedBox(
              width: double.infinity,
              height: 48,
              child: ElevatedButton.icon(
                onPressed: _allFilled && !_submitting ? _confirmSubmit : null,
                icon:
                    _submitting
                        ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2))
                        : const Icon(Icons.send),
                label: Text(_submitting ? 'Проверяем...' : 'Отправить отчёт'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _onChanged(String _) => setState(() {});
}

class _Field extends StatelessWidget {
  final String label;
  final String hint;
  final TextEditingController controller;
  final ValueChanged<String> onChanged;

  const _Field({required this.label, required this.hint, required this.controller, required this.onChanged});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: theme.textTheme.labelLarge?.copyWith(color: theme.colorScheme.primary)),
        const SizedBox(height: 4),
        TextField(
          controller: controller,
          maxLines: 2,
          onChanged: onChanged,
          decoration: InputDecoration(
            hintText: hint,
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
            contentPadding: const EdgeInsets.all(12),
          ),
        ),
      ],
    );
  }
}
