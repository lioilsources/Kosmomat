// Copies the bundled pack out of the app bundle so SQLite can open it.
import 'dart:convert';
import 'dart:io';

import 'package:flutter/services.dart';
import 'package:path_provider/path_provider.dart';

import 'pack.dart';

class PackLoader {
  static const _asset = 'assets/pack/pack.sqlite';
  static const _manifest = 'assets/pack/manifest.json';

  /// Opens the bundled pack. The file is copied once per pack version
  /// (sha256 from the bundled manifest), so a cold start only opens a file.
  static Future<Pack> openBundled() async {
    final manifest = jsonDecode(await rootBundle.loadString(_manifest)) as Map<String, dynamic>;
    final sha = manifest['sha256'] as String;
    final dir = Directory('${(await getApplicationSupportDirectory()).path}/pack');
    final file = File('${dir.path}/bundled-${sha.substring(0, 16)}.sqlite');
    if (!file.existsSync()) {
      await dir.create(recursive: true);
      for (final old in dir.listSync().whereType<File>()) {
        if (old.path.contains('bundled-')) old.deleteSync();
      }
      final data = await rootBundle.load(_asset);
      final tmp = File('${file.path}.tmp');
      await tmp.writeAsBytes(data.buffer.asUint8List(data.offsetInBytes, data.lengthInBytes), flush: true);
      await tmp.rename(file.path);
    }
    // Downloaded newer packs (M6) will be picked here before the bundled one.
    return Pack.open(file.path);
  }
}
