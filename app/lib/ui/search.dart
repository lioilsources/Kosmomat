// Offline FTS search with type and status filters (plan §6 "Hledat").
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../l10n/app_localizations.dart';
import '../state.dart';
import 'common.dart';
import 'labels.dart';

class SearchScreen extends ConsumerStatefulWidget {
  const SearchScreen({super.key, this.initial = ''});
  final String initial;

  @override
  ConsumerState<SearchScreen> createState() => _SearchScreenState();
}

class _SearchScreenState extends ConsumerState<SearchScreen> {
  late final _ctl = TextEditingController(text: widget.initial);
  String? _type;
  String? _status;

  @override
  void didUpdateWidget(SearchScreen old) {
    super.didUpdateWidget(old);
    if (widget.initial != old.initial && widget.initial.isNotEmpty) _ctl.text = widget.initial;
  }

  @override
  void dispose() {
    _ctl.dispose();
    super.dispose();
  }

  static const _types = ['theory', 'observation', 'claim', 'prediction', 'case', 'question', 'term', 'experiment', 'body', 'constellation'];
  static const _statuses = ['established', 'strong', 'contested', 'speculative', 'refuted', 'historical', 'philosophical'];

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final pack = ref.watch(packProvider);
    final lang = ref.watch(langProvider);
    final q = _ctl.text;
    final results = q.trim().isEmpty ? null : pack.search(q, lang, type: _type, status: _status);
    return Column(children: [
      Padding(
        padding: const EdgeInsets.fromLTRB(12, 12, 12, 4),
        child: TextField(
          controller: _ctl,
          autofocus: widget.initial.isEmpty,
          textInputAction: TextInputAction.search,
          decoration: InputDecoration(prefixIcon: const Icon(Icons.search), hintText: l.searchHint, border: const OutlineInputBorder()),
          onChanged: (_) => setState(() {}),
        ),
      ),
      SizedBox(
        height: 48,
        child: ListView(scrollDirection: Axis.horizontal, padding: const EdgeInsets.symmetric(horizontal: 12), children: [
          _dropdown<String>(l.filterAll, _type, _types, l.type, (v) => setState(() => _type = v)),
          const SizedBox(width: 8),
          _dropdown<String>(l.statusAxis, _status, _statuses, l.status, (v) => setState(() => _status = v)),
        ]),
      ),
      Expanded(
        child: results == null
            ? Center(child: Text(l.searchPrompt))
            : results.isEmpty
                ? Center(child: Text(l.searchEmpty))
                : ListView(children: [for (final h in results) NodeTile(head: h)]),
      ),
    ]);
  }

  Widget _dropdown<T>(String allLabel, T? value, List<T> items, String Function(T) label, ValueChanged<T?> onChanged) => DropdownButton<T?>(
        value: value,
        hint: Text(allLabel),
        underline: const SizedBox.shrink(),
        items: [
          DropdownMenuItem<T?>(value: null, child: Text(allLabel)),
          for (final i in items) DropdownMenuItem<T?>(value: i, child: Text(label(i))),
        ],
        onChanged: onChanged,
      );
}
