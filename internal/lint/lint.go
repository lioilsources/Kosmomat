// Package lint enforces docs/EDITORIAL.md on a loaded corpus.
// Rule ids L1–L11 follow the plan (§7); S-rules are structural checks.
package lint

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/lioilsources/kosmomat/internal/content"
)

// Severity of an issue. Errors fail the gate, warnings do not.
type Severity int

const (
	Warning Severity = iota
	Error
)

func (s Severity) String() string {
	if s == Error {
		return "ERROR"
	}
	return "WARN"
}

// Issue is a single finding.
type Issue struct {
	Sev  Severity
	Rule string
	Node string
	File string
	Line int
	Msg  string
}

func (i Issue) String() string {
	loc := i.File
	if i.Line > 0 {
		loc = fmt.Sprintf("%s:%d", i.File, i.Line)
	}
	return fmt.Sprintf("%-5s %-4s %-36s %s  %s", i.Sev, i.Rule, i.Node, loc, i.Msg)
}

// Options configure a lint run.
type Options struct {
	// Today is the reference date for freshness rules (L7) and the 10-year window (L3).
	Today time.Time
}

type linter struct {
	c      *content.Corpus
	opt    Options
	issues []Issue
}

// Run lints the corpus and returns issues sorted by node, file and line.
func Run(c *content.Corpus, opt Options) []Issue {
	l := &linter{c: c, opt: opt}
	for _, e := range c.LoadErrors {
		l.add(Error, "S0", "", e.File, 0, "%s", e.Err)
	}
	for _, k := range c.WorkKeys() {
		l.work(c.Works[k])
	}
	for _, id := range c.NodeIDs() {
		l.node(c.Nodes[id])
	}
	l.graph()
	l.media()
	sort.SliceStable(l.issues, func(i, j int) bool {
		a, b := l.issues[i], l.issues[j]
		if a.Node != b.Node {
			return a.Node < b.Node
		}
		if a.File != b.File {
			return a.File < b.File
		}
		return a.Line < b.Line
	})
	return l.issues
}

// Count returns the number of errors and warnings.
func Count(issues []Issue) (errs, warns int) {
	for _, i := range issues {
		if i.Sev == Error {
			errs++
		} else {
			warns++
		}
	}
	return errs, warns
}

func (l *linter) add(sev Severity, rule, node, file string, line int, format string, args ...any) {
	l.issues = append(l.issues, Issue{Sev: sev, Rule: rule, Node: node, File: file, Line: line, Msg: fmt.Sprintf(format, args...)})
}

var (
	typePrefix = map[string]string{
		"question": "q", "observation": "obs", "claim": "claim", "theory": "theory",
		"prediction": "pred", "case": "case", "person": "person", "mission": "mission",
		"term": "term", "experiment": "exp", "era": "era", "worldview": "worldview",
		"body": "body", "constellation": "constellation",
	}
	statusRequired = map[string]bool{"observation": true, "claim": true, "theory": true, "prediction": true}
	statuses       = set("established", "strong", "contested", "speculative", "refuted", "historical", "philosophical")
	literatures    = set("mainstream", "minority", "outside")
	testabilities  = set("empirical", "partial", "philosophical")
	caseKinds      = set("vindicated", "not_vindicated", "quick_fix", "pair", "myth")
	outcomes       = set("confirmed", "refuted", "partial", "open", "untestable")
	tiers          = set("T1", "T2", "T3", "T4", "T5", "T6")
	rels           = set("first_proposed_in", "proposed_by", "builds_on", "inspired_by", "supported_by",
		"contradicted_by", "interprets", "competes_with", "superseded_by", "predicts", "tested_by",
		"answers", "opposed_by", "measured_by", "related")
	sourcedRels = set("supported_by", "contradicted_by", "opposed_by")
	reSlug      = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	reWorkKey   = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

	requiredSections = map[string][]string{
		"theory":        {"claims", "predictions", "open_problems", "competitors", "would_refute"},
		"case":          {"context", "who", "objections", "decided", "duration", "lesson"},
		"question":      {"myth", "historical", "science"},
		"experiment":    {"needs", "steps", "you_will_see", "proves", "does_not_prove"},
		"body":          {"what", "distance", "sky_today", "missions", "open_question", "astrology_link"},
		"constellation": {"what", "distance", "sky_today", "missions", "open_question", "sign_vs_constellation", "astrology_link"},
		"observation":   {"observation", "interpretations"},
	}
	nociteAllowed = map[string]map[string]bool{
		"experiment":    set("needs", "steps", "you_will_see"),
		"body":          set("sky_today", "astrology_link"),
		"constellation": set("sky_today", "astrology_link"),
		"question":      set("checklist"),
	}
)

