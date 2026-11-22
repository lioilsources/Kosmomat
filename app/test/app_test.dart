import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:kosmomat/app.dart';
import 'package:kosmomat/data/pack.dart';
import 'package:kosmomat/state.dart';
import 'package:kosmomat/ui/source_sheet.dart';
import 'package:shared_preferences/shared_preferences.dart';

late Pack pack;

Future<void> pumpApp(WidgetTester tester, {String lang = 'cs', double textScale = 1}) async {
  SharedPreferences.setMockInitialValues({'lang': lang});
  final prefs = await SharedPreferences.getInstance();
  tester.view.physicalSize = const Size(1080, 2400);
  tester.view.devicePixelRatio = 2.6;
  addTearDown(tester.view.reset);
  await tester.pumpWidget(ProviderScope(
    overrides: [
      packProvider.overrideWithValue(pack),
      prefsProvider.overrideWithValue(prefs),
      clockProvider.overrideWithValue(() => DateTime(2026, 11, 22)),
    ],
    child: MediaQuery(
      data: MediaQueryData(textScaler: TextScaler.linear(textScale), size: const Size(415, 923)),
      child: const KosmoApp(),
    ),
  ));
  await tester.pumpAndSettle();
}

GoRouter routerOf(WidgetTester tester) => GoRouter.of(tester.element(find.byType(Scaffold).first));

void main() {
  setUpAll(() => pack = Pack.fromBytes(File('assets/pack/pack.sqlite').readAsBytesSync()));

  testWidgets('today shows cards', (tester) async {
    await pumpApp(tester);
    expect(find.text('Co se změnilo'), findsOneWidget);
    expect(find.text('Případ dne'), findsOneWidget);
  });

  for (final (lang, scale) in [('cs', 1.0), ('en', 1.0), ('cs', 1.6)]) {
    testWidgets('every card renders without errors ($lang, text ×$scale)', (tester) async {
      await pumpApp(tester, lang: lang, textScale: scale);
      for (final h in pack.heads(lang)) {
        routerOf(tester).push('/n/${h.node.id}');
        await tester.pumpAndSettle();
        expect(find.text(h.title), findsWidgets, reason: h.node.id);
        // Scroll to the bottom so lazily built sections (lineage, sources) render too.
        await tester.fling(find.byType(ListView).last, const Offset(0, -20000), 5000);
        await tester.pumpAndSettle();
        expect(tester.takeException(), isNull, reason: h.node.id);
        routerOf(tester).pop();
        await tester.pumpAndSettle();
      }
    });
  }

  testWidgets('citation opens the source sheet', (tester) async {
    await pumpApp(tester);
    routerOf(tester).push('/n/theory.electric-universe');
    await tester.pumpAndSettle();
    await tester.tap(find.text('[1]').first);
    await tester.pumpAndSettle();
    expect(find.byType(SourceDetails), findsOneWidget);
    expect(find.textContaining('T5'), findsWidgets);
  });

  testWidgets('theory lineage reaches the origin', (tester) async {
    await pumpApp(tester);
    routerOf(tester).push('/n/theory.mond');
    await tester.pumpAndSettle();
    await tester.scrollUntilVisible(find.textContaining('Počátek'), 300, scrollable: find.byType(Scrollable).last);
    expect(find.textContaining('Milgrom'), findsWidgets);
  });

  testWidgets('unknown id opens search prefilled', (tester) async {
    await pumpApp(tester);
    routerOf(tester).go('/n/body.saturn');
    await tester.pumpAndSettle();
    expect(find.widgetWithText(TextField, 'saturn'), findsOneWidget);
  });

  testWidgets('scheme kosmomat://node/<id> opens the card', (tester) async {
    await pumpApp(tester);
    routerOf(tester).go('kosmomat://node/body.mars');
    await tester.pumpAndSettle();
    expect(find.text('Mars'), findsWidgets);
  });

  testWidgets('search finds cards offline', (tester) async {
    await pumpApp(tester);
    routerOf(tester).go('/search');
    await tester.pumpAndSettle();
    await tester.enterText(find.byType(TextField), 'neutrina');
    await tester.pumpAndSettle();
    expect(find.byType(ListTile), findsWidgets);
  });
}
