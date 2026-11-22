package evalset

import "fmt"

// Verdict is what the verify step returns for one item (pipeline/verify.py).
type Verdict struct {
	ID        string `json:"id"`
	Verdict   string `json:"verdict"` // yes | partial | no
	NumbersOK bool   `json:"numbers_ok"`
	Reason    string `json:"reason,omitempty"`
	Model     string `json:"model,omitempty"`
}

// Gate thresholds from plan §10 M2.
const (
	GateCitation = 0.95
	GateNumbers  = 1.00
	GateCatch    = 0.90
)

// Score holds the three go/no-go metrics.
type Score struct {
	Genuine, Planted int
	// CitationPrecision: genuine items the verifier accepts as supported.
	CitationPrecision float64
	// NumberAgreement: items with numbers where the verifier's number check matches the truth.
	NumberAgreement float64
	// Catch: planted items the verifier rejects (verdict ≠ yes).
	Catch   float64
	Missing []string // items without a verdict
	Misses  []string // ids the verifier got wrong
}

// Pass reports whether the score clears the gate.
func (s Score) Pass() bool {
	return len(s.Missing) == 0 && s.CitationPrecision >= GateCitation && s.NumberAgreement >= GateNumbers && s.Catch >= GateCatch
}

func (s Score) String() string {
	mark := func(v, gate float64) string {
		if v >= gate {
			return "✓"
		}
		return "✗"
	}
	return fmt.Sprintf("citation precision %.1f%% %s (≥ %.0f%%)\nnumber agreement   %.1f%% %s (= 100%%)\ncatch rate         %.1f%% %s (≥ %.0f%%)\n%d genuine, %d planted, %d without verdict",
		100*s.CitationPrecision, mark(s.CitationPrecision, GateCitation), 100*GateCitation,
		100*s.NumberAgreement, mark(s.NumberAgreement, GateNumbers),
		100*s.Catch, mark(s.Catch, GateCatch), 100*GateCatch,
		s.Genuine, s.Planted, len(s.Missing))
}

// Evaluate scores verdicts against items. hasNumbers tells which items
// contain numbers (in sentence or passages) and whether they truly agree.
func Evaluate(items []Item, verdicts []Verdict, numbersTruth map[string]bool) Score {
	byID := map[string]Verdict{}
	for _, v := range verdicts {
		byID[v.ID] = v
	}
	var s Score
	var genOK, catch, numItems, numOK int
	for _, it := range items {
		v, ok := byID[it.ID]
		if !ok {
			s.Missing = append(s.Missing, it.ID)
			continue
		}
		if truth, has := numbersTruth[it.ID]; has {
			numItems++
			if v.NumbersOK == truth {
				numOK++
			} else {
				s.Misses = append(s.Misses, it.ID+" (numbers)")
			}
		}
		if it.Planted() {
			s.Planted++
			if v.Verdict != "yes" {
				catch++
			} else {
				s.Misses = append(s.Misses, it.ID+" (planted accepted)")
			}
			continue
		}
		s.Genuine++
		if v.Verdict == it.Expected {
			genOK++
		} else {
			s.Misses = append(s.Misses, fmt.Sprintf("%s (expected %s, got %s)", it.ID, it.Expected, v.Verdict))
		}
	}
	ratio := func(a, b int) float64 {
		if b == 0 {
			return 1
		}
		return float64(a) / float64(b)
	}
	s.CitationPrecision = ratio(genOK, s.Genuine)
	s.NumberAgreement = ratio(numOK, numItems)
	s.Catch = ratio(catch, s.Planted)
	return s
}