func set(xs ...string) map[string]bool {
	m := make(map[string]bool, len(xs))
	for _, x := range xs {
		m[x] = true
	}
	return m
}

func (l *linter) work(w *content.Work) {
	const node = ""
	if !reWorkKey.MatchString(w.Key) {
		l.add(Error, "S1", node, w.File, 0, "work key %q is not lowercase ascii with dashes", w.Key)
	}
	if !tiers[w.Tier] {
		l.add(Error, "S1", node, w.File, 0, "work %s: unknown tier %q", w.Key, w.Tier)
	}
	if w.Title == "" || len(w.Authors) == 0 || w.Year == 0 {
		l.add(Error, "S1", node, w.File, 0, "work %s: title, authors and year are required", w.Key)
	}
	if (w.Tier == "T1" || w.Tier == "T2") && w.IDs.DOI == "" && w.IDs.Arxiv == "" && w.IDs.Bibcode == "" {
		if w.Kind == "classic" {
			l.add(Warning, "L2", node, w.File, 0, "historical work %s has no DOI/bibcode; add a scan URL", w.Key)
		} else {
			l.add(Error, "L2", node, w.File, 0, "work %s is %s but has no DOI, arXiv id or bibcode", w.Key, w.Tier)
		}
	}
	if w.Tier == "T6" && w.IDs.URL == "" && w.IDs.ISBN == "" {
		l.add(Warning, "L2", node, w.File, 0, "work %s (T6) has no URL or ISBN", w.Key)
	}
	if !w.Verified {
		l.add(Warning, "V1", node, w.File, 0, "work %s: metadata not yet verified against Crossref/SciX/arXiv", w.Key)
	}
}

func (l *linter) node(n *content.Node) {
	file := n.Dir + "/node.yaml"
	l.structure(n, file)
	l.texts(n, file)
	l.edges(n, file)
	l.predictions(n, file)
	l.freshness(n, file)
	l.status(n, file)
	l.docs(n)
}

func (l *linter) structure(n *content.Node, file string) {
	prefix, ok := typePrefix[n.Type]
	if !ok {
		l.add(Error, "S1", n.ID, file, 0, "unknown type %q", n.Type)
		return
	}
	p, slug, _ := strings.Cut(n.ID, ".")
	if p != prefix || !reSlug.MatchString(slug) {
		l.add(Error, "S1", n.ID, file, 0, "id must be %s.<slug> for type %s", prefix, n.Type)
	}
	if want := "nodes/" + n.Type + "/" + slug; n.Dir != want {
		l.add(Error, "S1", n.ID, file, 0, "node lives in %s, expected %s", n.Dir, want)
	}
	if statusRequired[n.Type] && n.Status == "" {
		l.add(Error, "S2", n.ID, file, 0, "status is required for %s", n.Type)
	}
	if n.Status != "" && !statuses[n.Status] {
		l.add(Error, "S2", n.ID, file, 0, "unknown status %q", n.Status)
	}
	if statusRequired[n.Type] && n.Literature == "" {
		l.add(Error, "S2", n.ID, file, 0, "literature is required for %s (P2)", n.Type)
	}
	if n.Literature != "" && !literatures[n.Literature] {
		l.add(Error, "S2", n.ID, file, 0, "unknown literature %q", n.Literature)
	}
	if n.Type == "question" && !testabilities[n.Testability] {
		l.add(Error, "S2", n.ID, file, 0, "question needs testability empirical|partial|philosophical")
	}
	if n.Type == "case" && !caseKinds[n.CaseKind] {
		l.add(Error, "S2", n.ID, file, 0, "case needs case_kind (%s)", keys(caseKinds))
	}
	if (n.Type == "body" || n.Type == "constellation") && !n.Astromat {
		l.add(Warning, "S2", n.ID, file, 0, "%s card is not marked astromat: true (axis I entry card)", n.Type)
	}
	if len(n.Changelog) == 0 {
		l.add(Error, "S3", n.ID, file, 0, "changelog must have at least one entry (P12)")
	}
}

