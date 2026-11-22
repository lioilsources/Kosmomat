import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'app.dart';
import 'data/pack_loader.dart';
import 'state.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  final (pack, prefs) = await (PackLoader.openBundled(), SharedPreferences.getInstance()).wait;
  runApp(ProviderScope(
    overrides: [packProvider.overrideWithValue(pack), prefsProvider.overrideWithValue(prefs)],
    child: const KosmoApp(),
  ));
}
