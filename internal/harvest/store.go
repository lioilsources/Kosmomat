package harvest

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// Store is the local works_raw database. Plan §7 puts the production line in
// Postgres on JODA; the M2 spike keeps the same tables in a SQLite file so the
// harvester runs anywhere without infrastructure.
type Store struct{ db *sql.DB }

//go:embed store.sql
var storeSchema string

// OpenStore opens (and migrates) the SQLite database at path.
func OpenStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(storeSchema); err != nil {
		db.Close()
		return nil, fmt.Errorf("store schema: %w", err)
	}
	return &Store{db: db}, nil
}

// Close closes the database.
func (s *Store) Close() error { return s.db.Close() }

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// Upsert stores r, deduplicating on DOI, arXiv ID or bibcode. Fields already
// known are kept unless the new record fills an empty one. Returns true when
// the record was new.
func (s *Store) Upsert(r Record, now string) (bool, error) {
	r.DOI, r.Arxiv = NormDOI(r.DOI), NormArxiv(r.Arxiv)
	var id int64
	err := s.db.QueryRow(`SELECT id FROM works_raw WHERE doi = ? OR arxiv = ? OR bibcode = ? LIMIT 1`,
		nullable(r.DOI), nullable(r.Arxiv), nullable(r.Bibcode)).Scan(&id)
	authors, _ := json.Marshal(r.Authors)
	cats := strings.Join(r.Categories, " ")
	switch {
	case err == sql.ErrNoRows:
		_, err = s.db.Exec(`INSERT INTO works_raw
			(doi, arxiv, bibcode, title, authors, year, venue, categories, abstract, license, retracted, source, source_updated, first_seen, last_seen)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			nullable(r.DOI), nullable(r.Arxiv), nullable(r.Bibcode), r.Title, string(authors), r.Year(), r.Venue, cats,
			r.Abstract, r.License, r.Retracted, r.Source, r.Updated, now, now)
		return err == nil, err
	case err != nil:
		return false, err
	}
	_, err = s.db.Exec(`UPDATE works_raw SET
			title = CASE WHEN title = '' THEN ? ELSE title END,
			authors = CASE WHEN authors = 'null' OR authors = '[]' THEN ? ELSE authors END,
			year = CASE WHEN year = 0 THEN ? ELSE year END,
			venue = CASE WHEN venue = '' THEN ? ELSE venue END,
			categories = CASE WHEN categories = '' THEN ? ELSE categories END,
			abstract = CASE WHEN abstract = '' THEN ? ELSE abstract END,
			license = CASE WHEN license = '' THEN ? ELSE license END,
			retracted = retracted OR ?,
			source_updated = ?,
			last_seen = ?
		WHERE id = ?`,
		r.Title, string(authors), r.Year(), r.Venue, cats, r.Abstract, r.License, r.Retracted, r.Updated, now, id)
	if err != nil {
		return false, err
	}
	// Fill identifiers the row lacks; a clash with another row means two rows
	// describe the same work — leave both and let triage see it.
	for _, col := range []struct{ name, val string }{{"doi", r.DOI}, {"arxiv", r.Arxiv}, {"bibcode", r.Bibcode}} {
		if col.val == "" {
			continue
		}
		_, _ = s.db.Exec(`UPDATE works_raw SET `+col.name+` = ? WHERE id = ? AND `+col.name+` IS NULL`, col.val, id)
	}
	return false, nil
}

// Run is one harvest_runs row.
type Run struct {
	ID       int64
	Source   string
	Params   string
	Fetched  int
	Inserted int
}

// StartRun records the start of a harvest.
func (s *Store) StartRun(source, params, now string) (*Run, error) {
	res, err := s.db.Exec(`INSERT INTO harvest_runs (source, params, started) VALUES (?,?,?)`, source, params, now)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Run{ID: id, Source: source, Params: params}, nil
}

// FinishRun records counts and the error, if any.
func (s *Store) FinishRun(r *Run, now string, runErr error) error {
	msg := ""
	if runErr != nil {
		msg = runErr.Error()
	}
	_, err := s.db.Exec(`UPDATE harvest_runs SET finished = ?, fetched = ?, inserted = ?, error = ? WHERE id = ?`,
		now, r.Fetched, r.Inserted, msg, r.ID)
	return err
}

// Stats returns record counts per source and the total.
func (s *Store) Stats() (map[string]int, error) {
	rows, err := s.db.Query(`SELECT source, count(*) FROM works_raw GROUP BY source`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := map[string]int{}
	for rows.Next() {
		var src string
		var n int
		if err := rows.Scan(&src, &n); err != nil {
			return nil, err
		}
		m[src] = n
		m["total"] += n
	}
	return m, rows.Err()
}
