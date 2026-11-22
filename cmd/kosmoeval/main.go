// Command kosmoeval builds and scores the M2 evaluation set (plan §9).
//
//	kosmoeval candidates   seed sentences citing abstracts → eval/claims.yaml (draft)
//	kosmoeval plant        deliberately wrong variants   → eval/planted.yaml
//	kosmoeval export       items + passages              → out/eval/items.jsonl (for pipeline/verify.py)
//	kosmoeval score        verdicts vs. truth            → gate metrics; exit 1 when the gate fails
//	kosmoeval numbers      genuine sentences whose numbers are not in the cited passages
//
// Abstracts are fetched from Crossref / arXiv (SciX with ADS_API_TOKEN) into
// out/eval/abstracts.json and never committed; items pin them by sha.
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/lioilsources/kosmomat/internal/content"
	"github.com/lioilsources/kosmomat/internal/evalset"
	"github.com/lioilsources/kosmomat/internal/harvest"
)

const (
	claimsFile    = "eval/claims.yaml"
	plantedFile   = "eval/planted.yaml"
	abstractsFile = "out/eval/abstracts.json"
	itemsFile     = "out/eval/items.jsonl"
	verdictsFile  = "out/eval/verdicts.jsonl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: kosmoeval candidates|plant|export|score [flags]")
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]
	var err error
	switch cmd {
	case "candidates":
		err = candidates(args)
	case "plant":
		err = plant(args)
	case "export":
		err = export(args)
	case "score":
		err = score(args)
	case "numbers":
		err = numbers(args)
	default:
		err = fmt.Errorf("unknown command %q", cmd)
	}
	if err != nil {
		slog.Error(cmd, "err", err)
		os.Exit(1)
	}
}

type abstract struct {
	Source string `json:"source"`
	Text   string `json:"text"`
}

func loadAbstracts() map[string]abstract {
	m := map[string]abstract{}
	if b, err := os.ReadFile(abstractsFile); err == nil {
		json.Unmarshal(b, &m)
	}
	return m
}

// fetchAbstracts fills the cache for the given works (Crossref → arXiv → SciX).
func fetchAbstracts(ctx context.Context, c *content.Corpus, keys []string) (map[string]abstract, error) {
	cache := loadAbstracts()
	cr := &harvest.Crossref{C: harvest.NewClient(time.Second)}
	cr.C.Mailto = os.Getenv("CROSSREF_MAILTO")
	ax := &harvest.Arxiv{C: harvest.NewClient(3 * time.Second)}
	var sx *harvest.Scix
	if tok := os.Getenv("ADS_API_TOKEN"); tok != "" {
		cl := harvest.NewClient(time.Second)
		cl.Header = map[string]string{"Authorization": "Bearer " + tok}
		sx = &harvest.Scix{C: cl}
	}
	for _, k := range keys {
		if _, ok := cache[k]; ok {
			continue
		}
		w := c.Works[k]
		if w == nil {
			continue
		}
		var got abstract
		if w.IDs.DOI != "" {
			if r, err := cr.Work(ctx, w.IDs.DOI); err == nil && r.Abstract != "" {
				got = abstract{"crossref", r.Abstract}
			} else if err != nil && !errors.Is(err, harvest.ErrNotFound) {
				return cache, err
			}
		}
		if got.Text == "" && w.IDs.Arxiv != "" {
			m, err := ax.ByIDs(ctx, []string{w.IDs.Arxiv})
			if err != nil {
				return cache, err
			}
			if r, ok := m[harvest.NormArxiv(w.IDs.Arxiv)]; ok && r.Abstract != "" {
				got = abstract{"arxiv", r.Abstract}
			}
		}
		if got.Text == "" && sx != nil && w.IDs.Bibcode != "" {
			if r, err := sx.ByBibcode(ctx, w.IDs.Bibcode); err == nil && r.Abstract != "" {
				got = abstract{"scix", r.Abstract}
			}
		}
		cache[k] = got // empty = looked up, none available
		slog.Info("abstract", "work", k, "source", got.Source, "chars", len(got.Text))
	}
	if err := os.MkdirAll(filepath.Dir(abstractsFile), 0o755); err != nil {
		return cache, err
	}
	b, _ := json.MarshalIndent(cache, "", " ")
	return cache, os.WriteFile(abstractsFile, b, 0o644)
}

