package harvest

import (
	"html"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Record is one bibliographic record as a source reports it.
type Record struct {
	Source     string   `json:"source"` // arxiv | crossref | scix
	DOI        string   `json:"doi,omitempty"`
	Arxiv      string   `json:"arxiv,omitempty"`
	Bibcode    string   `json:"bibcode,omitempty"`
	Title      string   `json:"title"`
	Authors    []Author `json:"authors"`
	Years      []int    `json:"years"` // all plausible publication years (issued, print, online, first version)
	Venue      string   `json:"venue,omitempty"`
	Categories []string `json:"categories,omitempty"`
	Abstract   string   `json:"abstract,omitempty"`
	License    string   `json:"license,omitempty"`
	Updated    string   `json:"updated,omitempty"` // source datestamp, YYYY-MM-DD
	Retracted  bool     `json:"retracted,omitempty"`
}

// Author keeps the family name separately; that is what matching uses.
type Author struct {
	Family string `json:"family"`
	Given  string `json:"given,omitempty"`
}

// Year is the first reported year, or 0.
func (r *Record) Year() int {
	if len(r.Years) == 0 {
		return 0
	}
	return r.Years[0]
}

func addYear(ys []int, y int) []int {
	if y <= 0 {
		return ys
	}
	for _, v := range ys {
		if v == y {
			return ys
		}
	}
	return append(ys, y)
}

// NormDOI lower-cases a DOI and strips resolver prefixes.
func NormDOI(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	for _, p := range []string{"https://doi.org/", "http://doi.org/", "https://dx.doi.org/", "doi:"} {
		s = strings.TrimPrefix(s, p)
	}
	return s
}

var arxivVersion = regexp.MustCompile(`v\d+$`)

// NormArxiv strips the "arXiv:" prefix, abs URL and version suffix.
func NormArxiv(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "http://arxiv.org/abs/")
	s = strings.TrimPrefix(s, "https://arxiv.org/abs/")
	s = strings.TrimPrefix(strings.TrimPrefix(s, "arXiv:"), "arxiv:")
	return arxivVersion.ReplaceAllString(s, "")
}

var (
	// unitPower turns "km s −1 Mpc −1" into "km s Mpc" so it folds like "km/s/Mpc".
	unitPower = regexp.MustCompile(`([A-Za-z])\s*\^?\{?[−\-–]\s*1\}?\b`)
	greek     = strings.NewReplacer(
		"α", " alpha ", "β", " beta ", "γ", " gamma ", "δ", " delta ", "ε", " epsilon ", "ζ", " zeta ",
		"η", " eta ", "θ", " theta ", "κ", " kappa ", "λ", " lambda ", "μ", " mu ", "ν", " nu ",
		"ξ", " xi ", "π", " pi ", "ρ", " rho ", "σ", " sigma ", "τ", " tau ", "φ", " phi ",
		"χ", " chi ", "ψ", " psi ", "ω", " omega ")
	latexCmd = regexp.MustCompile(`\\[a-zA-Z]+`)
	tags     = regexp.MustCompile(`<[^>]+>`)
	spaces   = regexp.MustCompile(`\s+`)
)

// Fold lower-cases, removes diacritics, LaTeX commands, markup and
// punctuation, and collapses whitespace. Used for title/name comparison.
func Fold(s string) string {
	s = html.UnescapeString(tags.ReplaceAllString(s, " "))
	s = unitPower.ReplaceAllString(s, "$1")
	s = greek.Replace(strings.ToLower(s))
	s = latexCmd.ReplaceAllStringFunc(s, func(cmd string) string {
		// \Lambda, \Omega … keep the letter name, drop other commands.
		if name := strings.ToLower(cmd[1:]); isGreekName(name) {
			return " " + name + " "
		}
		return " "
	})
	s = norm.NFD.String(s)
	var b strings.Builder
	for _, r := range s {
		switch {
		case unicode.Is(unicode.Mn, r):
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(unicode.ToLower(r))
		default:
			b.WriteByte(' ')
		}
	}
	return strings.TrimSpace(spaces.ReplaceAllString(b.String(), " "))
}

func isGreekName(s string) bool {
	switch s {
	case "alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", "theta", "kappa", "lambda", "mu", "nu",
		"xi", "pi", "rho", "sigma", "tau", "phi", "chi", "psi", "omega":
		return true
	}
	return false
}

// CleanText strips markup (JATS/HTML), decodes entities and collapses whitespace.
func CleanText(s string) string {
	s = html.UnescapeString(tags.ReplaceAllString(s, " "))
	return strings.TrimSpace(spaces.ReplaceAllString(s, " "))
}
