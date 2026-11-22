-- works_raw / harvest_runs: same tables the Postgres production line will use (plan §7).
CREATE TABLE IF NOT EXISTS works_raw (
  id             INTEGER PRIMARY KEY,
  doi            TEXT UNIQUE,
  arxiv          TEXT UNIQUE,
  bibcode        TEXT UNIQUE,
  title          TEXT NOT NULL DEFAULT '',
  authors        TEXT NOT NULL DEFAULT '[]', -- JSON [{family, given}]
  year           INTEGER NOT NULL DEFAULT 0,
  venue          TEXT NOT NULL DEFAULT '',
  categories     TEXT NOT NULL DEFAULT '',   -- space separated arXiv categories
  abstract       TEXT NOT NULL DEFAULT '',
  license        TEXT NOT NULL DEFAULT '',
  retracted      INTEGER NOT NULL DEFAULT 0,
  source         TEXT NOT NULL,              -- first source that reported it
  source_updated TEXT NOT NULL DEFAULT '',
  first_seen     TEXT NOT NULL,
  last_seen      TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS works_raw_first_seen ON works_raw(first_seen);

CREATE TABLE IF NOT EXISTS harvest_runs (
  id       INTEGER PRIMARY KEY,
  source   TEXT NOT NULL,
  params   TEXT NOT NULL,
  started  TEXT NOT NULL,
  finished TEXT,
  fetched  INTEGER NOT NULL DEFAULT 0,
  inserted INTEGER NOT NULL DEFAULT 0,
  error    TEXT NOT NULL DEFAULT ''
);