func (l *linter) texts(n *content.Node, file string) {
	for _, lang := range content.Langs {
		if n.Title[lang] == "" {
			l.add(Error, "L6", n.ID, file, 0, "title.%s missing", lang)
		}
		s := n.Summary[lang]
		if s == "" {
			l.add(Error, "L6", n.ID, file, 0, "summary.%s missing", lang)
		} else if c := utf8.RuneCountInString(s); c > 280 {
			l.add(Error, "L6", n.ID, file, 0, "summary.%s has %d characters (max 280)", lang, c)
		}
		if _, ok := n.Docs[lang+".lay"]; !ok {
			l.add(Error, "L6", n.ID, n.Dir, 0, "%s.lay.md missing", lang)
		}
		l.phrases(n, file, 0, n.Summary[lang])
	}
}

func (l *linter) workRef(n *content.Node, file, rule, ref, what string) *content.Work {
	if ref == "" {
		return nil
	}
	w, ok := l.c.Works[content.WorkKey(ref)]
	if !ok {
		l.add(Error, rule, n.ID, file, 0, "%s cites unknown work %q", what, ref)
		return nil
	}
	if w.Retracted {
		l.add(Error, rule, n.ID, file, 0, "%s cites retracted work %q", what, ref)
	}
	return w
}

func (l *linter) edges(n *content.Node, file string) {
	if fp := n.FirstProposed; fp != nil {
		if fp.In == "" {
			l.add(Error, "L8", n.ID, file, 0, "first_proposed.in must name a work")
		}
		l.workRef(n, file, "L2", fp.In, "first_proposed")
	}
	for i, e := range n.Edges {
		what := fmt.Sprintf("edge #%d %s→%s", i+1, e.Rel, e.To)
		if !rels[e.Rel] {
			l.add(Error, "S4", n.ID, file, 0, "%s: unknown rel", what)
		}
		if e.To == n.ID {
			l.add(Error, "S4", n.ID, file, 0, "%s: self-loop", what)
		}
		if _, ok := l.c.Nodes[e.To]; !ok {
			l.add(Warning, "S5", n.ID, file, 0, "%s: target node not written yet", what)
		}
		if sourcedRels[e.Rel] {
			if e.Source == nil || e.Source.Work == "" {
				l.add(Error, "L11", n.ID, file, 0, "%s needs source {work, loc}", what)
				continue
			}
			if e.Source.Loc == "" {
				l.add(Error, "L11", n.ID, file, 0, "%s: source has no locator", what)
			}
		}
		if e.Source != nil {
			if w := l.workRef(n, file, "L2", e.Source.Work, what); w != nil && w.Tier == "T6" {
				l.add(Error, "P4", n.ID, file, 0, "%s: T6 work %s cannot support a claim", what, w.Key)
			}
		}
	}
}