func texts(m map[string]abstract) map[string]string {
	out := map[string]string{}
	for k, a := range m {
		out[k] = a.Text
	}
	return out
}

func candidates(args []string) error {
	fs := flag.NewFlagSet("candidates", flag.ExitOnError)
	root := fs.String("content", "content", "content root")
	out := fs.String("out", claimsFile, "output YAML")
	force := fs.Bool("force", false, "overwrite an existing file (loses human verification!)")
	fs.Parse(args)
	if _, err := os.Stat(*out); err == nil && !*force && *out != os.DevNull {
		return fmt.Errorf("%s exists; it may hold human verification — use -force or -out", *out)
	}
	c, err := content.Load(*root)
	if err != nil {
		return err
	}
	segs := evalset.Segments(c, "en")
	need := map[string]bool{}
	for _, s := range segs {
		for _, ct := range s.Cites {
			if ct.Loc == "abstract" {
				need[ct.W] = true
			}
		}
	}
	keys := make([]string, 0, len(need))
	for k := range need {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	abs, err := fetchAbstracts(context.Background(), c, keys)
	if err != nil {
		return err
	}
	items, missing := evalset.Candidates(segs, texts(abs))
	slog.Info("candidates", "segments", len(segs), "items", len(items), "works_without_abstract", len(missing))
	for _, m := range missing {
		fmt.Println("no abstract:", m)
	}
	return writeYAML(*out, items, "# Evaluation set — genuine items (plan §9). Generated by `kosmoeval candidates`.\n"+
		"# Every item is `verified_by: draft` until Ol1n checks that the sentence really follows from the cited\n"+
		"# abstract(s) and sets `expected` (yes|partial|no) and `verified_by: human`. Passages: out/eval/abstracts.json.\n")
}

func plant(args []string) error {
	fs := flag.NewFlagSet("plant", flag.ExitOnError)
	n := fs.Int("n", 20, "number of planted items")
	seed := fs.Uint64("seed", 1, "random seed")
	fs.Parse(args)
	var genuine []evalset.Item
	if err := readYAML(claimsFile, &genuine); err != nil {
		return err
	}
	p := evalset.Plant(genuine, *n, *seed)
	if len(p) < *n {
		slog.Warn("fewer planted items than asked", "got", len(p), "want", *n)
	}
	return writeYAML(plantedFile, p, fmt.Sprintf("# Planted errors (plan §9), `kosmoeval plant -n %d -seed %d`. Expected verdict: no.\n"+
		"# Check each still is wrong (a reversal can accidentally stay true) before relying on it.\n", *n, *seed))
}

type exportItem struct {
	ID       string    `json:"id"`
	Sentence string    `json:"sentence"`
	Passages []passage `json:"passages"`
}

type passage struct {
	Work string `json:"work"`
	Year int    `json:"year"` // publication year of the work; numbers equal to it are metadata, not passage
	Text string `json:"text"`
}

// workYears returns the publication years of the cited works as strings.
func workYears(c *content.Corpus, keys []string) []string {
	var ys []string
	for _, k := range keys {
		if w := c.Works[k]; w != nil {
			ys = append(ys, fmt.Sprint(w.Year))
		}
	}
	return ys
}

func numbers(args []string) error {
	fs := flag.NewFlagSet("numbers", flag.ExitOnError)
	root := fs.String("content", "content", "content root")
	fs.Parse(args)
	c, err := content.Load(*root)
	if err != nil {
		return err
	}
	var items []evalset.Item
	if err := readYAML(claimsFile, &items); err != nil {
		return err
	}
	abs := loadAbstracts()
	bad := 0
	for _, it := range items {
		ps, err := passagesOf(it, abs)
		if err != nil {
			return err
		}
		if ok, miss := evalset.NumbersAgree(it.Sentence, ps, workYears(c, it.Cites)...); !ok {
			bad++
			fmt.Printf("%-40s not in source: %v\n    %s\n", it.ID, miss, it.Sentence)
		}
	}
	fmt.Printf("%d of %d genuine items have numbers not found in the cited abstracts\n", bad, len(items))
	return nil
}

func allItems() ([]evalset.Item, error) {
	var a, b []evalset.Item
	if err := readYAML(claimsFile, &a); err != nil {
		return nil, err
	}
	if err := readYAML(plantedFile, &b); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return append(a, b...), nil
}

func passagesOf(it evalset.Item, abs map[string]abstract) ([]string, error) {
	var ps []string
	for _, k := range it.Cites {
		t := abs[k].Text
		if t == "" {
			return nil, fmt.Errorf("%s: no abstract for %s (run kosmoeval candidates -out /dev/null to refetch)", it.ID, k)
		}
		ps = append(ps, t)
	}
	if got := evalset.SHA(ps); got != it.PassageSHA {
		return nil, fmt.Errorf("%s: passage sha %s ≠ pinned %s (source text changed — re-check the item)", it.ID, got, it.PassageSHA)
	}
	return ps, nil
}

func export(args []string) error {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	out := fs.String("out", itemsFile, "output JSONL")
	root := fs.String("content", "content", "content root")
	fs.Parse(args)
	c, err := content.Load(*root)
	if err != nil {
		return err
	}
	items, err := allItems()
	if err != nil {
		return err
	}
	abs := loadAbstracts()
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return err
	}
	f, err := os.Create(*out)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	// Shuffle-free but planted items are not marked: the verifier only sees id, sentence, passages.
	for _, it := range items {
		ps, err := passagesOf(it, abs)
		if err != nil {
			return err
		}
		e := exportItem{ID: it.ID, Sentence: it.Sentence}
		for i, k := range it.Cites {
			e.Passages = append(e.Passages, passage{Work: k, Year: c.Works[k].Year, Text: ps[i]})
		}
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	slog.Info("exported", "items", len(items), "file", *out)
	return w.Flush()
}

func score(args []string) error {
	fs := flag.NewFlagSet("score", flag.ExitOnError)
	in := fs.String("verdicts", verdictsFile, "verdicts JSONL from pipeline/verify.py")
	root := fs.String("content", "content", "content root")
	fs.Parse(args)
	c, err := content.Load(*root)
	if err != nil {
		return err
	}
	items, err := allItems()
	if err != nil {
		return err
	}
	abs := loadAbstracts()
	truth := map[string]bool{}
	for _, it := range items {
		ps, err := passagesOf(it, abs)
		if err != nil {
			return err
		}
		if t, has := evalset.NumbersTruth(it, ps, workYears(c, it.Cites)...); has {
			truth[it.ID] = t
		}
	}
	f, err := os.Open(*in)
	if err != nil {
		return err
	}
	defer f.Close()
	var vs []evalset.Verdict
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		var v evalset.Verdict
		if err := json.Unmarshal(sc.Bytes(), &v); err != nil {
			return fmt.Errorf("%s: %w", *in, err)
		}
		vs = append(vs, v)
	}
	s := evalset.Evaluate(items, vs, truth)
	fmt.Println(s)
	for _, m := range s.Misses {
		fmt.Println("  miss:", m)
	}
	drafts := 0
	for _, it := range items {
		if it.VerifiedBy != "human" {
			drafts++
		}
	}
	if drafts > 0 {
		fmt.Printf("note: %d items are not human-verified yet — the score is provisional\n", drafts)
	}
	if !s.Pass() {
		return errors.New("gate: no-go")
	}
	fmt.Println("gate: go")
	return nil
}

func readYAML(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	return dec.Decode(v)
}

func writeYAML(path string, v any, header string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, append([]byte(header), b...), 0o644)
}
