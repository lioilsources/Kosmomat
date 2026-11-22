// Status and literature chips: two separate axes, never merged (EDITORIAL P2).
import 'package:flutter/material.dart';

import '../l10n/app_localizations.dart';
import 'labels.dart';

class AxisChips extends StatelessWidget {
  const AxisChips({super.key, this.status, this.literature});
  final String? status;
  final String? literature;

  static IconData statusIcon(String s) => switch (s) {
        'established' => Icons.verified_outlined,
        'strong' => Icons.check_circle_outline,
        'contested' => Icons.balance,
        'speculative' => Icons.help_outline,
        'refuted' => Icons.cancel_outlined,
        'historical' => Icons.history_edu,
        _ => Icons.psychology_alt_outlined,
      };

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Wrap(spacing: 8, runSpacing: 4, children: [
      if (status != null)
        ActionChip(
          // Icon + text so the label never depends on colour alone.
          avatar: Icon(statusIcon(status!), size: 18),
          label: Text('${l.statusAxis}: ${l.status(status!)}'),
          onPressed: () => _explain(context, l.statusAxis, l.status(status!), l.statusHelp(status!)),
        ),
      if (literature != null)
        ActionChip(
          avatar: const Icon(Icons.library_books_outlined, size: 18),
          label: Text('${l.litAxis}: ${l.literature(literature!)}'),
          onPressed: () => _explain(context, l.litAxis, l.literature(literature!), l.litHelp),
        ),
    ]);
  }

  void _explain(BuildContext context, String axis, String value, String help) => showModalBottomSheet(
        context: context,
        showDragHandle: true,
        builder: (ctx) => SafeArea(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(20, 0, 20, 24),
            child: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(axis, style: Theme.of(ctx).textTheme.labelLarge),
              Text(value, style: Theme.of(ctx).textTheme.titleLarge),
              const SizedBox(height: 8),
              Text(help, style: Theme.of(ctx).textTheme.bodyLarge),
            ]),
          ),
        ),
      );
}
