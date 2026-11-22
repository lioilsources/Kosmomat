// Package pack turns a corpus into the app's SQLite content pack,
// a manifest and a JSON export for the Astro web.
package pack

import (
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/klauspost/compress/zstd"
	_ "modernc.org/sqlite"

	"github.com/lioilsources/kosmomat/internal/content"
)

//go:embed schema.sql
var schemaSQL string

// Options configure a build.
type Options struct {
	OutDir  string
	Version string // pack version; empty = "dev-<content hash>"
	MinApp  int    // minimum app build that can read this pack
	BaseURL string // public base URL where the pack is served (Cloudflare)
	// Editorial maps a language to the Markdown of docs/EDITORIAL.md, shipped in meta.
	Editorial map[string]string
}

// Manifest is published next to the pack; the app polls it for updates.
type Manifest struct {
	Version   string `json:"version"`
	SHA256    string `json:"sha256"`
	Size      int64  `json:"size"`
	MinApp    int    `json:"min_app"`
	URL       string `json:"url"`
	Nodes     int    `json:"nodes"`
	BuiltFrom string `json:"built_from"` // sha256 of the content tree
}

// Build writes pack.sqlite, pack.sqlite.zst, manifest.json and web/ into opt.OutDir.
func Build(c *content.Corpus, opt Options) (*Manifest, error) {
	if err := os.MkdirAll(opt.OutDir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir out: %w", err)
	}
	hash, err := ContentHash(c.Root)
	if err != nil {
		return nil, err
	}
	if opt.Version == "" {
		opt.Version = "dev-" + hash[:12]
	}
	dbPath := filepath.Join(opt.OutDir, "pack.sqlite")
	_ = os.Remove(dbPath)
	if err := writeDB(c, dbPath, opt.Version, hash, opt.Editorial); err != nil {
		return nil, err
	}
	zPath := dbPath + ".zst"
	sum, size, err := compress(dbPath, zPath)
	if err != nil {
		return nil, err
	}
	m := &Manifest{
		Version: opt.Version, SHA256: sum, Size: size, MinApp: opt.MinApp,
		URL:   strings.TrimRight(opt.BaseURL, "/") + "/pack-" + opt.Version + ".sqlite.zst",
		Nodes: len(c.Nodes), BuiltFrom: hash,
	}
	if err := writeJSON(filepath.Join(opt.OutDir, "manifest.json"), m); err != nil {
		return nil, err
	}
	if err := exportWeb(c, filepath.Join(opt.OutDir, "web"), m); err != nil {
		return nil, err
	}
	return m, nil
}

