import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../l10n/app_localizations.dart';

/// Bottom navigation: Dnes · Otázky · Spory · Dějiny · Hledat; Moje top right.
class Shell extends StatelessWidget {
  const Shell({super.key, required this.shell});
  final StatefulNavigationShell shell;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context);
    return Scaffold(
      appBar: AppBar(
        title: Text(l.appTitle),
        actions: [
          IconButton(tooltip: l.mine, icon: const Icon(Icons.person_outline), onPressed: () => context.push('/mine')),
          IconButton(tooltip: l.settings, icon: const Icon(Icons.settings_outlined), onPressed: () => context.push('/settings')),
        ],
      ),
      body: shell,
      bottomNavigationBar: NavigationBar(
        selectedIndex: shell.currentIndex,
        onDestinationSelected: (i) => shell.goBranch(i, initialLocation: i == shell.currentIndex),
        destinations: [
          NavigationDestination(icon: const Icon(Icons.today_outlined), label: l.navToday),
          NavigationDestination(icon: const Icon(Icons.help_outline), label: l.navQuestions),
          NavigationDestination(icon: const Icon(Icons.balance), label: l.navDisputes),
          NavigationDestination(icon: const Icon(Icons.history_edu), label: l.navHistory),
          NavigationDestination(icon: const Icon(Icons.search), label: l.navSearch),
        ],
      ),
    );
  }
}
