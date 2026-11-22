import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../data/models.dart';
import '../l10n/app_localizations.dart';
import 'chips.dart';
import 'labels.dart';

/// One row in a list of cards.
class NodeTile extends StatelessWidget {
  const NodeTile({super.key, required this.head, this.overline});
  final NodeHead head;
  final String? overline;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    final n = head.node;
    final meta = [l.type(n.type), if (n.status != null) l.status(n.status!), if (n.literature != null) l.literature(n.literature!)].join(' · ');
    return ListTile(
      leading: n.status != null ? Icon(AxisChips.statusIcon(n.status!)) : const Icon(Icons.article_outlined),
      title: Text(head.title),
      subtitle: Text('${overline != null ? '$overline\n' : ''}$meta\n${head.summary}', maxLines: overline != null ? 4 : 3, overflow: TextOverflow.ellipsis),
      isThreeLine: true,
      onTap: () => context.push('/n/${n.id}'),
    );
  }
}

class IntroText extends StatelessWidget {
  const IntroText(this.text, {super.key});
  final String text;
  @override
  Widget build(BuildContext context) =>
      Padding(padding: const EdgeInsets.fromLTRB(16, 12, 16, 4), child: Text(text, style: Theme.of(context).textTheme.bodyMedium));
}
