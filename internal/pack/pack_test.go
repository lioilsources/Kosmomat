package pack

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/lioilsources/kosmomat/internal/content"
)

// TestBuildRealContent builds the repository's own content and queries the pack.
func TestBuildRealContent(t *testing.T) {
	c, err := content.Load("../../content")
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	m, err := Build(c, Options{OutDir: out, Version: "test", MinApp: 1, BaseURL: "https://x/packs"})
	if err != nil {
		t.Fatal(err)
	}
	if m.Nodes != len(c.Nodes) || m.Size == 0 || len(m.SHA256) != 64 {
		t.Fatalf("manifest: %+v", m)
	}
	db, err := sql.Open("sqlite", filepath.Join(out, "pack.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var n int
	if err := db.QueryRow(`SELECT count(*) FROM nodes`).Scan(&n); err != nil || n != len(c.Nodes) {
		t.Fatalf("nodes = %d, %v", n, err)
	}
	if err := db.QueryRow(`SELECT count(*) FROM node_i18n WHERE level='lay'`).Scan(&n); err != nil || n != 2*len(c.Nodes) {
		t.Fatalf("lay docs = %d, %v", n, err)
	}
	// FTS5 with diacritics removed: "neutrin" matches Czech text.
	if err := db.QueryRow(`SELECT count(*) FROM search WHERE search MATCH 'neutrina' AND lang='cs'`).Scan(&n); err != nil || n == 0 {
		t.Fatalf("fts neutrina = %d, %v", n, err)
	}
	if err := db.QueryRow(`SELECT count(*) FROM search WHERE search MATCH 'hmota' AND lang='cs'`).Scan(&n); err != nil || n == 0 {
		t.Fatalf("fts hmota = %d, %v", n, err)
	}
	// Every citation in the pack resolves to a work row.
	if err := db.QueryRow(`SELECT count(*) FROM node_sources s LEFT JOIN works w ON w.key=s.work_key WHERE w.key IS NULL`).Scan(&n); err != nil || n != 0 {
		t.Fatalf("dangling sources = %d, %v", n, err)
	}
	if err := db.QueryRow(`SELECT count(*) FROM predictions WHERE theory='theory.electric-universe'`).Scan(&n); err != nil || n < 3 {
		t.Fatalf("EU predictions = %d, %v", n, err)
	}
	var blocks string
	if err := db.QueryRow(`SELECT blocks_json FROM node_i18n WHERE node_id='theory.lcdm' AND lang='cs' AND level='lay'`).Scan(&blocks); err != nil {
		t.Fatal(err)
	}
	var parsed []content.Block
	if err := json.Unmarshal([]byte(blocks), &parsed); err != nil || len(parsed) == 0 {
		t.Fatalf("blocks json: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(out, "web", "n", "body.mars.json"))
	if err != nil {
		t.Fatal(err)
	}
	var w WebNode
	if err := json.Unmarshal(b, &w); err != nil || !w.Astromat || len(w.Works) == 0 {
		t.Fatalf("web export mars: astromat=%v works=%d err=%v", w.Astromat, len(w.Works), err)
	}
}

func TestBuildDeterministicHash(t *testing.T) {
	a, err := ContentHash("../../content")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := ContentHash("../../content")
	if a != b {
		t.Fatal("content hash not stable")
	}
}
