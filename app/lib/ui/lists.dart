// Questions, Disputes, History tabs. Rich layouts (dispute columns, era timeline) come in M4.
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../l10n/app_localizations.dart';
import '../state.dart';
import 'common.dart';

class TypeListScreen extends ConsumerWidget {
  const TypeListScreen({super.key, required this.types, required this.intro});
  final List<String> types;
  final String Function(AppLocalizations) intro;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context);
    final pack = ref.watch(packProvider);
    final lang = ref.watch(langProvider);
    final heads = [for (final t in types) ...pack.heads(lang, type: t)];
    return ListView(children: [
      IntroText(intro(l)),
      if (heads.isEmpty) Padding(padding: const EdgeInsets.all(16), child: Text(l.listEmpty)),
      for (final h in heads) NodeTile(head: h),
    ]);
  }
}