func (l *linter) predictions(n *content.Node, file string) {
	var preds []content.Prediction
	if n.Type == "prediction" {
		p := n.AsPrediction()
		if p.Theory == "" {
			l.add(Error, "P5", n.ID, file, 0, "prediction node needs theory:")
		} else if _, ok := l.c.Nodes[p.Theory]; !ok {
			l.add(Warning, "S5", n.ID, file, 0, "prediction theory %s not written yet", p.Theory)
		}
		preds = append(preds, p)
	}
	if len(n.Predictions) > 0 && n.Type != "theory" {
		l.add(Error, "P5", n.ID, file, 0, "inline predictions are only allowed on theory nodes")
	}
	preds = append(preds, n.Predictions...)
	for _, p := range preds {
		what := "prediction " + p.ID
		if !strings.HasPrefix(p.ID, "pred.") {
			l.add(Error, "P5", n.ID, file, 0, "%s: id must start with pred.", what)
		}
		if p.Inline {
			if _, clash := l.c.Nodes[p.ID]; clash {
				l.add(Error, "P5", n.ID, file, 0, "%s is declared inline and as a node", what)
			}
		}
		for _, lang := range content.Langs {
			if p.Statement[lang] == "" {
				l.add(Error, "P5", n.ID, file, 0, "%s: statement.%s missing", what, lang)
			}
		}
		if !outcomes[p.Outcome] {
			l.add(Error, "P5", n.ID, file, 0, "%s: outcome must be one of %s", what, keys(outcomes))
		}
		if p.Made == nil || p.Made.Year == 0 || p.Made.In == "" {
			l.add(Error, "P5", n.ID, file, 0, "%s: made {year, by, in} is required", what)
		} else {
			l.workRef(n, file, "L2", p.Made.In, what+" made.in")
		}
		decided := p.Outcome == "confirmed" || p.Outcome == "refuted" || p.Outcome == "partial"
		if decided && (p.Test == nil || p.Test.In == "") {
			l.add(Error, "P5", n.ID, file, 0, "%s: outcome %s needs test.in (a work)", what, p.Outcome)
		}
		if p.Test != nil {
			w := l.workRef(n, file, "L2", p.Test.In, what+" test.in")
			if w != nil && p.Outcome == "refuted" && w.Tier != "T1" {
				l.add(Error, "L3", n.ID, file, 0, "%s: refuted outcome needs a T1 test source, %s is %s", what, w.Key, w.Tier)
			}
			if p.MadeBeforeObservation && p.Made != nil && p.Test.Year != 0 && p.Made.Year > p.Test.Year {
				l.add(Error, "P5", n.ID, file, 0, "%s: made_before_observation but made %d after test %d", what, p.Made.Year, p.Test.Year)
			}
		}
		if p.Retrodiction && p.MadeBeforeObservation {
			l.add(Error, "P5", n.ID, file, 0, "%s: a retrodiction cannot be made before the observation", what)
		}
	}
}

func (l *linter) freshness(n *content.Node, file string) {
	if n.LastVerified == "" {
		l.add(Error, "L7", n.ID, file, 0, "last_verified missing (P10)")
	} else if d, err := time.Parse("2006-01-02", n.LastVerified); err != nil {
		l.add(Error, "L7", n.ID, file, 0, "last_verified %q is not YYYY-MM-DD", n.LastVerified)
	} else if (n.Status == "contested" || n.Status == "speculative") && d.AddDate(1, 0, 0).Before(l.opt.Today) {
		l.add(Warning, "L7", n.ID, file, 0, "%s card last verified %s (>12 months): awaiting review", n.Status, n.LastVerified)
	} else if d.After(l.opt.Today) {
		l.add(Error, "L7", n.ID, file, 0, "last_verified %s is in the future", n.LastVerified)
	}
	switch n.VerifiedBy {
	case "human":
	case "draft", "":
		l.add(Warning, "V2", n.ID, file, 0, "verified_by=%q: card awaits human review", n.VerifiedBy)
	default:
		l.add(Error, "L7", n.ID, file, 0, "verified_by must be human or draft")
	}
	for _, t := range n.TodoVerify {
		l.add(Warning, "V3", n.ID, file, 0, "to verify: %s", t)
	}
}

// citedWorks returns all works cited by the node (bodies, edges, origin, predictions).
func (l *linter) citedWorks(n *content.Node) map[string]bool {
	out := map[string]bool{}
	for _, d := range n.Docs {
		for _, b := range d.Blocks {
			for _, para := range b.Paragraphs() {
				for _, c := range content.Cites(para) {
					out[c.W] = true
				}
			}
		}
	}
	for _, e := range n.Edges {
		if e.Source != nil {
			out[content.WorkKey(e.Source.Work)] = true
		}
	}
	if n.FirstProposed != nil && n.FirstProposed.In != "" {
		out[content.WorkKey(n.FirstProposed.In)] = true
	}
	return out
}

