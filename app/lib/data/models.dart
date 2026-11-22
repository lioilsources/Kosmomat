// Data classes mirroring the content pack (internal/pack/schema.sql).
import 'dart:convert';

Map<String, String> _i18n(Object? v) {
  if (v == null) return const {};
  final m = (v is String ? jsonDecode(v) : v) as Map<String, dynamic>;
  return m.map((k, v) => MapEntry(k, v as String));
}

List<String> _strings(Object? v) {
  if (v == null) return const [];
  final l = (v is String ? jsonDecode(v) : v) as List<dynamic>;
  return l.cast<String>();
}

/// Pick a language with English fallback.
String pick(Map<String, String> m, String lang) => m[lang] ?? m['en'] ?? m.values.firstOrNull ?? '';

class NodeRow {
  NodeRow({
    required this.id,
    required this.type,
    this.status,
    this.literature,
    this.testability,
    this.caseKind,
    this.astromat = false,
    this.firstYear,
    this.firstBy = const [],
    this.firstWork,
    this.firstLoc,
    required this.lastVerified,
    required this.verifiedBy,
    this.extra = const {},
  });

  final String id;
  final String type;
  final String? status;
  final String? literature;
  final String? testability;
  final String? caseKind;
  final bool astromat;
  final int? firstYear;
  final List<String> firstBy;
  final String? firstWork;
  final String? firstLoc;
  final String lastVerified;
  final String verifiedBy;
  final Map<String, dynamic> extra;

  factory NodeRow.fromRow(Map<String, Object?> r) => NodeRow(
        id: r['id']! as String,
        type: r['type']! as String,
        status: r['status'] as String?,
        literature: r['literature'] as String?,
        testability: r['testability'] as String?,
        caseKind: r['case_kind'] as String?,
        astromat: (r['astromat'] as int? ?? 0) == 1,
        firstYear: r['first_year'] as int?,
        firstBy: _strings(r['first_by']),
        firstWork: r['first_work'] as String?,
        firstLoc: r['first_loc'] as String?,
        lastVerified: r['last_verified']! as String,
        verifiedBy: r['verified_by']! as String,
        extra: r['extra_json'] == null ? const {} : jsonDecode(r['extra_json']! as String) as Map<String, dynamic>,
      );

  /// Card awaits review: written by the pipeline, or contested/speculative and
  /// last verified more than 12 months ago (EDITORIAL P10).
  bool awaitsReview(DateTime now) {
    if (verifiedBy != 'human') return true;
    if (status != 'contested' && status != 'speculative') return false;
    final d = DateTime.tryParse(lastVerified);
    return d == null || d.add(const Duration(days: 365)).isBefore(now);
  }
}

/// Title + summary of a node in one language (for lists and links).
class NodeHead {
  NodeHead({required this.node, required this.title, required this.summary});
  final NodeRow node;
  final String title;
  final String summary;
}

class Span {
  Span({required this.t, this.v = '', this.b = false, this.i = false, this.w, this.loc, this.to});
  final String t; // text|cite|link|term
  final String v;
  final bool b;
  final bool i;
  final String? w;
  final String? loc;
  final String? to;

  factory Span.fromJson(Map<String, dynamic> j) => Span(
        t: j['t'] as String,
        v: j['v'] as String? ?? '',
        b: j['b'] as bool? ?? false,
        i: j['i'] as bool? ?? false,
        w: j['w'] as String?,
        loc: j['loc'] as String?,
        to: j['to'] as String?,
      );
}

class Block {
  Block({required this.t, this.id, this.text, this.spans = const [], this.items = const []});
  final String t; // h|p|ul|q
  final String? id;
  final String? text;
  final List<Span> spans;
  final List<List<Span>> items;

  factory Block.fromJson(Map<String, dynamic> j) => Block(
        t: j['t'] as String,
        id: j['id'] as String?,
        text: j['text'] as String?,
        spans: [for (final s in (j['s'] as List<dynamic>? ?? const [])) Span.fromJson(s as Map<String, dynamic>)],
        items: [
          for (final it in (j['items'] as List<dynamic>? ?? const []))
            [for (final s in it as List<dynamic>) Span.fromJson(s as Map<String, dynamic>)]
        ],
      );

  static List<Block> listFromJson(String s) =>
      [for (final b in jsonDecode(s) as List<dynamic>) Block.fromJson(b as Map<String, dynamic>)];
}

class NodeDoc {
  NodeDoc({required this.lang, required this.level, required this.title, required this.summary, required this.blocks});
  final String lang;
  final String level;
  final String title;
  final String summary;
  final List<Block> blocks;
}