func js(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func writeDB(c *content.Corpus, path, version, hash string, editorial map[string]string) (err error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("open pack: %w", err)
	}
	defer func() {
		if cerr := db.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()
	if _, err := db.Exec(schemaSQL); err != nil {
		return fmt.Errorf("schema: %w", err)
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	exec := func(q string, args ...any) {
		if err != nil {
			return
		}
		if _, e := tx.Exec(q, args...); e != nil {
			err = fmt.Errorf("%s: %w", strings.Fields(q)[2], e)
		}
	}

	exec(`INSERT INTO meta VALUES ('version', ?), ('built_from', ?), ('schema', '1')`, version, hash)
	for _, lang := range sortedKeys(editorial) {
		exec(`INSERT INTO meta VALUES (?, ?)`, "editorial_"+lang, editorial[lang])
	}

	for _, k := range c.WorkKeys() {
		w := c.Works[k]
		exec(`INSERT INTO works VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			w.Key, w.Tier, w.Kind, w.Title, js(w.Authors), w.Year, w.Venue,
			nullable(w.IDs.DOI), nullable(w.IDs.Arxiv), nullable(w.IDs.Bibcode), nullable(w.IDs.ISBN), nullable(w.IDs.URL),
			js(w.Country), b2i(w.Retracted), b2i(w.Verified))
	}

	for _, id := range c.NodeIDs() {
		n := c.Nodes[id]
		var fy any
		var fby, fwork, floc any
		if fp := n.FirstProposed; fp != nil {
			fy, fby, fwork, floc = fp.Year, js(fp.By), nullable(content.WorkKey(fp.In)), nullable(fp.Loc)
		}
		extra := map[string]any{}
		if n.Difficulty != "" {
			extra["difficulty"] = n.Difficulty
		}
		if len(n.Needs) > 0 {
			extra["needs"] = n.Needs
		}
		if len(n.TodoVerify) > 0 {
			extra["todo_verify"] = len(n.TodoVerify)
		}
		exec(`INSERT INTO nodes VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			n.ID, n.Type, nullable(n.Status), nullable(n.Literature), nullable(n.Testability), nullable(n.CaseKind),
			nullable(n.Era), b2i(n.Astromat), fy, fby, fwork, floc, n.LastVerified, n.VerifiedBy, js(extra))

		sources := map[[2]string]bool{}
		addSource := func(work, loc string) {
			if work == "" {
				return
			}
			sources[[2]string{content.WorkKey(work), loc}] = true
		}
		if n.FirstProposed != nil {
			addSource(n.FirstProposed.In, n.FirstProposed.Loc)
		}
		for _, key := range sortedKeys(n.Docs) {
			d := n.Docs[key]
			blocks, _ := json.Marshal(d.Blocks)
			exec(`INSERT INTO node_i18n VALUES (?,?,?,?,?,?)`, n.ID, d.Lang, d.Level, n.Title[d.Lang], n.Summary[d.Lang], string(blocks))
			var body []string
			for _, b := range d.Blocks {
				if b.T == "h" {
					body = append(body, b.Text)
				}
				for _, p := range b.Paragraphs() {
					body = append(body, content.PlainText(p, func(id string) string { return c.Title(id, d.Lang) }))
					for _, ct := range content.Cites(p) {
						addSource(ct.W, ct.Loc)
					}
				}
			}
			if d.Level == "lay" {
				exec(`INSERT INTO search VALUES (?,?,?,?,?)`, n.ID, d.Lang, n.Title[d.Lang], n.Summary[d.Lang], strings.Join(body, "\n"))
			}
		}
		for _, e := range n.Edges {
			var wk, loc any
			if e.Source != nil {
				wk, loc = nullable(content.WorkKey(e.Source.Work)), nullable(e.Source.Loc)
				addSource(e.Source.Work, e.Source.Loc)
			}
			_, exists := c.Nodes[e.To]
			var note any
			if len(e.Note) > 0 {
				note = js(e.Note)
			}
			exec(`INSERT INTO edges VALUES (?,?,?,?,?,?,?)`, n.ID, e.Rel, e.To, wk, loc, note, b2i(exists))
		}
		srcKeys := make([][2]string, 0, len(sources))
		for k := range sources {
			if _, ok := c.Works[k[0]]; ok {
				srcKeys = append(srcKeys, k)
			}
		}
		sort.Slice(srcKeys, func(i, j int) bool {
			if srcKeys[i][0] != srcKeys[j][0] {
				return srcKeys[i][0] < srcKeys[j][0]
			}
			return srcKeys[i][1] < srcKeys[j][1]
		})
		for _, k := range srcKeys {
			exec(`INSERT INTO node_sources VALUES (?,?,?)`, n.ID, k[0], k[1])
		}
		for _, ch := range n.Changelog {
			exec(`INSERT INTO changelog VALUES (?,?,?,?,?)`, n.ID, ch.Date, ch.Kind, ch.Note, nullable(content.WorkKey(ch.Work)))
		}
	}

	for _, p := range c.AllPredictions() {
		var my, mby, mw, ml any
		if p.Made != nil {
			my, mby, mw, ml = p.Made.Year, js(p.Made.By), nullable(content.WorkKey(p.Made.In)), nullable(p.Made.Loc)
		}
		var ty, tobs, tby, tw, tl any
		if p.Test != nil {
			ty, tobs, tby, tw, tl = p.Test.Year, nullable(p.Test.Obs), js(p.Test.By), nullable(content.WorkKey(p.Test.In)), nullable(p.Test.Loc)
		}
		var note any
		if len(p.Note) > 0 {
			note = js(p.Note)
		}
		exec(`INSERT INTO predictions VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			p.ID, p.Theory, js(p.Statement), my, mby, mw, ml, b2i(p.MadeBeforeObservation), b2i(p.Retrodiction),
			ty, tobs, tby, tw, tl, p.Outcome, note, b2i(!p.Inline))
	}

	for _, m := range c.Media {
		exec(`INSERT INTO media VALUES (?,?,?,?,?,?,?,?,?)`, m.SHA256[:min(16, len(m.SHA256))], m.Node, m.SHA256, m.Ext, m.Credit, m.License, m.SourceURL, m.W, m.H)
	}
	for _, p := range c.Paths {
		exec(`INSERT INTO paths VALUES (?,?,?)`, p.Slug, js(p.Title), js(p.Nodes))
	}
	for _, old := range sortedKeys(c.Redirects) {
		exec(`INSERT INTO redirects VALUES (?,?)`, old, c.Redirects[old])
	}
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	if _, err = db.Exec(`INSERT INTO search(search) VALUES ('optimize'); VACUUM;`); err != nil {
		return fmt.Errorf("optimize: %w", err)
	}
	return nil
}

func compress(src, dst string) (sum string, size int64, err error) {
	in, err := os.Open(src)
	if err != nil {
		return "", 0, err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return "", 0, err
	}
	defer out.Close()
	h := sha256.New()
	enc, err := zstd.NewWriter(io.MultiWriter(out, h), zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return "", 0, err
	}
	if _, err := io.Copy(enc, in); err != nil {
		return "", 0, fmt.Errorf("zstd: %w", err)
	}
	if err := enc.Close(); err != nil {
		return "", 0, fmt.Errorf("zstd: %w", err)
	}
	st, err := out.Stat()
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), st.Size(), nil
}

// ContentHash hashes every file under root (path + bytes), in sorted order.
func ContentHash(root string) (string, error) {
	var files []string
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && !strings.HasPrefix(d.Name(), ".") {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("hash content: %w", err)
	}
	sort.Strings(files)
	h := sha256.New()
	for _, f := range files {
		rel, _ := filepath.Rel(root, f)
		b, err := os.ReadFile(f)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(h, "%s\x00%d\x00", filepath.ToSlash(rel), len(b))
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func sortedKeys[V any](m map[string]V) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// CopyToApp copies the uncompressed pack and its manifest into the app's
// asset folder, where it ships as the bundled pack v1 (plan §6).
func CopyToApp(outDir, assetDir string) error {
	if err := os.MkdirAll(assetDir, 0o755); err != nil {
		return fmt.Errorf("mkdir app assets: %w", err)
	}
	for _, name := range []string{"pack.sqlite", "manifest.json"} {
		b, err := os.ReadFile(filepath.Join(outDir, name))
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		if err := os.WriteFile(filepath.Join(assetDir, name), b, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	return nil
}
