// Read-only access to the bundled SQLite content pack.
import 'dart:io';

import 'package:sqlite3/sqlite3.dart';

import 'models.dart';

/// Relations that point from a node to its ancestor (the idea it came from).
const ancestorRels = {'builds_on', 'inspired_by', 'interprets', 'first_proposed_in'};

/// Relations shown under "who disagrees".
const opposingRels = {'opposed_by', 'contradicted_by'};

class Pack {
  Pack._(this._db);

  final Database _db;

  /// Opens a pack file read-only.
  factory Pack.open(String path) => Pack._(sqlite3.open(path, mode: OpenMode.readOnly));

  /// Opens a pack from bytes in memory (tests, or before the file copy finishes).
  factory Pack.fromBytes(List<int> bytes) {
    final dir = Directory.systemTemp.createTempSync('kosmomat');
    final f = File('${dir.path}/pack.sqlite')..writeAsBytesSync(bytes);
    return Pack.open(f.path);
  }

  void dispose() => _db.close();

  String? meta(String key) {
    final r = _db.select('SELECT value FROM meta WHERE key = ?', [key]);
    return r.isEmpty ? null : r.first['value'] as String;
  }

  String get version => _db.select("SELECT value FROM meta WHERE key='version'").first['value'] as String;

  NodeRow? node(String id) {
    final r = _db.select('SELECT * FROM nodes WHERE id = ?', [id]);
    return r.isEmpty ? null : NodeRow.fromRow(r.first);
  }

  /// All reading levels available for a node in [lang].
  List<String> levels(String id, String lang) => [
        for (final r in _db.select(
            "SELECT level FROM node_i18n WHERE node_id = ? AND lang = ? ORDER BY CASE level WHEN 'kid' THEN 0 WHEN 'lay' THEN 1 ELSE 2 END",
            [id, lang]))
          r['level'] as String
      ];

  NodeDoc? doc(String id, String lang, String level) {
    var r = _db.select('SELECT * FROM node_i18n WHERE node_id = ? AND lang = ? AND level = ?', [id, lang, level]);
    if (r.isEmpty) {
      r = _db.select("SELECT * FROM node_i18n WHERE node_id = ? AND lang = ? AND level = 'lay'", [id, lang]);
    }
    if (r.isEmpty) return null;
    final row = r.first;
    return NodeDoc(
      lang: row['lang'] as String,
      level: row['level'] as String,
      title: row['title'] as String,
      summary: row['summary'] as String,
      blocks: Block.listFromJson(row['blocks_json'] as String),
    );
  }

  /// Heads (title + summary) of nodes, optionally filtered by type.
  List<NodeHead> heads(String lang, {String? type, Iterable<String>? ids}) {
    final where = <String>["i.lang = ?", "i.level = 'lay'"];
    final args = <Object?>[lang];
    if (type != null) {
      where.add('n.type = ?');
      args.add(type);
    }
    if (ids != null) {
      final list = ids.toList();
      if (list.isEmpty) return const [];
      where.add('n.id IN (${List.filled(list.length, '?').join(',')})');
      args.addAll(list);
    }
    final rows = _db.select(
        'SELECT n.*, i.title, i.summary FROM nodes n JOIN node_i18n i ON i.node_id = n.id WHERE ${where.join(' AND ')} ORDER BY i.title',
        args);
    return [for (final r in rows) NodeHead(node: NodeRow.fromRow(r), title: r['title'] as String, summary: r['summary'] as String)];
  }

  NodeHead? head(String id, String lang) => heads(lang, ids: [id]).firstOrNull;

  String title(String id, String lang) => head(id, lang)?.title ?? id;

  List<Edge> edgesFrom(String id) => [for (final r in _db.select('SELECT * FROM edges WHERE src = ?', [id])) Edge.fromRow(r)];

  List<Edge> edgesTo(String id) => [for (final r in _db.select('SELECT * FROM edges WHERE dst = ?', [id])) Edge.fromRow(r)];

  Work? work(String key) {
    final r = _db.select('SELECT * FROM works WHERE key = ?', [key]);
    return r.isEmpty ? null : Work.fromRow(r.first);
  }

  /// Works cited by a node, with all locators used on the card.
  List<(Work, List<String>)> sources(String id) {
    final rows = _db.select(
        'SELECT w.*, group_concat(s.locator, char(31)) AS locs FROM node_sources s JOIN works w ON w.key = s.work_key '
        'WHERE s.node_id = ? GROUP BY w.key ORDER BY w.year, w.key',
        [id]);
    return [
      for (final r in rows) (Work.fromRow(r), (r['locs'] as String? ?? '').split('\u001f').where((s) => s.isNotEmpty).toList())
    ];
  }

  List<Prediction> predictions(String theory) => [
        for (final r in _db.select('SELECT * FROM predictions WHERE theory = ? OR id = ? ORDER BY made_year, id', [theory, theory]))
          Prediction.fromRow(r)
      ];

  List<ChangelogEntry> changelog({String? nodeId, int limit = 50}) => [
        for (final r in nodeId == null
            ? _db.select('SELECT * FROM changelog ORDER BY date DESC, node_id LIMIT ?', [limit])
            : _db.select('SELECT * FROM changelog WHERE node_id = ? ORDER BY date DESC', [nodeId]))
          ChangelogEntry.fromRow(r)
      ];

