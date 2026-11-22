package harvest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/lioilsources/kosmomat/internal/content"
)

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// server replays the recorded responses; anything else is 404.
func server(t *testing.T) (*httptest.Server, *[]string) {
	var seen []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.URL.RequestURI())
		if ua := r.Header.Get("User-Agent"); !strings.HasPrefix(ua, "metadata-harvester/") {
			t.Errorf("user agent %q", ua)
		}
		q := r.URL.Query()
		switch {
		case strings.EqualFold(r.URL.Path, "/works/10.1103/PhysRev.43.491"):
			w.Write(fixture(t, "crossref_work.json"))
		case r.URL.Path == "/api" && q.Get("id_list") == "1112.3960":
			w.Write(fixture(t, "arxiv_atom.xml"))
		case r.URL.Path == "/oai" && q.Get("resumptionToken") == "tok|1":
			w.Write(fixture(t, "oai_page2.xml"))
		case r.URL.Path == "/oai" && q.Get("set") == "physics:gr-qc":
			w.Write(fixture(t, "oai_page1.xml"))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &seen
}

func testClient() *Client {
	c := NewClient(0)
	c.Sleep = func(time.Duration) {}
	return c
}

func TestCrossrefWork(t *testing.T) {
	srv, _ := server(t)
	cr := &Crossref{C: testClient(), Base: srv.URL}
	r, err := cr.Work(context.Background(), "https://doi.org/10.1103/PhysRev.43.491")
	if err != nil {
		t.Fatal(err)
	}
	if r.Title != "The Positive Electron" || r.Year() != 1933 || r.Authors[0].Family != "Anderson" || r.Venue != "Physical Review" {
		t.Fatalf("record %+v", r)
	}
	if _, err := cr.Work(context.Background(), "10.0/none"); err != ErrNotFound {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestArxivByIDs(t *testing.T) {
	srv, _ := server(t)
	ax := &Arxiv{C: testClient(), APIBase: srv.URL + "/api"}
	m, err := ax.ByIDs(context.Background(), []string{"arXiv:1112.3960v1"})
	if err != nil {
		t.Fatal(err)
	}
	r, ok := m["1112.3960"]
	if !ok {
		t.Fatalf("missing record: %v", m)
	}
	if r.Authors[0].Family != "Famaey" || r.Authors[1].Family != "McGaugh" || r.Year() != 2011 || r.DOI != "10.12942/lrr-2012-10" {
		t.Fatalf("record %+v", r)
	}
	if !slices.Contains(r.Categories, "astro-ph.CO") {
		t.Fatalf("categories %v", r.Categories)
	}
}

func TestArxivHarvestPagesAndFilters(t *testing.T) {
	srv, seen := server(t)
	ax := &Arxiv{C: testClient(), OAIBase: srv.URL + "/oai"}
	var ids []string
	err := ax.Harvest(context.Background(), "physics:gr-qc", "2026-11-19", "2026-11-19", DefaultCategories, func(r Record) error {
		if r.Title == "" || len(r.Authors) == 0 || r.Abstract == "" {
			t.Errorf("incomplete record %+v", r)
		}
		ids = append(ids, r.Arxiv)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// 2505.24743 is in page 2 but re-labelled hep-th only in the fixture → filtered out.
	want := []string{"2407.02076", "2407.04743", "2508.05536"}
	if !slices.Equal(ids, want) {
		t.Fatalf("ids %v, want %v", ids, want)
	}
	if len(*seen) != 2 {
		t.Fatalf("requests %v", *seen)
	}
}

func TestFoldAndTitles(t *testing.T) {
	if got := Fold(`The $\Lambda$CDM model: Café <i>test</i>`); got != "the lambda cdm model cafe test" {
		t.Fatalf("fold %q", got)
	}
	if titlesAgree("Planck 2018 results. VI. Cosmological parameters", "Planck 2018 results") {
		t.Fatal("short prefix (< half) must not agree")
	}
	if !titlesAgree("Die Rotverschiebung von extragalaktischen Nebeln", "Die Rotverschiebung von extragalaktischen Nebeln.") {
		t.Fatal("punctuation must not matter")
	}
	for a, b := range map[string]string{
		"Measurements of Omega and Lambda from 42 High-Redshift Supernovae": "Measurements of Ω and Λ from 42 High‐Redshift Supernovae",
		"with 1 km/s/Mpc Uncertainty":                                       "with 1 km s −1 Mpc −1 Uncertainty",
		"Planck 2018 results. VI. Cosmological parameters":                  "Planck 2018 results: VI. Cosmological parameters",
		"Astronomy & Astrophysics":                                          "Astronomy &amp; Astrophysics",
	} {
		if !titlesAgree(a, b) {
			t.Errorf("%q vs %q should agree (%q / %q)", a, b, Fold(a), Fold(b))
		}
	}
	if titlesAgree("Measurements of Omega", "Measurements of Lambda") {
		t.Fatal("different Greek letters must differ")
	}
	if !titlesAgree("The Positive Electron", "The positive electron: a study") {
		t.Fatal("subtitle should agree")
	}
}

func TestSplitName(t *testing.T) {
	for in, fam := range map[string]string{"Benoit Famaey": "Famaey", "Floor van Leeuwen": "van Leeuwen", "Planck": "Planck"} {
		if got := splitName(in).Family; got != fam {
			t.Errorf("%q → %q, want %q", in, got, fam)
		}
	}
}

func writeCorpus(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"works/1933/anderson-1933.yaml": `key: anderson-1933
tier: T1
kind: article
title: "The Positive Electron"
authors: ["Anderson, C. D."]
year: 1933
venue: "Physical Review 43, 491–494"
ids: { doi: "10.1103/PhysRev.43.491", arxiv: "", bibcode: "", isbn: "", url: "" }
retracted: false
verified: false
`,
		// arXiv-only with the journal year: the arXiv record says 2011 → mismatch.
		"works/2012/famaey-mcgaugh-2012.yaml": `key: famaey-mcgaugh-2012
tier: T1
kind: article
title: "Modified Newtonian Dynamics (MOND): Observational Phenomenology and Relativistic Extensions"
authors: ["Famaey, B.", "McGaugh, S. S."]
year: 2012
venue: "Living Reviews in Relativity 15, 10"
ids: { doi: "", arxiv: "1112.3960", bibcode: "", isbn: "", url: "" }
retracted: false
verified: false
`,
		"works/1546/nobody-1546.yaml": `key: nobody-1546
tier: T5
kind: classic
title: "Untitled"
authors: ["Nobody, A."]
year: 1546
venue: ""
ids: { doi: "", arxiv: "", bibcode: "", isbn: "", url: "" }
retracted: false
verified: false
`,
		"works/2000/gone-2000.yaml": `key: gone-2000
tier: T1
kind: article
title: "Missing"
authors: ["Gone, B."]
year: 2000
venue: ""
ids: { doi: "10.0/missing", arxiv: "", bibcode: "", isbn: "", url: "" }
retracted: false
verified: false
`,
	}
	for p, s := range files {
		full := filepath.Join(root, p)
		os.MkdirAll(filepath.Dir(full), 0o755)
		if err := os.WriteFile(full, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	os.MkdirAll(filepath.Join(root, "nodes"), 0o755)
	return root
}

func TestVerifierWritesOnlyFullMatches(t *testing.T) {
	srv, _ := server(t)
	root := writeCorpus(t)
	c, err := content.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	v := &Verifier{
		Crossref: &Crossref{C: testClient(), Base: srv.URL},
		Arxiv:    &Arxiv{C: testClient(), APIBase: srv.URL + "/api"},
	}
	checks, err := v.Run(context.Background(), c, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]Status{}
	for _, ch := range checks {
		got[ch.Key] = ch.Status
	}
	want := map[string]Status{
		"anderson-1933":       StatusOK,
		"famaey-mcgaugh-2012": StatusMismatch,
		"nobody-1546":         StatusSkipped,
		"gone-2000":           StatusNotFound,
	}
	for k, st := range want {
		if got[k] != st {
			t.Errorf("%s: %s, want %s", k, got[k], st)
		}
	}
	read := func(p string) string {
		b, _ := os.ReadFile(filepath.Join(root, p))
		return string(b)
	}
	if s := read("works/1933/anderson-1933.yaml"); !strings.Contains(s, "\nverified: true\n") || !strings.Contains(s, `venue: "Physical Review 43, 491–494"`) {
		t.Fatalf("anderson not flipped cleanly:\n%s", s)
	}
	for _, p := range []string{"works/2012/famaey-mcgaugh-2012.yaml", "works/1546/nobody-1546.yaml", "works/2000/gone-2000.yaml"} {
		if !strings.Contains(read(p), "\nverified: false\n") {
			t.Errorf("%s must stay unverified", p)
		}
	}
}

func TestStoreDedup(t *testing.T) {
	st, err := OpenStore(filepath.Join(t.TempDir(), "h.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	ax := Record{Source: "arxiv", Arxiv: "1112.3960", Title: "MOND", Years: []int{2011}}
	cr := Record{Source: "crossref", DOI: "10.12942/lrr-2012-10", Arxiv: "1112.3960", Title: "MOND", Venue: "Living Reviews in Relativity"}
	if n, err := st.Upsert(ax, "t1"); err != nil || !n {
		t.Fatalf("first insert: %v %v", n, err)
	}
	if n, err := st.Upsert(cr, "t2"); err != nil || n {
		t.Fatalf("second should merge: %v %v", n, err)
	}
	if n, err := st.Upsert(Record{Source: "crossref", DOI: "https://doi.org/10.12942/LRR-2012-10"}, "t3"); err != nil || n {
		t.Fatalf("un-normalised DOI should merge: %v %v", n, err)
	}
	m, err := st.Stats()
	if err != nil {
		t.Fatal(err)
	}
	if m["arxiv"] != 1 {
		t.Fatalf("stats %v", m)
	}
	var doi, venue string
	st.db.QueryRow(`SELECT doi, venue FROM works_raw WHERE arxiv = '1112.3960'`).Scan(&doi, &venue)
	if doi != "10.12942/lrr-2012-10" || venue != "Living Reviews in Relativity" {
		t.Fatalf("merge lost fields: %q %q", doi, venue)
	}
}

func TestMatchCollaborationAndYear(t *testing.T) {
	w := &content.Work{Title: "Planck 2018 results. VI. Cosmological parameters", Authors: []string{"Planck Collaboration"},
		Year: 2020, IDs: content.WorkIDs{DOI: "10.1051/0004-6361/201833910"}}
	r := &Record{Title: "Planck 2018 results", Authors: []Author{{Family: "Aghanim"}}, Years: []int{2020},
		DOI: "10.1051/0004-6361/201833910"}
	r.Title += ": VI. Cosmological parameters" // as Crossref title + subtitle
	probs, notes := Match(w, r)
	if len(probs) != 0 || len(notes) != 1 {
		t.Fatalf("probs %v notes %v", probs, notes)
	}
	w.Year = 2018
	if probs, _ := Match(w, r); len(probs) != 1 || !strings.HasPrefix(probs[0], "year:") {
		t.Fatalf("year mismatch not reported: %v", probs)
	}
	w.Authors = []string{"Smith, J."}
	if probs, _ := Match(w, r); len(probs) != 2 {
		t.Fatalf("author mismatch not reported: %v", probs)
	}
}
