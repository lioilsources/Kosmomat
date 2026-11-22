package evalset

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/lioilsources/kosmomat/internal/content"
)

// testdata/numbers.json is shared with pipeline/test_verify.py so the Go and
// Python number checks cannot drift apart.
func TestNumbersAgreeSharedVectors(t *testing.T) {
	b, err := os.ReadFile("testdata/numbers.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Sentence string   `json:"sentence"`
		Passages []string `json:"passages"`
		Allowed  []string `json:"allowed"`
		OK       bool     `json:"ok"`
		Missing  []string `json:"missing"`
	}
	if err := json.Unmarshal(b, &cases); err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		ok, miss := NumbersAgree(c.Sentence, c.Passages, c.Allowed...)
		if ok != c.OK || !slices.Equal(miss, c.Missing) && len(miss)+len(c.Missing) > 0 {
			t.Errorf("%q: ok=%v missing=%v, want %v %v", c.Sentence, ok, miss, c.OK, c.Missing)
		}
	}
}

func corpus() *content.Corpus {
	md := `## A {#a}

Speeds stay roughly constant in 21 galaxies out to 122 kpc [@w1 abstract]. A second claim is more general [@w2 abstract; @w3 abstract].

Background without abstract [@w4 p. 3].
`
	c := &content.Corpus{Nodes: map[string]*content.Node{}, Works: map[string]*content.Work{}}
	for _, id := range []string{"obs.a", "obs.b"} {
		c.Nodes[id] = &content.Node{ID: id, Docs: map[string]*content.Doc{"en.lay": {Blocks: content.ParseMarkdown(md)}}}
	}
	return c
}

func TestSegmentsAndCandidates(t *testing.T) {
	segs := Segments(corpus(), "en")
	if len(segs) != 6 {
		t.Fatalf("segments %d: %+v", len(segs), segs)
	}
	if segs[0].Text != "Speeds stay roughly constant in 21 galaxies out to 122 kpc." || segs[1].Text != "A second claim is more general." {
		t.Fatalf("texts %q / %q", segs[0].Text, segs[1].Text)
	}
	abs := map[string]string{"w1": "Rotation speeds are flat for 21 galaxies to 122 kpc.", "w2": "x", "w4": "y"}
	items, missing := Candidates(segs, abs)
	// w3 has no abstract → second segment skipped; w4 is cited by page → not a candidate.
	if len(items) != 2 || !slices.Equal(missing, []string{"w3"}) {
		t.Fatalf("items %+v missing %v", items, missing)
	}
	if items[0].ID != "obs.a#1" || items[0].Expected != "yes" || items[0].VerifiedBy != "draft" || items[0].PassageSHA != SHA([]string{abs["w1"]}) {
		t.Fatalf("item %+v", items[0])
	}
}

func TestPlantDeterministicAndWrong(t *testing.T) {
	var gen []Item
	for i, s := range []string{
		"Speeds stay roughly constant in 21 galaxies out to 122 kpc.",
		"The signal is found at 4.2 sigma in 1998.",
		"Lensing mass is larger than the gas mass.",
		"The core radius is 1,830 km.",
		"Neutrinos were detected before the model.",
		"The effect shows in 153 galaxies.",
	} {
		gen = append(gen, Item{ID: "n#" + string(rune('a'+i)), Node: []string{"n1", "n2", "n3"}[i%3], Sentence: s,
			Cites: []string{"w" + string(rune('a'+i))}, PassageSHA: "sha" + string(rune('a'+i)), Expected: "yes"})
	}
	a, b := Plant(gen, 6, 7), Plant(gen, 6, 7)
	if len(a) != 6 {
		t.Fatalf("planted %d", len(a))
	}
	for i := range a {
		if a[i].Sentence != b[i].Sentence || a[i].PassageSHA != b[i].PassageSHA {
			t.Fatal("not deterministic")
		}
		p := a[i]
		if p.Expected != "no" || p.Mutation == nil {
			t.Fatalf("planted item %+v", p)
		}
		switch p.Mutation.Kind {
		case "number":
			if isYear(p.Mutation.From) || p.Mutation.From == p.Mutation.To || !strings.Contains(p.Sentence, p.Mutation.To) {
				t.Errorf("bad number mutation %+v", p.Mutation)
			}
		case "reversal":
			if !strings.Contains(p.Sentence, p.Mutation.To) {
				t.Errorf("bad reversal %+v", p)
			}
		case "wrong-cite":
			var orig Item
			for _, g := range gen {
				if g.ID == p.Mutation.Of {
					orig = g
				}
			}
			if orig.Node == "" || p.Sentence != orig.Sentence || slices.Equal(p.Cites, orig.Cites) {
				t.Errorf("bad wrong-cite %+v", p)
			}
		}
	}
}

func TestEvaluate(t *testing.T) {
	items := []Item{
		{ID: "g1", Sentence: "21 galaxies", Expected: "yes"},
		{ID: "g2", Sentence: "flat", Expected: "yes"},
		{ID: "p1", Sentence: "42 galaxies", Expected: "no", Mutation: &Mutation{Kind: "number"}},
		{ID: "p2", Sentence: "not flat", Expected: "no", Mutation: &Mutation{Kind: "reversal"}},
	}
	truth := map[string]bool{"g1": true, "p1": false}
	s := Evaluate(items, []Verdict{
		{ID: "g1", Verdict: "yes", NumbersOK: true},
		{ID: "g2", Verdict: "yes"},
		{ID: "p1", Verdict: "partial", NumbersOK: false},
		{ID: "p2", Verdict: "no"},
	}, truth)
	if !s.Pass() || s.CitationPrecision != 1 || s.Catch != 1 || s.NumberAgreement != 1 {
		t.Fatalf("score %+v", s)
	}
	s = Evaluate(items, []Verdict{{ID: "g1", Verdict: "yes", NumbersOK: true}, {ID: "p1", Verdict: "yes", NumbersOK: true}}, truth)
	if s.Pass() || len(s.Missing) != 2 || s.Catch != 0 || s.NumberAgreement != 0.5 {
		t.Fatalf("score %+v", s)
	}
}
