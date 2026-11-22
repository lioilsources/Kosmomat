// Bottom sheet with full details of a cited work (plan §6 "list zdroje").
import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

import '../data/models.dart';
import '../data/pack.dart';
import '../l10n/app_localizations.dart';
import 'labels.dart';

Future<void> showSourceSheet(BuildContext context, Pack pack, String workKey, {String? locator, required String lang}) {
  final w = pack.work(workKey);
  return showModalBottomSheet(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (ctx) => SafeArea(child: w == null ? Padding(padding: const EdgeInsets.all(24), child: Text(workKey)) : SourceDetails(work: w, locator: locator)),
  );
}

class SourceDetails extends StatelessWidget {
  const SourceDetails({super.key, required this.work, this.locator});
  final Work work;
  final String? locator;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final t = Theme.of(context).textTheme;
    final link = work.link;
    return SingleChildScrollView(
      padding: const EdgeInsets.fromLTRB(20, 0, 20, 20),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Wrap(spacing: 8, runSpacing: 4, children: [
          Chip(label: Text(l.tier(work.tier))),
          if (work.preprint) Chip(avatar: const Icon(Icons.info_outline, size: 18), label: Text(l.sourceUnreviewed)),
          if (work.retracted) Chip(avatar: const Icon(Icons.block, size: 18), label: Text(l.sourceRetracted)),
        ]),
        const SizedBox(height: 8),
        Text(work.title, style: t.titleMedium),
        const SizedBox(height: 4),
        Text('${work.authors.join(', ')} (${work.year})', style: t.bodyMedium),
        if (work.venue != null) Text(work.venue!, style: t.bodySmall),
        const SizedBox(height: 12),
        if (locator != null && locator!.isNotEmpty) _row(context, l.sourceLocator, locator!),
        if (work.country.isNotEmpty) _row(context, l.sourceCountry, work.country.join(', ')),
        if (work.doi != null) _row(context, 'DOI', work.doi!),
        if (work.arxiv != null) _row(context, 'arXiv', work.arxiv!),
        if (work.bibcode != null) _row(context, 'SciX', work.bibcode!),
        if (work.isbn != null) _row(context, 'ISBN', work.isbn!),
        if (!work.verified)
          Padding(padding: const EdgeInsets.only(top: 8), child: Text(l.sourceUnverified, style: t.bodySmall!.copyWith(fontStyle: FontStyle.italic))),
        if (link != null)
          Padding(
            padding: const EdgeInsets.only(top: 12),
            child: FilledButton.tonalIcon(
              onPressed: () => launchUrl(link, mode: LaunchMode.externalApplication),
              icon: const Icon(Icons.open_in_new),
              label: Text(l.sourceOpen),
            ),
          ),
      ]),
    );
  }

  Widget _row(BuildContext context, String k, String v) => Padding(
        padding: const EdgeInsets.only(bottom: 4),
        child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
          SizedBox(width: 130, child: Text(k, style: Theme.of(context).textTheme.labelMedium)),
          Expanded(child: SelectableText(v)),
        ]),
      );
}
