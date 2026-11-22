// "Dnes": 3–5 cards, no endless feed (plan §6).
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../data/models.dart';
import '../l10n/app_localizations.dart';
import '../state.dart';
import 'common.dart';

/// Deterministic pick for a given day, so the card stays the same all day.
T? pickOfDay<T>(List<T> items, DateTime day) {
  if (items.isEmpty) return null;
  final n = DateTime.utc(day.year, day.month, day.day).millisecondsSinceEpoch ~/ Duration.millisecondsPerDay;
  return items[n % items.length];
}

class TodayScreen extends ConsumerWidget {
  const TodayScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context);
    final pack = ref.watch(packProvider);
    final lang = ref.watch(langProvider);
    final now = ref.watch(clockProvider)();

    final changed = <String, ChangelogEntry>{};
    for (final c in pack.changelog(limit: 30)) {
      changed.putIfAbsent(c.nodeId, () => c);
    }
    final changedHeads = {for (final h in pack.heads(lang, ids: changed.keys.take(3))) h.node.id: h};
    final caseOfDay = pickOfDay(pack.heads(lang, type: 'case'), now);
    final question = pickOfDay(pack.heads(lang, type: 'question'), now);
    final dispute = pickOfDay(pack.heads(lang, type: 'theory'), now);

    Widget section(String title, List<Widget> children) => Card(
          margin: const EdgeInsets.fromLTRB(12, 12, 12, 0),
          child: Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
            Padding(padding: const EdgeInsets.fromLTRB(16, 12, 16, 0), child: Text(title, style: Theme.of(context).textTheme.titleMedium)),
            ...children,
          ]),
        );

    return ListView(padding: const EdgeInsets.only(bottom: 24), children: [
      if (changedHeads.isNotEmpty)
        section(l.todayChanged, [
          for (final id in changed.keys.take(3))
            if (changedHeads[id] != null) NodeTile(head: changedHeads[id]!, overline: '${changed[id]!.date} · ${changed[id]!.note}'),
        ]),
      if (caseOfDay != null) section(l.todayCase, [NodeTile(head: caseOfDay)]),
      if (question != null) section(l.todayQuestion, [NodeTile(head: question)]),
      if (dispute != null)
        section(l.todayDispute, [
          NodeTile(head: dispute),
          Align(alignment: Alignment.centerRight, child: TextButton(onPressed: () => context.go('/disputes'), child: Text(l.navDisputes))),
        ]),
    ]);
  }
}
