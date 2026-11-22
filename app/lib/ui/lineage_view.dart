// Vertical lineage timeline drawn with CustomPaint (plan §6 "Rodokmen"):
// origin and ancestors on top, this card in the middle, descendants below,
// competitors as side branches.
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../data/pack.dart';
import '../l10n/app_localizations.dart';
import 'source_sheet.dart';

class LineageView extends StatelessWidget {
  const LineageView({super.key, required this.pack, required this.nodeId, required this.lang});
  final Pack pack;
  final String nodeId;
  final String lang;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final lin = pack.lineage(nodeId);
    final rows = <_Row>[
      for (final o in lin.origins.reversed)
        _Row(
          kind: _Kind.origin,
          year: o.year,
          label: '${l.lineageOrigin}: ${pack.work(o.workKey)?.authors.firstOrNull ?? o.workKey}',
          sub: pack.work(o.workKey)?.title,
          onTap: () => showSourceSheet(context, pack, o.workKey, locator: o.locator, lang: lang),
        ),
      for (final a in lin.ancestors)
        _Row(kind: _Kind.ancestor, year: a.year, label: pack.title(a.id, lang), exists: a.exists, onTap: a.exists ? () => context.push('/n/${a.id}') : null),
      _Row(kind: _Kind.self, year: pack.node(nodeId)?.firstYear, label: pack.title(nodeId, lang)),
      for (final d in lin.descendants)
        _Row(kind: _Kind.descendant, year: d.year, label: pack.title(d.id, lang), exists: d.exists, onTap: d.exists ? () => context.push('/n/${d.id}') : null),
    ];
    if (rows.length == 1 && lin.competitors.isEmpty) return const SizedBox.shrink();
    final theme = Theme.of(context);
    return Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
      Text(l.lineage, style: theme.textTheme.titleMedium),
      const SizedBox(height: 8),
      for (var i = 0; i < rows.length; i++)
        IntrinsicHeight(
          child: Row(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
            SizedBox(
              width: 56,
              child: Padding(
                padding: const EdgeInsets.only(top: 12),
                child: Text(rows[i].year?.toString() ?? '', style: theme.textTheme.labelSmall, textAlign: TextAlign.right),
              ),
            ),
            SizedBox(
              width: 32,
              child: CustomPaint(
                painter: _RailPainter(
                  color: theme.colorScheme.outline,
                  dot: rows[i].kind == _Kind.self ? theme.colorScheme.primary : theme.colorScheme.outline,
                  big: rows[i].kind == _Kind.self,
                  top: i > 0,
                  bottom: i < rows.length - 1,
                  hollow: !rows[i].exists,
                ),
              ),
            ),
            Expanded(child: rows[i].build(context, l)),
          ]),
        ),
      if (lin.competitors.isNotEmpty) ...[
        const SizedBox(height: 12),
        Text(l.lineageCompetitors, style: theme.textTheme.labelLarge),
        const SizedBox(height: 4),
        Wrap(spacing: 8, runSpacing: 4, children: [
          for (final c in lin.competitors)
            ActionChip(
              avatar: const Icon(Icons.call_split, size: 18),
              label: Text(c.exists ? pack.title(c.id, lang) : '${c.id} · ${l.notWrittenYet}'),
              onPressed: c.exists ? () => context.push('/n/${c.id}') : null,
            ),
        ]),
      ],
    ]);
  }
}

enum _Kind { origin, ancestor, self, descendant }

class _Row {
  _Row({required this.kind, this.year, required this.label, this.sub, this.exists = true, this.onTap});
  final _Kind kind;
  final int? year;
  final String label;
  final String? sub;
  final bool exists;
  final VoidCallback? onTap;

  Widget build(BuildContext context, AppLocalizations l) {
    final t = Theme.of(context).textTheme;
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text(label, style: kind == _Kind.self ? t.titleSmall : t.bodyMedium),
          if (sub != null) Text(sub!, style: t.bodySmall, maxLines: 2, overflow: TextOverflow.ellipsis),
          if (!exists) Text(l.notWrittenYet, style: t.bodySmall!.copyWith(fontStyle: FontStyle.italic)),
        ]),
      ),
    );
  }
}

class _RailPainter extends CustomPainter {
  _RailPainter({required this.color, required this.dot, required this.big, required this.top, required this.bottom, required this.hollow});
  final Color color;
  final Color dot;
  final bool big;
  final bool top;
  final bool bottom;
  final bool hollow;

  @override
  void paint(Canvas canvas, Size size) {
    final x = size.width / 2;
    const y = 19.0;
    final line = Paint()
      ..color = color
      ..strokeWidth = 2;
    if (top) canvas.drawLine(Offset(x, 0), Offset(x, y), line);
    if (bottom) canvas.drawLine(Offset(x, y), Offset(x, size.height), line);
    final r = big ? 8.0 : 5.0;
    canvas.drawCircle(
        Offset(x, y),
        r,
        Paint()
          ..color = dot
          ..style = hollow ? PaintingStyle.stroke : PaintingStyle.fill
          ..strokeWidth = 2);
  }

  @override
  bool shouldRepaint(_RailPainter old) =>
      old.color != color || old.dot != dot || old.big != big || old.top != top || old.bottom != bottom || old.hollow != hollow;
}
