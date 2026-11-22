package evalset

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
)

var number = regexp.MustCompile(`\b\d+(?:[.,]\d+)?\b`)

// Numbers returns the numeric tokens of s, normalised ("1,5" → "1.5", "10,000" → "10000").
func Numbers(s string) []string {
	var out []string
	for _, m := range number.FindAllString(s, -1) {
		out = append(out, normNumber(m))
	}
	return out
}

func normNumber(m string) string {
	if i := strings.IndexByte(m, ','); i >= 0 {
		if len(m)-i-1 == 3 { // thousands separator
			return strings.ReplaceAll(m, ",", "")
		}
		return strings.Replace(m, ",", ".", 1)
	}
	return m
}

// NumbersAgree reports whether every number in the sentence also occurs in the
// passages or among allowed (the cited works' publication years, which come from
// citation metadata, not from the abstract). This is the deterministic
// "číselná shoda" check; verify must not leave it to a model.
func NumbersAgree(sentence string, passages []string, allowed ...string) (bool, []string) {
	have := map[string]bool{}
	for _, a := range allowed {
		have[a] = true
	}
	for _, p := range passages {
		for _, n := range Numbers(p) {
			have[n] = true
		}
	}
	var missing []string
	for _, n := range Numbers(sentence) {
		if !have[n] {
			missing = append(missing, n)
		}
	}
	return len(missing) == 0, missing
}

// reversals flip the claim; the first that matches a word in the sentence is used.
var reversals = [][2]string{
	{"stays roughly constant", "falls steadily"}, {"constant", "falling"}, {"before", "after"},
	{"more", "less"}, {"less", "more"}, {"higher", "lower"}, {"lower", "higher"},
	{"larger", "smaller"}, {"smaller", "larger"}, {"faster", "slower"}, {"slower", "faster"},
	{"increase", "decrease"}, {"decrease", "increase"}, {"confirmed", "ruled out"},
	{"found", "did not find"}, {"shows", "does not show"}, {"predicted", "failed to predict"},
	{"is", "is not"}, {"are", "are not"}, {"was", "was not"}, {"were", "were not"},
}

func reverse(s string) (from, to string, ok bool) {
	for _, r := range reversals {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(r[0]) + `\b`)
		if re.MatchString(s) {
			return r[0], r[1], true
		}
	}
	return "", "", false
}

func isYear(n string) bool {
	v, err := strconv.Atoi(n)
	return err == nil && len(n) == 4 && v >= 1500 && v <= 2100
}

func bump(n string) string {
	if f, err := strconv.ParseFloat(n, 64); err == nil && strings.Contains(n, ".") {
		return strconv.FormatFloat(f*1.4, 'f', len(n)-strings.IndexByte(n, '.')-1, 64)
	}
	v, _ := strconv.Atoi(n)
	switch {
	case isYear(n):
		return strconv.Itoa(v + 7)
	case v < 3:
		return strconv.Itoa(v + 3)
	default:
		return strconv.Itoa(v * 2)
	}
}

// Plant makes n deliberately wrong variants of genuine items, cycling
// number → reversal → wrong-cite. Deterministic for a given seed and input.
// A wrong-cite variant takes the citations of an item from another node.
func Plant(genuine []Item, n int, seed uint64) []Item {
	r := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	order := r.Perm(len(genuine))
	kinds := []string{"number", "reversal", "wrong-cite"}
	used := map[string]bool{}
	var out []Item
	for len(out) < n {
		progressed := false
		for _, kind := range kinds {
			if len(out) >= n {
				break
			}
			for _, i := range order {
				g := genuine[i]
				if used[g.ID] || g.Planted() {
					continue
				}
				p, ok := plantOne(g, kind, genuine, r)
				if !ok {
					continue
				}
				used[g.ID] = true
				p.ID = fmt.Sprintf("planted-%02d", len(out)+1)
				out = append(out, p)
				progressed = true
				break
			}
		}
		if !progressed {
			break
		}
	}
	return out
}

// NumbersTruth is the expected outcome of the number check: genuine items
// (human-verified) agree, number mutations do not, other plants are computed.
func NumbersTruth(it Item, passages []string, allowed ...string) (truth, has bool) {
	if len(Numbers(it.Sentence)) == 0 {
		return false, false
	}
	switch {
	case !it.Planted():
		return true, true
	case it.Mutation.Kind == "number":
		return false, true
	}
	ok, _ := NumbersAgree(it.Sentence, passages, allowed...)
	return ok, true
}

func plantOne(g Item, kind string, all []Item, r *rand.Rand) (Item, bool) {
	p := g
	p.Expected, p.VerifiedBy, p.Note = "no", "draft", ""
	switch kind {
	case "number":
		// Prefer quantities over years: a changed year is a metadata error, not a passage error.
		var nums [][]int
		for _, loc := range number.FindAllStringIndex(g.Sentence, -1) {
			if !isYear(g.Sentence[loc[0]:loc[1]]) {
				nums = append(nums, loc)
			}
		}
		if len(nums) == 0 {
			return p, false
		}
		loc := nums[r.IntN(len(nums))]
		from := g.Sentence[loc[0]:loc[1]]
		to := bump(normNumber(from))
		if strings.Contains(from, ",") && !strings.Contains(to, ".") { // keep "1,830" style readable
			to = thousands(to)
		}
		p.Sentence = g.Sentence[:loc[0]] + to + g.Sentence[loc[1]:]
		p.Mutation = &Mutation{Kind: kind, From: from, To: to, Of: g.ID}
	case "reversal":
		from, to, ok := reverse(g.Sentence)
		if !ok {
			return p, false
		}
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(from) + `\b`)
		done := false
		p.Sentence = re.ReplaceAllStringFunc(g.Sentence, func(m string) string {
			if done {
				return m
			}
			done = true
			return to
		})
		p.Mutation = &Mutation{Kind: kind, From: from, To: to, Of: g.ID}
	case "wrong-cite":
		var others []Item
		for _, o := range all {
			if o.Node != g.Node && !o.Planted() && !sameCites(o.Cites, g.Cites) {
				others = append(others, o)
			}
		}
		if len(others) == 0 {
			return p, false
		}
		o := others[r.IntN(len(others))]
		p.Cites, p.PassageSHA = o.Cites, o.PassageSHA
		p.Mutation = &Mutation{Kind: kind, From: strings.Join(g.Cites, ";"), To: strings.Join(o.Cites, ";"), Of: g.ID}
	}
	return p, true
}

func thousands(n string) string {
	for i := len(n) - 3; i > 0; i -= 3 {
		n = n[:i] + "," + n[i:]
	}
	return n
}

func sameCites(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				return true
			}
		}
	}
	return false
}
