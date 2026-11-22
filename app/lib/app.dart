import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'data/pack.dart';
import 'l10n/app_localizations.dart';
import 'state.dart';
import 'ui/lists.dart';
import 'ui/node_screen.dart';
import 'ui/search.dart';
import 'ui/settings.dart';
import 'ui/shell.dart';
import 'ui/today.dart';

/// Routes. `/n/<id>` matches the universal link `https://<web>/n/<id>`
/// and the scheme `kosmomat://node/<id>` (plan §5 I).
GoRouter buildRouter(Pack pack) => GoRouter(
      initialLocation: '/today',
      redirect: (context, state) {
        final uri = state.uri;
        // kosmomat://node/<id> arrives as host "node" with path "/<id>".
        if (uri.scheme == 'kosmomat' && uri.host == 'node') return '/n${uri.path}';
        if (uri.pathSegments.length == 2 && uri.pathSegments.first == 'n') {
          final id = uri.pathSegments[1];
          if (pack.node(id) != null) return null;
          final moved = pack.redirect(id);
          if (moved != null) return '/n/$moved';
          final q = id.contains('.') ? id.split('.').last.replaceAll('-', ' ') : id;
          return '/search?q=${Uri.encodeQueryComponent(q)}';
        }
        return null;
      },
      routes: [
        StatefulShellRoute.indexedStack(
          builder: (context, state, shell) => Shell(shell: shell),
          branches: [
            StatefulShellBranch(routes: [GoRoute(path: '/today', builder: (_, _) => const TodayScreen())]),
            StatefulShellBranch(routes: [
              GoRoute(path: '/questions', builder: (_, _) => TypeListScreen(types: const ['question'], intro: (l) => l.listQuestionsIntro)),
            ]),
            StatefulShellBranch(routes: [
              GoRoute(path: '/disputes', builder: (_, _) => TypeListScreen(types: const ['theory', 'claim', 'observation'], intro: (l) => l.listDisputesIntro)),
            ]),
            StatefulShellBranch(routes: [
              GoRoute(path: '/history', builder: (_, _) => TypeListScreen(types: const ['case', 'worldview', 'era'], intro: (l) => l.listHistoryIntro)),
            ]),
            StatefulShellBranch(routes: [
              GoRoute(path: '/search', builder: (_, s) => SearchScreen(initial: s.uri.queryParameters['q'] ?? '')),
            ]),
          ],
        ),
        GoRoute(path: '/n/:id', builder: (_, s) => NodeScreen(id: s.pathParameters['id']!)),
        GoRoute(path: '/mine', builder: (_, _) => const MineScreen()),
        GoRoute(path: '/settings', builder: (_, _) => const SettingsScreen()),
        GoRoute(path: '/credits', builder: (_, _) => const CreditsScreen()),
        GoRoute(path: '/editorial', builder: (_, _) => const EditorialScreen()),
      ],
    );

class KosmoApp extends ConsumerStatefulWidget {
  const KosmoApp({super.key});

  @override
  ConsumerState<KosmoApp> createState() => _KosmoAppState();
}

class _KosmoAppState extends ConsumerState<KosmoApp> {
  late final GoRouter _router = buildRouter(ref.read(packProvider));

  @override
  Widget build(BuildContext context) {
    final lang = ref.watch(langProvider);
    ThemeData theme(Brightness b) => ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF1F3A5F), brightness: b),
          useMaterial3: true,
        );
    return MaterialApp.router(
      onGenerateTitle: (c) => AppLocalizations.of(c).appTitle,
      debugShowCheckedModeBanner: false,
      theme: theme(Brightness.light),
      darkTheme: theme(Brightness.dark),
      locale: Locale(lang),
      supportedLocales: AppLocalizations.supportedLocales,
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      routerConfig: _router,
    );
  }
}
