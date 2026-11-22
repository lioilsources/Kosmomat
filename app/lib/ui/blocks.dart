// Renders pack blocks (never Markdown): paragraphs with tappable citations and links.
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart' hide Block;

import '../data/models.dart';
import '../data/pack.dart';
import 'source_sheet.dart';

class BlocksView extends StatelessWidget {
  const BlocksView({super.key, required this.blocks, required this.pack, required this.lang, required this.nodeId});

  final List<Block> blocks;
  final Pack pack;
  final String lang;
  final String nodeId;

  @override
  Widget build(BuildContext context) {
    final cites = _CiteNumbers(blocks);
    final theme = Theme.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        for (final b in blocks)
          switch (b.t) {
            'h' => Padding(
                padding: const EdgeInsets.only(top: 20, bottom: 6),
                child: Semantics(header: true, child: Text(b.text ?? '', style: theme.textTheme.titleMedium)),
              ),
            'ul' => Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  for (final it in b.items)
                    Padding(
                      padding: const EdgeInsets.only(bottom: 8),
                      child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
                        Text('•  ', style: theme.textTheme.bodyLarge),
                        Expanded(child: _rich(context, it, cites)),
                      ]),
                    ),
                ],
              ),
            'q' => Container(
                margin: const EdgeInsets.symmetric(vertical: 8),
                padding: const EdgeInsets.only(left: 12),
                decoration: BoxDecoration(border: Border(left: BorderSide(color: theme.colorScheme.primary, width: 3))),
                child: _rich(context, b.spans, cites, italic: true),
              ),
            _ => Padding(padding: const EdgeInsets.only(bottom: 10), child: _rich(context, b.spans, cites)),
          },
      ],
    );
  }

  Widget _rich(BuildContext context, List<Span> spans, _CiteNumbers cites, {bool italic = false}) {
    final theme = Theme.of(context);
    final base = theme.textTheme.bodyLarge!.copyWith(height: 1.5, fontStyle: italic ? FontStyle.italic : null);
    final link = base.copyWith(color: theme.colorScheme.primary, decoration: TextDecoration.underline);
    return Text.rich(
      TextSpan(style: base, children: [
        for (final s in spans)
          switch (s.t) {
            'cite' => WidgetSpan(
                alignment: PlaceholderAlignment.top,
                child: Semantics(
                  button: true,
                  label: 'zdroj ${cites.of(s.w!)}',
                  child: InkWell(
                    onTap: () => showSourceSheet(context, pack, s.w!, locator: s.loc, lang: lang),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 2),
                      child: Text('[${cites.of(s.w!)}]',
                          style: theme.textTheme.labelSmall!.copyWith(color: theme.colorScheme.primary, fontWeight: FontWeight.bold)),
                    ),
                  ),
                ),
              ),
            'link' || 'term' => TextSpan(
                text: s.v.isNotEmpty ? s.v : pack.title(s.to!, lang),
                style: s.t == 'term' ? link.copyWith(decorationStyle: TextDecorationStyle.dotted) : link,
                recognizer: TapGestureRecognizer()..onTap = () => context.push('/n/${s.to}'),
              ),
            _ => TextSpan(
                text: s.v,
                style: TextStyle(fontWeight: s.b ? FontWeight.bold : null, fontStyle: s.i ? FontStyle.italic : null),
              ),
          }
      ]),
    );
  }
}

/// Numbers citations by first appearance on the card: [1], [2]…
class _CiteNumbers {
  _CiteNumbers(List<Block> blocks) {
    for (final b in blocks) {
      for (final s in [...b.spans, ...b.items.expand((i) => i)]) {
        if (s.t == 'cite') _n.putIfAbsent(s.w!, () => _n.length + 1);
      }
    }
  }
  final _n = <String, int>{};
  int of(String w) => _n[w] ?? 0;
}
