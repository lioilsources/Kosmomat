import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../l10n/app_localizations.dart';
import '../state.dart';
import 'common.dart';
import 'labels.dart';
import 'source_sheet.dart';

class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context);
    final s = ref.watch(settingsProvider);
    final n = ref.read(settingsProvider.notifier);
    final pack = ref.watch(packProvider);
    return Scaffold(
      appBar: AppBar(title: Text(l.settings)),
      body: ListView(children: [
        ListTile(
          title: Text(l.settingsLanguage),
          trailing: DropdownButton<String?>(
            value: s.lang,
            underline: const SizedBox.shrink(),
            items: [
              DropdownMenuItem(value: null, child: Text(l.settingsLanguageSystem)),
              const DropdownMenuItem(value: 'cs', child: Text('Čeština')),
              const DropdownMenuItem(value: 'en', child: Text('English')),
            ],
            onChanged: n.setLang,
          ),
        ),
        ListTile(
          title: Text(l.settingsLevel),
          trailing: DropdownButton<String>(
            value: s.level,
            underline: const SizedBox.shrink(),
            items: [for (final lv in ['kid', 'lay', 'student']) DropdownMenuItem(value: lv, child: Text(l.level(lv)))],
            onChanged: (v) => n.setLevel(v!),
          ),
        ),
        const Divider(),
        ListTile(leading: const Icon(Icons.gavel), title: Text(l.settingsEditorial), onTap: () => context.push('/editorial')),
        ListTile(leading: const Icon(Icons.copyright), title: Text(l.settingsCredits), onTap: () => context.push('/credits')),
        ListTile(leading: const Icon(Icons.inventory_2_outlined), title: Text(l.settingsPack), subtitle: Text(pack.version)),
      ]),
    );
  }
}

class MineScreen extends ConsumerWidget {
  const MineScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context);
    final pack = ref.watch(packProvider);
    final lang = ref.watch(langProvider);
    final ids = ref.watch(settingsProvider).bookmarks;
    final heads = {for (final h in pack.heads(lang, ids: ids)) h.node.id: h};
    return Scaffold(
      appBar: AppBar(title: Text(l.mine)),
      body: ListView(children: [
        IntroText(l.mineBookmarks),
        if (ids.isEmpty) Padding(padding: const EdgeInsets.all(16), child: Text(l.mineEmpty)),
        for (final id in ids)
          if (heads[id] != null) NodeTile(head: heads[id]!),
      ]),
    );
  }
}

class CreditsScreen extends ConsumerWidget {
  const CreditsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context);
    final pack = ref.watch(packProvider);
    final lang = ref.watch(langProvider);
    final media = pack.mediaCredits();
    final works = {for (final h in pack.heads(lang)) for (final (w, _) in pack.sources(h.node.id)) w.key: w}.values.toList()
      ..sort((a, b) => a.year.compareTo(b.year));
    final t = Theme.of(context).textTheme;
    return Scaffold(
      appBar: AppBar(title: Text(l.settingsCredits)),
      body: ListView(padding: const EdgeInsets.symmetric(vertical: 8), children: [
        Padding(padding: const EdgeInsets.all(16), child: Text(l.creditsMedia, style: t.titleMedium)),
        if (media.isEmpty) Padding(padding: const EdgeInsets.symmetric(horizontal: 16), child: Text(l.creditsMediaNone)),
        for (final m in media) ListTile(title: Text(m.credit), subtitle: Text('${m.license} · ${m.sourceUrl}')),
        Padding(padding: const EdgeInsets.all(16), child: Text(l.creditsSoftware, style: t.titleMedium)),
        ListTile(title: Text(l.creditsEphemeris)),
        ListTile(
          leading: const Icon(Icons.description_outlined),
          title: Text(l.creditsSoftware),
          onTap: () => showLicensePage(context: context, applicationName: l.appTitle, applicationVersion: pack.version),
        ),
        Padding(padding: const EdgeInsets.all(16), child: Text('${l.creditsWorks} (${works.length})', style: t.titleMedium)),
        for (final w in works)
          ListTile(
            dense: true,
            title: Text('${w.authors.firstOrNull ?? ''} (${w.year}) · ${w.tier}'),
            subtitle: Text(w.title, maxLines: 2, overflow: TextOverflow.ellipsis),
            onTap: () => showSourceSheet(context, pack, w.key, lang: lang),
          ),
      ]),
    );
  }
}

/// Shows docs/EDITORIAL.md shipped in the pack. A tiny line renderer:
/// headings, bullets and tables as plain rows — no Markdown package.
class EditorialScreen extends ConsumerWidget {
  const EditorialScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context);
    final pack = ref.watch(packProvider);
    final lang = ref.watch(langProvider);
    final md = pack.meta('editorial_$lang') ?? pack.meta('editorial_cs') ?? '';
    final t = Theme.of(context).textTheme;
    final widgets = <Widget>[];
    for (final raw in md.split('\n')) {
      final line = raw.trim();
      if (line.isEmpty || RegExp(r'^\|[-| ]+\|$').hasMatch(line)) continue;
      final clean = line.replaceAll('**', '').replaceAll('`', '').replaceAllMapped(RegExp(r'\[([^\]]+)\]\([^)]+\)'), (m) => m[1]!);
      if (line.startsWith('#')) {
        final level = line.indexOf(' ');
        widgets.add(Padding(
          padding: const EdgeInsets.only(top: 16, bottom: 4),
          child: Text(clean.substring(level + 1), style: level == 1 ? t.headlineSmall : t.titleMedium),
        ));
      } else if (line.startsWith('- ')) {
        widgets.add(Padding(padding: const EdgeInsets.only(left: 8, bottom: 4), child: Text('•  ${clean.substring(2)}', style: t.bodyMedium)));
      } else if (line.startsWith('|')) {
        final cells = clean.split('|').map((c) => c.trim()).where((c) => c.isNotEmpty).toList();
        widgets.add(Padding(padding: const EdgeInsets.only(bottom: 4), child: Text(cells.join(' — '), style: t.bodySmall)));
      } else {
        widgets.add(Padding(padding: const EdgeInsets.only(bottom: 6), child: Text(clean, style: t.bodyMedium)));
      }
    }
    return Scaffold(
      appBar: AppBar(title: Text(l.settingsEditorial)),
      body: ListView(padding: const EdgeInsets.fromLTRB(16, 0, 16, 32), children: widgets),
    );
  }
}
