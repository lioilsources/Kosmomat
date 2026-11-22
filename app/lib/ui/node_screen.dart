// Card detail (plan §6 "Karta").
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../data/models.dart';
import '../data/pack.dart';
import '../l10n/app_localizations.dart';
import '../state.dart';
import 'blocks.dart';
import 'chips.dart';
import 'labels.dart';
import 'lineage_view.dart';
import 'source_sheet.dart';

class NodeScreen extends ConsumerStatefulWidget {
  const NodeScreen({super.key, required this.id});
  final String id;

  @override
  ConsumerState<NodeScreen> createState() => _NodeScreenState();
}

class _NodeScreenState extends ConsumerState<NodeScreen> {
  String? _level;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final pack = ref.watch(packProvider);
    final lang = ref.watch(langProvider);
    final settings = ref.watch(settingsProvider);
    final id = pack.node(widget.id) == null ? (pack.redirect(widget.id) ?? widget.id) : widget.id;
    final node = pack.node(id);
    if (node == null) return UnknownNodeScreen(id: widget.id);

    final levels = pack.levels(id, lang);
    final level = _level ?? (levels.contains(settings.level) ? settings.level : 'lay');
    final doc = pack.doc(id, lang, level)!;
    final now = ref.watch(clockProvider)();
    final saved = settings.bookmarks.contains(id);
    final t = Theme.of(context).textTheme;

    return Scaffold(
      appBar: AppBar(
        title: Text(l.type(node.type)),
        actions: [
          IconButton(
            tooltip: saved ? l.bookmarked : l.bookmark,
            icon: Icon(saved ? Icons.star : Icons.star_border),
            onPressed: () => ref.read(settingsProvider.notifier).toggleBookmark(id),
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.fromLTRB(16, 8, 16, 32),
        children: [
          Text(doc.title, style: t.headlineSmall),
          const SizedBox(height: 8),
          Text(doc.summary, style: t.bodyLarge!.copyWith(fontStyle: FontStyle.italic)),
          const SizedBox(height: 12),
          AxisChips(status: node.status, literature: node.literature),
          if (node.awaitsReview(now))
            Padding(
              padding: const EdgeInsets.only(top: 8),
              child: Row(children: [
                const Icon(Icons.pending_outlined, size: 18),
                const SizedBox(width: 6),
                Expanded(child: Text(node.verifiedBy == 'human' ? l.awaitsReview : l.awaitsReviewDraft, style: t.bodySmall)),
              ]),
            ),
          if (levels.length > 1)
            Padding(
              padding: const EdgeInsets.only(top: 12),
              child: SegmentedButton<String>(
                segments: [for (final lv in levels) ButtonSegment(value: lv, label: Text(l.level(lv)))],
                selected: {level},
                onSelectionChanged: (s) => setState(() => _level = s.first),
              ),
            ),
          const SizedBox(height: 8),
          BlocksView(blocks: doc.blocks, pack: pack, lang: lang, nodeId: id),
          if (node.type == 'theory' || node.type == 'prediction') ...[
            const SizedBox(height: 24),
            PredictionBook(pack: pack, theoryId: id, lang: lang),
          ],
          const SizedBox(height: 24),
          LineageView(pack: pack, nodeId: id, lang: lang),
          _EdgeSection(pack: pack, lang: lang, title: l.supportedBy, edges: pack.edgesFrom(id).where((e) => e.rel == 'supported_by').toList()),
          _EdgeSection(pack: pack, lang: lang, title: l.opposedBy, edges: pack.edgesFrom(id).where((e) => opposingRels.contains(e.rel)).toList()),
          const SizedBox(height: 24),
          _Sources(pack: pack, id: id, lang: lang),
          const SizedBox(height: 24),
          Text(l.cardHistory, style: t.titleMedium),
          for (final c in pack.changelog(nodeId: id)) ListTile(dense: true, contentPadding: EdgeInsets.zero, title: Text(c.note), subtitle: Text(c.date)),
          Text(l.lastVerified(node.lastVerified), style: t.bodySmall),
        ],
      ),
    );
  }
}

class _EdgeSection extends StatelessWidget {
  const _EdgeSection({required this.pack, required this.lang, required this.title, required this.edges});
  final Pack pack;
  final String lang;
  final String title;
  final List<Edge> edges;

