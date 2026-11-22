package harvest

import (
	"fmt"
	"slices"
	"strings"

	"github.com/lioilsources/kosmomat/internal/content"
)

// Match compares a curated work with a source record and lists every
// disagreement in title, first author, year and cross identifiers.
// An empty result means the metadata agrees. Notes explain checks that were
// relaxed (collaboration authorship) so a human can see why it passed.
func Match(w *content.Work, r *Record) (probs, notes []string) {
	if !titlesAgree(w.Title, r.Title) {
		probs = append(probs, fmt.Sprintf("title: %q ≠ %q", w.Title, r.Title))
	}
	if len(w.Authors) == 0 {
		probs = append(probs, "authors: none in work")
	} else if isCollaboration(w.Authors[0]) {
		// Sources list either the collaboration or its alphabetically first member;
		// title + year + identifier carry the match.
		notes = append(notes, fmt.Sprintf("collaboration author %q not compared", w.Authors[0]))
	} else if !firstAuthorAgrees(w.Authors[0], r.Authors) {
		got := "∅"
		if len(r.Authors) > 0 {
			got = r.Authors[0].Family
		}
		probs = append(probs, fmt.Sprintf("first author: %q ≠ %q", w.Authors[0], got))
	}
	if !slices.Contains(r.Years, w.Year) {
		probs = append(probs, fmt.Sprintf("year: %d ∉ %v", w.Year, r.Years))
	}
	if w.IDs.DOI != "" && r.DOI != "" && NormDOI(w.IDs.DOI) != r.DOI {
		probs = append(probs, fmt.Sprintf("doi: %q ≠ %q", w.IDs.DOI, r.DOI))
	}
	if w.IDs.Arxiv != "" && r.Arxiv != "" && NormArxiv(w.IDs.Arxiv) != r.Arxiv {
		probs = append(probs, fmt.Sprintf("arxiv: %q ≠ %q", w.IDs.Arxiv, r.Arxiv))
	}
	if r.Retracted && !w.Retracted {
		probs = append(probs, "retracted at source but not in work")
	}
	return probs, notes
}

func isCollaboration(a string) bool {
	return strings.Contains(strings.ToLower(a), "collaboration")
}

// titlesAgree accepts equal folded titles, or one being the other plus a
// subtitle (the shorter must keep at least half the length).
func titlesAgree(a, b string) bool {
	fa, fb := Fold(a), Fold(b)
	if fa == "" || fb == "" {
		return false
	}
	if fa == fb {
		return true
	}
	short, long := fa, fb
	if len(short) > len(long) {
		short, long = long, short
	}
	return strings.HasPrefix(long, short+" ") && 2*len(short) >= len(long)
}

// surname takes "Anderson, C. D." → "anderson"; a name without a comma
// (e.g. "Planck Collaboration") is used whole.
func surname(s string) string {
	fam, _, _ := strings.Cut(s, ",")
	return Fold(fam)
}

// firstAuthorAgrees looks for the work's first author among the first three
// record authors (collaboration papers often list a person first at one source
// and the collaboration at another).
func firstAuthorAgrees(first string, authors []Author) bool {
	want := surname(first)
	if want == "" {
		return false
	}
	for i, a := range authors {
		if i == 3 {
			break
		}
		if Fold(a.Family) == want {
			return true
		}
	}
	return false
}
