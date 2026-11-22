package harvest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/lioilsources/kosmomat/internal/content"
)

// Status of one work check.
type Status string

const (
	StatusOK        Status = "ok"
	StatusMismatch  Status = "mismatch"
	StatusNotFound  Status = "not-found"
	StatusSkipped   Status = "skipped"
	StatusError     Status = "error"
	StatusRetracted Status = "retracted"
)

// Check is the outcome for one work.
type Check struct {
	Key      string   `json:"key"`
	File     string   `json:"file"`
	Status   Status   `json:"status"`
	Source   string   `json:"source,omitempty"`
	Problems []string `json:"problems,omitempty"`
	Notes    []string `json:"notes,omitempty"`
	Venue    string   `json:"venue,omitempty"` // as the source reports it, for the human
	Written  bool     `json:"written,omitempty"`
}

// Verifier checks content/works against Crossref, arXiv and (with a token) SciX.
type Verifier struct {
	Crossref *Crossref
	Arxiv    *Arxiv
	Scix     *Scix // nil without token
	Log      func(format string, args ...any)
}

func (v *Verifier) logf(f string, a ...any) {
	if v.Log != nil {
		v.Log(f, a...)
	}
}

// Run checks every not-yet-verified work (or the given keys). With write,
// a work whose metadata fully agrees gets `verified: true` in its file.
// Nothing else is ever written: retractions and mismatches are for a human.
func (v *Verifier) Run(ctx context.Context, c *content.Corpus, keys []string, write bool) ([]Check, error) {
	if len(keys) == 0 {
		for k, w := range c.Works {
			if !w.Verified {
				keys = append(keys, k)
			}
		}
	}
	sort.Strings(keys)
	var out []Check
	for _, k := range keys {
		w, ok := c.Works[k]
		if !ok {
			return out, fmt.Errorf("unknown work %q", k)
		}
		if err := ctx.Err(); err != nil {
			return out, err
		}
		ch := v.check(ctx, w)
		v.logf("%-10s %-32s %s %v", ch.Status, k, ch.Source, ch.Problems)
		if write && ch.Status == StatusOK {
			if err := markVerified(filepath.Join(c.Root, w.File)); err != nil {
				return out, err
			}
			ch.Written = true
		}
		out = append(out, ch)
	}
	return out, nil
}

func (v *Verifier) check(ctx context.Context, w *content.Work) Check {
	ch := Check{Key: w.Key, File: w.File}
	type attempt struct {
		source string
		fetch  func() (*Record, error)
	}
	var tries []attempt
	if w.IDs.DOI != "" {
		tries = append(tries, attempt{"crossref", func() (*Record, error) { return v.Crossref.Work(ctx, w.IDs.DOI) }})
	}
	if w.IDs.Arxiv != "" {
		tries = append(tries, attempt{"arxiv", func() (*Record, error) {
			m, err := v.Arxiv.ByIDs(ctx, []string{w.IDs.Arxiv})
			if err != nil {
				return nil, err
			}
			r, ok := m[NormArxiv(w.IDs.Arxiv)]
			if !ok {
				return nil, ErrNotFound
			}
			return &r, nil
		}})
	}
	if w.IDs.Bibcode != "" && v.Scix != nil {
		tries = append(tries, attempt{"scix", func() (*Record, error) { return v.Scix.ByBibcode(ctx, w.IDs.Bibcode) }})
	}
	if len(tries) == 0 {
		ch.Status = StatusSkipped
		switch {
		case w.IDs.Bibcode != "":
			ch.Problems = []string{"bibcode only; set ADS_API_TOKEN to check against SciX"}
		case w.IDs.ISBN != "":
			ch.Problems = []string{"book (ISBN): check by hand"}
		default:
			ch.Problems = []string{"no DOI/arXiv/bibcode: check by hand"}
		}
		return ch
	}
	// The first source that has the record decides; a not-found falls through to the next.
	ch.Status = StatusNotFound
	for _, t := range tries {
		r, err := t.fetch()
		ch.Source = t.source
		if errors.Is(err, ErrNotFound) {
			ch.Problems = append(ch.Problems, t.source+": not found")
			continue
		}
		if err != nil {
			ch.Status, ch.Problems = StatusError, append(ch.Problems, err.Error())
			return ch
		}
		ch.Venue = r.Venue
		ch.Problems, ch.Notes = Match(w, r)
		switch {
		case r.Retracted && !w.Retracted:
			ch.Status = StatusRetracted
		case len(ch.Problems) == 0:
			ch.Status = StatusOK
		default:
			ch.Status = StatusMismatch
		}
		return ch
	}
	return ch
}

var verifiedFalse = regexp.MustCompile(`(?m)^verified: false[ \t]*$`)

// markVerified flips the single top-level `verified: false` line, keeping the rest byte-identical.
func markVerified(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if n := len(verifiedFalse.FindAllIndex(b, -1)); n != 1 {
		return fmt.Errorf("%s: expected one `verified: false` line, found %d", path, n)
	}
	return os.WriteFile(path, verifiedFalse.ReplaceAll(b, []byte("verified: true")), 0o644)
}

// Counts tallies checks by status.
func Counts(cs []Check) map[Status]int {
	m := map[Status]int{}
	for _, c := range cs {
		m[c.Status]++
	}
	return m
}
