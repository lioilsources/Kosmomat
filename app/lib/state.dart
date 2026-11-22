import 'package:flutter/widgets.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'data/pack.dart';

/// The open content pack; overridden in main().
final packProvider = Provider<Pack>((ref) => throw UnimplementedError('pack not opened'));

/// Local preferences; overridden in main().
final prefsProvider = Provider<SharedPreferences>((ref) => throw UnimplementedError('prefs not loaded'));

/// Clock, overridable in tests.
final clockProvider = Provider<DateTime Function()>((ref) => DateTime.now);

class Settings {
  const Settings({this.lang, this.level = 'lay', this.bookmarks = const []});

  /// Content language; null = follow the system ("cs" for Czech, otherwise "en").
  final String? lang;
  final String level;
  final List<String> bookmarks;

  Settings copyWith({String? Function()? lang, String? level, List<String>? bookmarks}) =>
      Settings(lang: lang != null ? lang() : this.lang, level: level ?? this.level, bookmarks: bookmarks ?? this.bookmarks);
}

class SettingsNotifier extends Notifier<Settings> {
  SharedPreferences get _p => ref.read(prefsProvider);

  @override
  Settings build() {
    final p = ref.watch(prefsProvider);
    return Settings(
      lang: p.getString('lang'),
      level: p.getString('level') ?? 'lay',
      bookmarks: p.getStringList('bookmarks') ?? const [],
    );
  }

  void setLang(String? lang) {
    lang == null ? _p.remove('lang') : _p.setString('lang', lang);
    state = state.copyWith(lang: () => lang);
  }

  void setLevel(String level) {
    _p.setString('level', level);
    state = state.copyWith(level: level);
  }

  void toggleBookmark(String id) {
    final b = [...state.bookmarks];
    b.contains(id) ? b.remove(id) : b.insert(0, id);
    _p.setStringList('bookmarks', b);
    state = state.copyWith(bookmarks: b);
  }
}

final settingsProvider = NotifierProvider<SettingsNotifier, Settings>(SettingsNotifier.new);

/// Content language resolved from settings and the system locale.
String contentLang(Settings s, Locale system) => s.lang ?? (system.languageCode == 'cs' ? 'cs' : 'en');

/// Content and UI language ("cs" | "en").
final langProvider = Provider<String>((ref) {
  final s = ref.watch(settingsProvider);
  return contentLang(s, WidgetsBinding.instance.platformDispatcher.locale);
});