class Edge {
  Edge({required this.src, required this.rel, required this.dst, this.workKey, this.locator, this.note = const {}, required this.dstExists});
  final String src;
  final String rel;
  final String dst;
  final String? workKey;
  final String? locator;
  final Map<String, String> note;
  final bool dstExists;

  factory Edge.fromRow(Map<String, Object?> r) => Edge(
        src: r['src']! as String,
        rel: r['rel']! as String,
        dst: r['dst']! as String,
        workKey: r['work_key'] as String?,
        locator: r['locator'] as String?,
        note: _i18n(r['note_json']),
        dstExists: (r['dst_exists'] as int? ?? 0) == 1,
      );
}

class Work {
  Work({
    required this.key,
    required this.tier,
    required this.kind,
    required this.title,
    required this.authors,
    required this.year,
    this.venue,
    this.doi,
    this.arxiv,
    this.bibcode,
    this.isbn,
    this.url,
    this.country = const [],
    this.retracted = false,
    this.verified = false,
  });

  final String key;
  final String tier;
  final String kind;
  final String title;
  final List<String> authors;
  final int year;
  final String? venue;
  final String? doi;
  final String? arxiv;
  final String? bibcode;
  final String? isbn;
  final String? url;
  final List<String> country;
  final bool retracted;
  final bool verified;

  bool get peerReviewed => tier == 'T1';
  bool get preprint => tier == 'T2';

  /// Best public link: DOI, arXiv, SciX bibcode, then plain URL.
  Uri? get link {
    if (doi != null) return Uri.https('doi.org', '/$doi');
    if (arxiv != null) return Uri.https('arxiv.org', '/abs/$arxiv');
    if (bibcode != null) return Uri.https('ui.adsabs.harvard.edu', '/abs/$bibcode');
    if (url != null) return Uri.tryParse(url!);
    return null;
  }

  factory Work.fromRow(Map<String, Object?> r) => Work(
        key: r['key']! as String,
        tier: r['tier']! as String,
        kind: r['kind']! as String,
        title: r['title']! as String,
        authors: _strings(r['authors']),
        year: r['year']! as int,
        venue: r['venue'] as String?,
        doi: r['doi'] as String?,
        arxiv: r['arxiv'] as String?,
        bibcode: r['bibcode'] as String?,
        isbn: r['isbn'] as String?,
        url: r['url'] as String?,
        country: _strings(r['country']),
        retracted: (r['retracted'] as int? ?? 0) == 1,
        verified: (r['verified'] as int? ?? 0) == 1,
      );
}

class Prediction {
  Prediction({
    required this.id,
    required this.theory,
    required this.statement,
    this.madeYear,
    this.madeWork,
    this.madeLoc,
    required this.madeBeforeObservation,
    required this.retrodiction,
    this.testYear,
    this.testWork,
    this.testLoc,
    required this.outcome,
    this.note = const {},
    required this.hasCard,
  });

  final String id;
  final String theory;
  final Map<String, String> statement;
  final int? madeYear;
  final String? madeWork;
  final String? madeLoc;
  final bool madeBeforeObservation;
  final bool retrodiction;
  final int? testYear;
  final String? testWork;
  final String? testLoc;
  final String outcome;
  final Map<String, String> note;
  final bool hasCard;

  factory Prediction.fromRow(Map<String, Object?> r) => Prediction(
        id: r['id']! as String,
        theory: r['theory']! as String,
        statement: _i18n(r['statement_json']),
        madeYear: r['made_year'] as int?,
        madeWork: r['made_work'] as String?,
        madeLoc: r['made_loc'] as String?,
        madeBeforeObservation: (r['made_before_observation'] as int? ?? 0) == 1,
        retrodiction: (r['retrodiction'] as int? ?? 0) == 1,
        testYear: r['test_year'] as int?,
        testWork: r['test_work'] as String?,
        testLoc: r['test_loc'] as String?,
        outcome: r['outcome']! as String,
        note: _i18n(r['note_json']),
        hasCard: (r['has_card'] as int? ?? 0) == 1,
      );
}

class ChangelogEntry {
  ChangelogEntry({required this.nodeId, required this.date, required this.kind, required this.note, this.workKey});
  final String nodeId;
  final String date;
  final String kind;
  final String note;
  final String? workKey;

  factory ChangelogEntry.fromRow(Map<String, Object?> r) => ChangelogEntry(
        nodeId: r['node_id']! as String,
        date: r['date']! as String,
        kind: r['kind']! as String,
        note: r['note']! as String,
        workKey: r['work_key'] as String?,
      );
}

class MediaCredit {
  MediaCredit({required this.nodeId, required this.credit, required this.license, required this.sourceUrl});
  final String nodeId;
  final String credit;
  final String license;
  final String sourceUrl;
}