func (l *linter) tierOf(ref string) string {
	if w, ok := l.c.Works[content.WorkKey(ref)]; ok {
		return w.Tier
	}
	return ""
}

func (l *linter) yearOf(ref string) int {
	if w, ok := l.c.Works[content.WorkKey(ref)]; ok {
		return w.Year
	}
	return 0
}

// theoryPredictions returns inline predictions plus prediction nodes pointing at theory id.
func (l *linter) theoryPredictions(id string) []content.Prediction {
	var out []content.Prediction
	for _, p := range l.c.AllPredictions() {
		if p.Theory == id {
			out = append(out, p)
		}
	}
	return out
}

// status applies the mechanical rules of P3 (L3).
func (l *linter) status(n *content.Node, file string) {
	if n.Status == "" {
		return
	}
	var t1Support, t1Contra, t1ContraRecent int
	for _, e := range n.Edges {
		if e.Source == nil || l.tierOf(e.Source.Work) != "T1" {
			continue
		}
		switch e.Rel {
		case "supported_by":
			t1Support++
		case "contradicted_by", "opposed_by":
			t1Contra++
			if l.yearOf(e.Source.Work) >= l.opt.Today.Year()-10 {
				t1ContraRecent++
			}
		}
	}
	preds := l.theoryPredictions(n.ID)
	if n.Type == "prediction" {
		preds = []content.Prediction{n.AsPrediction()}
	}
	refutedPred, decidedPred := false, false
	for _, p := range preds {
		if p.Outcome == "refuted" {
			refutedPred = true
		}
		if p.Outcome == "confirmed" || p.Outcome == "refuted" {
			decidedPred = true
		}
	}
	t1Cited := 0
	for k := range l.citedWorks(n) {
		if l.tierOf(k) == "T1" {
			t1Cited++
		}
	}
	switch n.Status {
	case "established":
		if t1Cited < 2 {
			l.add(Error, "L3", n.ID, file, 0, "established needs ≥ 2 independent T1 sources, has %d", t1Cited)
		}
		if t1ContraRecent > 0 {
			l.add(Error, "L3", n.ID, file, 0, "established but has %d T1 contradicted_by/opposed_by from the last 10 years", t1ContraRecent)
		}
	case "strong":
		if t1Cited < 1 {
			l.add(Error, "L3", n.ID, file, 0, "strong needs ≥ 1 T1 source")
		}
	case "contested":
		if t1Support == 0 || t1Contra == 0 {
			l.add(Error, "L3", n.ID, file, 0, "contested needs T1 sources on both sides (supported_by: %d, contradicted_by/opposed_by: %d)", t1Support, t1Contra)
		}
	case "speculative":
		if decidedPred {
			l.add(Error, "L3", n.ID, file, 0, "speculative but a prediction already has a decisive outcome")
		}
	case "refuted":
		if !refutedPred && t1Contra == 0 {
			l.add(Error, "L3", n.ID, file, 0, "refuted needs a refuted prediction or a T1 contradicted_by edge")
		}
	}
	if refutedPred && n.Status != "refuted" && n.Type == "theory" {
		l.add(Warning, "L3", n.ID, file, 0, "has a refuted prediction but status is %s; justify in open_problems", n.Status)
	}
}

var forbidden = []struct{ phrase, lang string }{
	{"vědci dokázali", "cs"}, {"věda dokázala", "cs"}, {"je jisté", "cs"}, {"je dokázáno", "cs"},
	{"mainstream tají", "cs"}, {"vědci tají", "cs"}, {"vědci se mýlili", "cs"},
	{"scientists have proven", "en"}, {"scientists proved", "en"}, {"science has proven", "en"},
	{"it is certain", "en"}, {"mainstream hides", "en"}, {"scientists were wrong", "en"},
}

var attributionOnly = []string{"pavěd", "pseudověd", "pseudoscien"}

// phrases implements L10 on a piece of text.
func (l *linter) phrases(n *content.Node, file string, line int, text string) {
	low := strings.ToLower(text)
	for _, f := range forbidden {
		if strings.Contains(low, f.phrase) {
			l.add(Error, "L10", n.ID, file, line, "forbidden phrase %q (P9)", f.phrase)
		}
	}
	for _, f := range attributionOnly {
		if strings.Contains(low, f) && !strings.Contains(text, "[@") {
			l.add(Error, "L10", n.ID, file, line, "%q only with attribution (a citation in the same paragraph)", f)
		}
	}
}