  @override
  Widget build(BuildContext context) {
    if (edges.isEmpty) return const SizedBox.shrink();
    final l = AppLocalizations.of(context);
    return Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
      const SizedBox(height: 24),
      Text(title, style: Theme.of(context).textTheme.titleMedium),
      for (final e in edges)
        ListTile(
          contentPadding: EdgeInsets.zero,
          title: Text(e.dstExists ? pack.title(e.dst, lang) : e.dst),
          subtitle: Text([
            if (!e.dstExists) l.notWrittenYet,
            if (e.note.isNotEmpty) pick(e.note, lang),
            if (e.workKey != null) '${pack.work(e.workKey!)?.authors.firstOrNull ?? e.workKey} ${pack.work(e.workKey!)?.year ?? ''}, ${e.locator ?? ''}',
          ].join(' · ')),
          onTap: e.dstExists
              ? () => context.push('/n/${e.dst}')
              : e.workKey != null
                  ? () => showSourceSheet(context, pack, e.workKey!, locator: e.locator, lang: lang)
                  : null,
        ),
    ]);
  }
}

class _Sources extends StatelessWidget {
  const _Sources({required this.pack, required this.id, required this.lang});
  final Pack pack;
  final String id;
  final String lang;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final src = pack.sources(id);
    return Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
      Text(l.sources, style: Theme.of(context).textTheme.titleMedium),
      for (final (w, locs) in src)
        ListTile(
          contentPadding: EdgeInsets.zero,
          leading: CircleAvatar(radius: 16, child: Text(w.tier, style: Theme.of(context).textTheme.labelSmall)),
          title: Text('${w.authors.firstOrNull ?? ''} (${w.year})'),
          subtitle: Text(w.title, maxLines: 2, overflow: TextOverflow.ellipsis),
          onTap: () => showSourceSheet(context, pack, w.key, locator: locs.join('; '), lang: lang),
        ),
    ]);
  }
}

/// Prediction book table: row = prediction, outcome + year + before/after measurement.
class PredictionBook extends StatelessWidget {
  const PredictionBook({super.key, required this.pack, required this.theoryId, required this.lang});
  final Pack pack;
  final String theoryId;
  final String lang;

  static IconData icon(String o) => switch (o) {
        'confirmed' => Icons.check_circle,
        'refuted' => Icons.cancel,
        'partial' => Icons.adjust,
        'open' => Icons.radio_button_unchecked,
        _ => Icons.block,
      };

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final preds = pack.predictions(theoryId);
    if (preds.isEmpty) return const SizedBox.shrink();
    final t = Theme.of(context).textTheme;
    return Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
      Text(l.predictionBook, style: t.titleMedium),
      const SizedBox(height: 8),
      for (final p in preds)
        Card(
          margin: const EdgeInsets.only(bottom: 8),
          child: InkWell(
            onTap: p.hasCard && p.id != theoryId ? () => context.push('/n/${p.id}') : null,
            child: Padding(
              padding: const EdgeInsets.all(12),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Wrap(spacing: 8, runSpacing: 2, crossAxisAlignment: WrapCrossAlignment.center, children: [
                  Icon(icon(p.outcome), size: 20),
                  Text(l.outcome(p.outcome), style: t.labelLarge),
                  Text(
                    '· ${p.retrodiction ? l.predRetrodiction : (p.madeBeforeObservation ? l.predBefore : l.predAfter)}',
                    style: t.labelSmall,
                  ),
                ]),
                const SizedBox(height: 6),
                Text(pick(p.statement, lang), style: t.bodyMedium),
                const SizedBox(height: 6),
                Wrap(spacing: 12, children: [
                  if (p.madeYear != null)
                    _ref(context, '${l.predMade} ${p.madeYear}', p.madeWork, p.madeLoc),
                  if (p.testYear != null)
                    _ref(context, '${l.predTested} ${p.testYear}', p.testWork, p.testLoc),
                ]),
                if (p.note.isNotEmpty) Padding(padding: const EdgeInsets.only(top: 6), child: Text(pick(p.note, lang), style: t.bodySmall)),
              ]),
            ),
          ),
        ),
    ]);
  }

  Widget _ref(BuildContext context, String label, String? work, String? loc) => TextButton(
        style: TextButton.styleFrom(padding: EdgeInsets.zero, minimumSize: const Size(0, 32)),
        onPressed: work == null ? null : () => showSourceSheet(context, pack, work, locator: loc, lang: lang),
        child: Text(label),
      );
}

class UnknownNodeScreen extends StatelessWidget {
  const UnknownNodeScreen({super.key, required this.id});
  final String id;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final query = id.contains('.') ? id.split('.').last.replaceAll('-', ' ') : id;
    return Scaffold(
      appBar: AppBar(),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(mainAxisSize: MainAxisSize.min, children: [
            Text(l.unknownNode(id), textAlign: TextAlign.center),
            const SizedBox(height: 12),
            FilledButton.icon(onPressed: () => context.go('/search?q=${Uri.encodeQueryComponent(query)}'), icon: const Icon(Icons.search), label: Text(l.navSearch)),
          ]),
        ),
      ),
    );
  }
}
