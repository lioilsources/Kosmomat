package pack

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/lioilsources/kosmomat/internal/content"
)

// WebNode is the per-card JSON the Astro site renders at /n/<id>.
type WebNode struct {
	ID            string                                `json:"id"`
	Type          string                                `json:"type"`
	Status        string                                `json:"status,omitempty"`
	Literature    string                                `json:"literature,omitempty"`
	Testability   string                                `json:"testability,omitempty"`
	CaseKind      string                                `json:"case_kind,omitempty"`
	Astromat      bool                                  `json:"astromat,omitempty"`
	Title         content.I18n                          `json:"title"`
	Summary       content.I18n                          `json:"summary"`
	FirstProposed *content.Origin                       `json:"first_proposed,omitempty"`
	Body          map[string]map[string][]content.Block `json:"body"` // lang → level → blocks
	EdgesOut      []content.Edge                        `json:"edges_out"`
	EdgesIn       []WebInEdge                           `json:"edges_in"`
	Predictions   []content.Prediction                  `json:"predictions,omitempty"`
	Works         []*content.Work                       `json:"works"`
	Changelog     []content.ChangelogEntry              `json:"changelog"`
	LastVerified  string                                `json:"last_verified"`
	VerifiedBy    string                                `json:"verified_by"`
}

// WebInEdge is an edge pointing at this node from another one.
type WebInEdge struct {
	From string `json:"from"`
	Rel  string `json:"rel"`
}

// WebIndexEntry is one row of web/index.json.
type WebIndexEntry struct {
	ID         string       `json:"id"`
	Type       string       `json:"type"`
	Status     string       `json:"status,omitempty"`
	Literature string       `json:"literature,omitempty"`
	Title      content.I18n `json:"title"`
	Summary    content.I18n `json:"summary"`
}

func exportWeb(c *content.Corpus, dir string, m *Manifest) error {
	if err := os.MkdirAll(filepath.Join(dir, "n"), 0o755); err != nil {
		return fmt.Errorf("mkdir web: %w", err)
	}
	in := map[string][]WebInEdge{}
	for _, id := range c.NodeIDs() {
		for _, e := range c.Nodes[id].Edges {
			in[e.To] = append(in[e.To], WebInEdge{From: id, Rel: e.Rel})
		}
		if n := c.Nodes[id]; n.Type == "prediction" && n.Theory != "" {
			in[n.Theory] = append(in[n.Theory], WebInEdge{From: id, Rel: "predicted_by"})
		}
	}
	preds := c.AllPredictions()
	var index []WebIndexEntry
	for _, id := range c.NodeIDs() {
		n := c.Nodes[id]
		w := WebNode{
			ID: n.ID, Type: n.Type, Status: n.Status, Literature: n.Literature, Testability: n.Testability,
			CaseKind: n.CaseKind, Astromat: n.Astromat, Title: n.Title, Summary: n.Summary,
			FirstProposed: n.FirstProposed, Body: map[string]map[string][]content.Block{},
			EdgesOut: n.Edges, EdgesIn: in[id], Changelog: n.Changelog,
			LastVerified: n.LastVerified, VerifiedBy: n.VerifiedBy,
		}
		if w.EdgesOut == nil {
			w.EdgesOut = []content.Edge{}
		}
		if w.EdgesIn == nil {
			w.EdgesIn = []WebInEdge{}
		}
		cited := map[string]bool{}
		for _, d := range n.Docs {
			if w.Body[d.Lang] == nil {
				w.Body[d.Lang] = map[string][]content.Block{}
			}
			w.Body[d.Lang][d.Level] = d.Blocks
			for _, b := range d.Blocks {
				for _, p := range b.Paragraphs() {
					for _, ct := range content.Cites(p) {
						cited[ct.W] = true
					}
				}
			}
		}
		for _, e := range n.Edges {
			if e.Source != nil {
				cited[content.WorkKey(e.Source.Work)] = true
			}
		}
		for _, p := range preds {
			if p.Theory == id || p.ID == id {
				w.Predictions = append(w.Predictions, p)
				for _, ref := range []string{originWork(p.Made), testWork(p.Test)} {
					if ref != "" {
						cited[content.WorkKey(ref)] = true
					}
				}
			}
		}
		if n.FirstProposed != nil && n.FirstProposed.In != "" {
			cited[content.WorkKey(n.FirstProposed.In)] = true
		}
		for _, k := range sortedKeys(cited) {
			if wk, ok := c.Works[k]; ok {
				w.Works = append(w.Works, wk)
			}
		}
		sort.Slice(w.Works, func(i, j int) bool { return w.Works[i].Key < w.Works[j].Key })
		if err := writeJSON(filepath.Join(dir, "n", id+".json"), w); err != nil {
			return err
		}
		index = append(index, WebIndexEntry{ID: id, Type: n.Type, Status: n.Status, Literature: n.Literature, Title: n.Title, Summary: n.Summary})
	}
	return writeJSON(filepath.Join(dir, "index.json"), map[string]any{
		"version":   m.Version,
		"nodes":     index,
		"redirects": c.Redirects,
	})
}

func originWork(o *content.Origin) string {
	if o == nil {
		return ""
	}
	return o.In
}

func testWork(t *content.Test) string {
	if t == nil {
		return ""
	}
	return t.In
}