func (l *linter) docs(n *content.Node) {
	philosophical := n.Status == "philosophical" || n.Testability == "philosophical"
	for _, key := range sortedDocKeys(n.Docs) {
		d := n.Docs[key]
		quotesPerWork := map[string]int{}
		for _, b := range d.Blocks {
			if b.HasClass("nocite") && !nociteAllowed[n.Type][b.Section] && b.T == "h" {
				l.add(Error, "L1", n.ID, d.File, b.Line, "section %q may not be .nocite for %s", b.Section, n.Type)
			}
			nocite := philosophical || (b.HasClass("nocite") && nociteAllowed[n.Type][b.Section])
			for _, para := range b.Paragraphs() {
				line := b.Line
				if len(para) > 0 {
					line = para[0].Line
				}
				cites := content.Cites(para)
				if len(cites) == 0 && !nocite {
					l.add(Error, "L1", n.ID, d.File, line, "paragraph without a citation (P7)")
				}
				for _, c := range cites {
					w, ok := l.c.Works[c.W]
					if !ok {
						l.add(Error, "L2", n.ID, d.File, line, "unknown work [@%s]", c.W)
						continue
					}
					if c.Loc == "" {
						l.add(Error, "L1", n.ID, d.File, line, "[@%s] has no locator (P7)", c.W)
					}
					if w.Retracted {
						l.add(Error, "L2", n.ID, d.File, line, "[@%s] is retracted", c.W)
					}
				}
				if len(cites) > 0 && !nocite && b.T != "q" && allTier(l, cites, "T6") {
					l.add(Error, "P4", n.ID, d.File, line, "paragraph rests only on T6 sources")
				}
				for _, sp := range para {
					if (sp.T == "link" || sp.T == "term") && l.c.Nodes[sp.To] == nil {
						l.add(Warning, "S5", n.ID, d.File, line, "link to %s: node not written yet", sp.To)
					}
				}
				l.phrases(n, d.File, line, paraText(para))
			}
			if b.T == "q" {
				l.quote(n, d, b, quotesPerWork)
			}
		}
		l.sections(n, d)
	}
}

func allTier(l *linter, cites []content.Span, tier string) bool {
	for _, c := range cites {
		if l.tierOf(c.W) != tier {
			return false
		}
	}
	return true
}

func paraText(spans []content.Span) string {
	var b strings.Builder
	for _, s := range spans {
		switch s.T {
		case "text", "link", "term":
			b.WriteString(s.V)
		case "cite":
			b.WriteString("[@" + s.W + "]")
		}
	}
	return b.String()
}

// quote implements L9 (P11).
func (l *linter) quote(n *content.Node, d *content.Doc, b content.Block, perWork map[string]int) {
	cites := content.Cites(b.S)
	if len(cites) != 1 {
		l.add(Error, "L9", n.ID, d.File, b.Line, "a quote needs exactly one citation, has %d", len(cites))
		return
	}
	words := len(strings.Fields(content.PlainText(b.S, nil)))
	if words > 20 {
		l.add(Error, "L9", n.ID, d.File, b.Line, "quote has %d words (max 20)", words)
	}
	perWork[cites[0].W]++
	if perWork[cites[0].W] > 1 {
		l.add(Error, "L9", n.ID, d.File, b.Line, "second quote from %s on one card", cites[0].W)
	}
}

