package lint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/lioilsources/kosmomat/internal/content"
)

var today = time.Date(2026, 11, 22, 0, 0, 0, 0, time.UTC)

const work = `key: %s
tier: %s
kind: article
title: "T"
authors: ["A"]
year: 2020
venue: "V"
ids: { doi: "%s", arxiv: "", bibcode: "", isbn: "", url: "" }
retracted: false
verified: true
`

func write(t *testing.T, root, rel, body string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T, nodeYAML, csBody, enBody string) []Issue {
	t.Helper()
	root := t.TempDir()
	write(t, root, "works/2020/good-2020.yaml", sprintf(work, "good-2020", "T1", "10.1/x"))
	write(t, root, "works/2020/other-2020.yaml", sprintf(work, "other-2020", "T1", "10.1/y"))
	write(t, root, "works/2020/blog-2020.yaml", sprintf(work, "blog-2020", "T6", ""))
	write(t, root, "nodes/term/peer/node.yaml", `id: term.peer
type: term
title: { cs: "P", en: "P" }
summary: { cs: "P", en: "P" }
edges: [ { rel: related, to: obs.x } ]
changelog: [ { date: 2026-11-22, kind: created, note: "x" } ]
last_verified: 2026-11-22
verified_by: human
`)
	write(t, root, "nodes/term/peer/cs.lay.md", "Text [@good-2020 s. 1].\n")
	write(t, root, "nodes/term/peer/en.lay.md", "Text [@good-2020 p. 1].\n")
	write(t, root, "nodes/observation/x/node.yaml", nodeYAML)
	write(t, root, "nodes/observation/x/cs.lay.md", csBody)
	write(t, root, "nodes/observation/x/en.lay.md", enBody)
	c, err := content.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	return Run(c, Options{Today: today})
}

func sprintf(f string, a ...any) string {
	r := f
	for _, x := range a {
		r = strings.Replace(r, "%s", x.(string), 1)
	}
	return r
}

const okNode = `id: obs.x
type: observation
status: established
literature: mainstream
title: { cs: "X", en: "X" }
summary: { cs: "X", en: "X" }
edges: [ { rel: related, to: term.peer } ]
changelog: [ { date: 2026-11-22, kind: created, note: "x" } ]
last_verified: 2026-11-22
verified_by: human
`

const okBody = "## O {#observation}\n\nFakt [@good-2020 s. 1; @other-2020 s. 2].\n\n## V {#interpretations}\n\n- výklad [@other-2020 s. 3]\n"

func rules(issues []Issue, sev Severity) map[string]int {
	m := map[string]int{}
	for _, i := range issues {
		if i.Sev == sev {
			m[i.Rule]++
		}
	}
	return m
}

func TestCleanFixture(t *testing.T) {
	is := fixture(t, okNode, okBody, okBody)
	if errs := rules(is, Error); len(errs) > 0 {
		for _, i := range is {
			t.Log(i)
		}
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestRulesFire(t *testing.T) {
	cases := []struct {
		name       string
		node, body string
		want       string
	}{
		{"L1 uncited paragraph", okNode, okBody + "\nBez zdroje.\n", "L1"},
		{"L1 missing locator", okNode, strings.Replace(okBody, "s. 3", "", 1), "L1"},
		{"L1 nocite not allowed here", okNode, strings.Replace(okBody, "{#interpretations}", "{#interpretations .nocite}", 1), "L1"},
		{"L2 unknown work", okNode, okBody + "\nX [@nope-1999 s. 1].\n", "L2"},
		{"L3 established with one source", okNode, "## O {#observation}\n\nX [@good-2020 s. 1].\n\n## V {#interpretations}\n\nY [@good-2020 s. 2].\n", "L3"},
		{"L3 contested one-sided", strings.Replace(okNode, "established", "contested", 1), okBody, "L3"},
		{"L4 missing section", okNode, "## O {#observation}\n\nX [@good-2020 s. 1; @other-2020 s. 1].\n", "L4"},
		{"L6 long summary", strings.Replace(okNode, `summary: { cs: "X"`, `summary: { cs: "`+strings.Repeat("x", 281)+`"`, 1), okBody, "L6"},
		{"L7 missing last_verified", strings.Replace(okNode, "last_verified: 2026-11-22\n", "", 1), okBody, "L7"},
		{"L9 long quote", okNode, okBody + "\n> " + strings.Repeat("slovo ", 21) + "[@good-2020 s. 9]\n", "L9"},
		{"L9 two quotes one work", okNode, okBody + "\n> a [@good-2020 s. 9]\n\nmezi [@good-2020 s. 1]\n\n> b [@good-2020 s. 10]\n", "L9"},
		{"L10 forbidden phrase", okNode, okBody + "\nVědci dokázali, že ano [@good-2020 s. 1].\n", "L10"},
		{"L10 pseudoscience unattributed", okNode, okBody + "\nJe to pavěda.\n", "L10"},
		{"L11 unsourced support", strings.Replace(okNode, "edges: [ { rel: related, to: term.peer } ]", "edges: [ { rel: supported_by, to: term.peer } ]", 1), okBody, "L11"},
		{"P4 T6 only", okNode, okBody + "\nOhlas [@blog-2020 odst. 2].\n", "P4"},
		{"S1 bad type", strings.Replace(okNode, "type: observation", "type: nonsense", 1), okBody, "S1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := fixture(t, tc.node, tc.body, tc.body)
			if rules(is, Error)[tc.want] == 0 {
				for _, i := range is {
					t.Log(i)
				}
				t.Fatalf("rule %s did not fire", tc.want)
			}
		})
	}
}

func TestL8Orphan(t *testing.T) {
	node := strings.Replace(okNode, "edges: [ { rel: related, to: term.peer } ]", "edges: []", 1)
	is := fixture(t, node, okBody, okBody)
	// term.peer still links to obs.x, so obs.x is not an orphan.
	if rules(is, Error)["L8"] != 0 {
		t.Fatalf("obs.x is linked from term.peer and must not be an orphan")
	}
}

func TestL7StaleContested(t *testing.T) {
	node := strings.Replace(strings.Replace(okNode, "established", "contested", 1), "last_verified: 2026-11-22", "last_verified: 2025-01-01", 1)
	is := fixture(t, node, okBody, okBody)
	if rules(is, Warning)["L7"] == 0 {
		t.Fatal("stale contested card should warn")
	}
}

func TestL2WorkWithoutIdentifier(t *testing.T) {
	c := &content.Corpus{Nodes: map[string]*content.Node{}, Works: map[string]*content.Work{
		"a-2020": {Key: "a-2020", Tier: "T1", Kind: "article", Title: "T", Authors: []string{"A"}, Year: 2020, File: "w.yaml"},
		"b-1859": {Key: "b-1859", Tier: "T1", Kind: "classic", Title: "T", Authors: []string{"B"}, Year: 1859, File: "w2.yaml"},
	}}
	is := Run(c, Options{Today: today})
	if rules(is, Error)["L2"] != 1 || rules(is, Warning)["L2"] != 1 {
		t.Fatalf("want 1 L2 error (article) and 1 L2 warning (classic): %v", is)
	}
}