  List<MediaCredit> mediaCredits() => [
        for (final r in _db.select('SELECT node_id, credit, license, source_url FROM media ORDER BY node_id'))
          MediaCredit(
            nodeId: r['node_id'] as String,
            credit: r['credit'] as String,
            license: r['license'] as String,
            sourceUrl: r['source_url'] as String,
          )
      ];

  String? redirect(String id) {
    final r = _db.select('SELECT new_id FROM redirects WHERE old_id = ?', [id]);
    return r.isEmpty ? null : r.first['new_id'] as String;
  }

  /// Offline full-text search over title, summary and body (lay level).
  List<NodeHead> search(String query, String lang, {String? type, String? status}) {
    final match = ftsQuery(query);
    if (match.isEmpty) return const [];
    final where = <String>['search MATCH ?', 's.lang = ?'];
    final args = <Object?>[match, lang];
    if (type != null) {
      where.add('n.type = ?');
      args.add(type);
    }
    if (status != null) {
      where.add('n.status = ?');
      args.add(status);
    }
    final rows = _db.select(
        'SELECT n.*, i.title, i.summary FROM search s JOIN nodes n ON n.id = s.node_id '
        "JOIN node_i18n i ON i.node_id = n.id AND i.lang = ? AND i.level = 'lay' "
        'WHERE ${where.join(' AND ')} ORDER BY bm25(search, 10.0, 4.0, 1.0)',
        [lang, ...args]);
    return [for (final r in rows) NodeHead(node: NodeRow.fromRow(r), title: r['title'] as String, summary: r['summary'] as String)];
  }

  /// Turns free user input into a safe FTS5 prefix query ("tok"* "tok2"*).
  static String ftsQuery(String input) {
    final tokens = input
        .split(RegExp(r'[^\p{L}\p{N}]+', unicode: true))
        .where((t) => t.isNotEmpty)
        .map((t) => '"${t.replaceAll('"', '')}"*');
    return tokens.join(' ');
  }

  /// Lineage of a node: ancestors up to the origin, descendants and competitors.
  Lineage lineage(String id) {
    final ancestors = <LineageItem>[];
    final seen = <String>{id};
    var frontier = [id];
    for (var depth = 1; depth <= 6 && frontier.isNotEmpty; depth++) {
      final next = <String>[];
      for (final src in frontier) {
        for (final e in edgesFrom(src)) {
          if (!ancestorRels.contains(e.rel) || !seen.add(e.dst)) continue;
          final n = node(e.dst);
          ancestors.add(LineageItem(id: e.dst, rel: e.rel, depth: depth, year: n?.firstYear, exists: n != null, via: src));
          if (n != null && n.firstWork == null) next.add(e.dst);
        }
      }
      frontier = next;
    }
    final descendants = <LineageItem>[
      for (final e in edgesTo(id))
        if (ancestorRels.contains(e.rel)) LineageItem(id: e.src, rel: e.rel, depth: 1, year: node(e.src)?.firstYear, exists: true, via: id),
      for (final e in edgesFrom(id))
        if (e.rel == 'superseded_by')
          LineageItem(id: e.dst, rel: e.rel, depth: 1, year: node(e.dst)?.firstYear, exists: e.dstExists, via: id),
    ];
    final competitors = <String>{
      for (final e in edgesFrom(id))
        if (e.rel == 'competes_with') e.dst,
      for (final e in edgesTo(id))
        if (e.rel == 'competes_with') e.src,
    };
    final origins = <Origin>[];
    for (final nid in [id, ...ancestors.where((a) => a.exists).map((a) => a.id)]) {
      final n = node(nid);
      if (n?.firstWork != null) {
        origins.add(Origin(nodeId: nid, year: n!.firstYear, workKey: n.firstWork!, locator: n.firstLoc, by: n.firstBy));
      }
    }
    ancestors.sort((a, b) => (a.year ?? 9999).compareTo(b.year ?? 9999));
    return Lineage(
      ancestors: ancestors,
      descendants: descendants,
      competitors: [for (final c in competitors) LineageItem(id: c, rel: 'competes_with', depth: 0, year: node(c)?.firstYear, exists: node(c) != null, via: id)],
      origins: origins,
    );
  }
}

class LineageItem {
  LineageItem({required this.id, required this.rel, required this.depth, this.year, required this.exists, required this.via});
  final String id;
  final String rel;
  final int depth;
  final int? year;
  final bool exists;
  final String via;
}

class Origin {
  Origin({required this.nodeId, this.year, required this.workKey, this.locator, this.by = const []});
  final String nodeId;
  final int? year;
  final String workKey;
  final String? locator;
  final List<String> by;
}

class Lineage {
  Lineage({required this.ancestors, required this.descendants, required this.competitors, required this.origins});
  final List<LineageItem> ancestors;
  final List<LineageItem> descendants;
  final List<LineageItem> competitors;
  final List<Origin> origins;

  /// Acceptance M1: every claim/theory reaches a node with first_proposed.
  bool get reachesOrigin => origins.isNotEmpty;
}