// sections implements L4 and the per-type required sections.
func (l *linter) sections(n *content.Node, d *content.Doc) {
	have := set(d.Sections()...)
	for _, s := range requiredSections[n.Type] {
		if !have[s] {
			l.add(Error, "L4", n.ID, d.File, 0, "required section {#%s} missing", s)
		}
	}
	if n.Type != "theory" || d.Level != "lay" {
		return
	}
	items := 0
	for _, b := range d.Blocks {
		if b.Section == "open_problems" && b.T == "ul" {
			items += len(b.Items)
		}
	}
	if items < 2 {
		l.add(Error, "L4", n.ID, d.File, 0, "open_problems needs ≥ 2 list items, has %d", items)
	}
	if d.Lang != "cs" {
		return
	}
	file := n.Dir + "/node.yaml"
	if c := len(l.theoryPredictions(n.ID)); c < 3 {
		l.add(Error, "L4", n.ID, file, 0, "theory needs ≥ 3 predictions, has %d", c)
	}
	competitors := 0
	for _, e := range n.Edges {
		if e.Rel == "competes_with" {
			competitors++
		}
	}
	if competitors == 0 {
		l.add(Error, "L4", n.ID, file, 0, "theory needs ≥ 1 competes_with edge")
	}
}

// graph implements L8: no orphans, and every claim/theory reaches a first_proposed.
func (l *linter) graph() {
	degree := map[string]int{}
	adj := map[string][]string{}
	link := func(a, b string) {
		degree[a]++
		degree[b]++
		adj[a] = append(adj[a], b)
	}
	for _, id := range l.c.NodeIDs() {
		n := l.c.Nodes[id]
		for _, e := range n.Edges {
			link(id, e.To)
		}
		if n.Type == "prediction" && n.Theory != "" {
			link(id, n.Theory)
		}
		for _, d := range n.Docs {
			for _, b := range d.Blocks {
				for _, para := range b.Paragraphs() {
					for _, sp := range para {
						if sp.T == "link" || sp.T == "term" {
							degree[id]++
							degree[sp.To]++
						}
					}
				}
			}
		}
	}
	for _, id := range l.c.NodeIDs() {
		n := l.c.Nodes[id]
		file := n.Dir + "/node.yaml"
		if degree[id] == 0 && len(l.c.Nodes) > 1 {
			l.add(Error, "L8", id, file, 0, "orphan node: no edges or links in or out")
		}
		if n.Type != "claim" && n.Type != "theory" {
			continue
		}
		if !l.reachesOrigin(id, adj, map[string]bool{}) {
			l.add(Error, "L8", id, file, 0, "no path to a node with first_proposed (via builds_on/interprets/inspired_by/…)")
		}
	}
}

func (l *linter) reachesOrigin(id string, adj map[string][]string, seen map[string]bool) bool {
	if seen[id] {
		return false
	}
	seen[id] = true
	n, ok := l.c.Nodes[id]
	if !ok {
		return false
	}
	if n.FirstProposed != nil && n.FirstProposed.In != "" {
		return true
	}
	for _, e := range n.Edges {
		switch e.Rel {
		case "builds_on", "interprets", "inspired_by", "first_proposed_in", "superseded_by":
			if l.reachesOrigin(e.To, adj, seen) {
				return true
			}
		}
	}
	return false
}

// MediaLicenses is the allowlist from docs/SOURCES.md (L5).
var MediaLicenses = set("CC0-1.0", "PD-NASA", "PD-old", "CC-BY-4.0", "CC-BY-SA-3.0-IGO", "CC-BY-SA-4.0", "own")

func (l *linter) media() {
	const file = "media/media.yaml"
	for _, m := range l.c.Media {
		what := "media " + m.SHA256
		if m.Credit == "" || m.SourceURL == "" {
			l.add(Error, "L5", m.Node, file, 0, "%s needs credit and source_url", what)
		}
		if !MediaLicenses[m.License] {
			l.add(Error, "L5", m.Node, file, 0, "%s: license %q not in allowlist", what, m.License)
		}
		if m.AIEdited && (strings.Contains(m.Credit, "ESO") || strings.Contains(m.Credit, "NASA") || strings.Contains(m.Credit, "ESA")) {
			l.add(Error, "L5", m.Node, file, 0, "%s: agency images must not be AI-edited", what)
		}
		if _, ok := l.c.Nodes[m.Node]; !ok {
			l.add(Error, "L5", m.Node, file, 0, "%s belongs to unknown node", what)
		}
	}
}

func keys(m map[string]bool) string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return strings.Join(ks, "|")
}

func sortedDocKeys(m map[string]*content.Doc) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}
