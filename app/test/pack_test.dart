import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:kosmomat/data/pack.dart';

void main() {
  late Pack pack;
  setUpAll(() => pack = Pack.fromBytes(File('assets/pack/pack.sqlite').readAsBytesSync()));
  tearDownAll(() => pack.dispose());

  test('all nodes readable in cs and en', () {
    for (final lang in ['cs', 'en']) {
      final heads = pack.heads(lang);
      expect(heads.length, greaterThanOrEqualTo(14));
      for (final h in heads) {
        final d = pack.doc(h.node.id, lang, 'lay');
        expect(d, isNotNull, reason: h.node.id);
        expect(d!.blocks, isNotEmpty, reason: h.node.id);
      }
    }
  });

  test('every citation resolves to a work', () {
    for (final h in pack.heads('cs')) {
      for (final lang in ['cs', 'en']) {
        for (final b in pack.doc(h.node.id, lang, 'lay')!.blocks) {
          for (final s in [...b.spans, ...b.items.expand((i) => i)]) {
            if (s.t == 'cite') expect(pack.work(s.w!), isNotNull, reason: '${h.node.id} ${s.w}');
          }
        }
      }
    }
  });

  test('every claim and theory reaches its origin', () {
    for (final h in pack.heads('cs')) {
      if (h.node.type != 'claim' && h.node.type != 'theory') continue;
      expect(pack.lineage(h.node.id).reachesOrigin, isTrue, reason: h.node.id);
    }
  });

  test('fts5 search ignores diacritics and handles odd input', () {
    expect(pack.search('neutrina', 'cs'), isNotEmpty);
    expect(pack.search('temna hmota', 'cs').map((h) => h.node.id), contains('term.dark-matter'));
    expect(pack.search('"; DROP', 'cs'), isA<List>());
    expect(pack.search('   ', 'cs'), isEmpty);
    expect(pack.search('mars', 'en', type: 'body').single.node.id, 'body.mars');
  });

  test('prediction book of EU has refuted entries', () {
    final p = pack.predictions('theory.electric-universe');
    expect(p.length, greaterThanOrEqualTo(3));
    expect(p.where((x) => x.outcome == 'refuted'), isNotEmpty);
  });
}
