// Package content loads the KosmoMat content tree (content/) into memory.
// The format is documented in docs/SCHEMA.md.
package content

// I18n is a per-language string map ("cs", "en").
type I18n map[string]string

// I18nList is a per-language list of strings.
type I18nList map[string][]string

// Source points to a work with a locator (page, section, figure…).
type Source struct {
	Work string `yaml:"work" json:"work"`
	Loc  string `yaml:"loc" json:"loc"`
}

// Edge is a typed relation from the owning node to another node.
type Edge struct {
	Rel    string  `yaml:"rel" json:"rel"`
	To     string  `yaml:"to" json:"to"`
	Source *Source `yaml:"source" json:"source,omitempty"`
	Note   I18n    `yaml:"note" json:"note,omitempty"`
}

// Origin describes when and where an idea first appeared.
type Origin struct {
	Year int      `yaml:"year" json:"year"`
	By   []string `yaml:"by" json:"by,omitempty"`
	In   string   `yaml:"in" json:"in,omitempty"`
	Loc  string   `yaml:"loc" json:"loc,omitempty"`
}

// Test describes how a prediction was tested.
type Test struct {
	Year int      `yaml:"year" json:"year"`
	Obs  string   `yaml:"obs" json:"obs,omitempty"`
	By   []string `yaml:"by" json:"by,omitempty"`
	In   string   `yaml:"in" json:"in,omitempty"`
	Loc  string   `yaml:"loc" json:"loc,omitempty"`
}

// Prediction is one entry of the prediction book (EDITORIAL P5).
type Prediction struct {
	ID                    string  `yaml:"id" json:"id"`
	Theory                string  `yaml:"theory" json:"theory"`
	Statement             I18n    `yaml:"statement" json:"statement"`
	Made                  *Origin `yaml:"made" json:"made,omitempty"`
	MadeBeforeObservation bool    `yaml:"made_before_observation" json:"made_before_observation"`
	Retrodiction          bool    `yaml:"retrodiction" json:"retrodiction"`
	Test                  *Test   `yaml:"test" json:"test,omitempty"`
	Outcome               string  `yaml:"outcome" json:"outcome"`
	Note                  I18n    `yaml:"note" json:"note,omitempty"`
	// Inline is true for predictions declared inside a theory's node.yaml.
	Inline bool `yaml:"-" json:"inline"`
}

// ChangelogEntry records a public change of the card (EDITORIAL P12).
type ChangelogEntry struct {
	Date string `yaml:"date" json:"date"`
	Kind string `yaml:"kind" json:"kind"`
	Note string `yaml:"note" json:"note"`
	Work string `yaml:"work" json:"work,omitempty"`
}

// Node is the parsed node.yaml.
type Node struct {
	ID            string           `yaml:"id"`
	Type          string           `yaml:"type"`
	Status        string           `yaml:"status"`
	Literature    string           `yaml:"literature"`
	Testability   string           `yaml:"testability"`
	CaseKind      string           `yaml:"case_kind"`
	Difficulty    string           `yaml:"difficulty"`
	Needs         I18nList         `yaml:"needs"`
	Astromat      bool             `yaml:"astromat"`
	Era           string           `yaml:"era"`
	Title         I18n             `yaml:"title"`
	Summary       I18n             `yaml:"summary"`
	FirstProposed *Origin          `yaml:"first_proposed"`
	Edges         []Edge           `yaml:"edges"`
	Predictions   []Prediction     `yaml:"predictions"`
	Changelog     []ChangelogEntry `yaml:"changelog"`
	LastVerified  string           `yaml:"last_verified"`
	VerifiedBy    string           `yaml:"verified_by"`
	TodoVerify    []string         `yaml:"todo_verify"`
	// RequiredSections is informational (plan §4.3 example); the linter derives
	// required sections from the type.
	RequiredSections []string `yaml:"required_sections"`

	// Prediction fields, used when Type == "prediction".
	Theory                string  `yaml:"theory"`
	Statement             I18n    `yaml:"statement"`
	Made                  *Origin `yaml:"made"`
	MadeBeforeObservation bool    `yaml:"made_before_observation"`
	Retrodiction          bool    `yaml:"retrodiction"`
	Test                  *Test   `yaml:"test"`
	Outcome               string  `yaml:"outcome"`
	Note                  I18n    `yaml:"note"`

	// Filled by the loader.
	Dir  string          `yaml:"-"`
	Docs map[string]*Doc `yaml:"-"` // key: "<lang>.<level>", e.g. "cs.lay"
}

// AsPrediction returns the prediction fields of a prediction node.
func (n *Node) AsPrediction() Prediction {
	return Prediction{
		ID: n.ID, Theory: n.Theory, Statement: n.Statement, Made: n.Made,
		MadeBeforeObservation: n.MadeBeforeObservation, Retrodiction: n.Retrodiction,
		Test: n.Test, Outcome: n.Outcome, Note: n.Note,
	}
}

// WorkIDs holds external identifiers of a work.
type WorkIDs struct {
	DOI     string `yaml:"doi" json:"doi,omitempty"`
	Arxiv   string `yaml:"arxiv" json:"arxiv,omitempty"`
	Bibcode string `yaml:"bibcode" json:"bibcode,omitempty"`
	ISBN    string `yaml:"isbn" json:"isbn,omitempty"`
	URL     string `yaml:"url" json:"url,omitempty"`
}

// Work is a bibliographic record (content/works/<year>/<key>.yaml).
type Work struct {
	Key       string   `yaml:"key" json:"key"`
	Tier      string   `yaml:"tier" json:"tier"`
	Kind      string   `yaml:"kind" json:"kind"`
	Title     string   `yaml:"title" json:"title"`
	Authors   []string `yaml:"authors" json:"authors"`
	Year      int      `yaml:"year" json:"year"`
	Venue     string   `yaml:"venue" json:"venue"`
	IDs       WorkIDs  `yaml:"ids" json:"ids"`
	Country   []string `yaml:"country" json:"country,omitempty"`
	Retracted bool     `yaml:"retracted" json:"retracted"`
	Verified  bool     `yaml:"verified" json:"verified"`

	File string `yaml:"-" json:"-"`
}

// Media is one entry of content/media/media.yaml.
type Media struct {
	SHA256    string `yaml:"sha256"`
	Ext       string `yaml:"ext"`
	Node      string `yaml:"node"`
	Credit    string `yaml:"credit"`
	License   string `yaml:"license"`
	SourceURL string `yaml:"source_url"`
	W         int    `yaml:"w"`
	H         int    `yaml:"h"`
	AIEdited  bool   `yaml:"ai_edited"`
}

// Path is a reader's route (content/paths/<slug>.yaml).
type Path struct {
	Slug  string   `yaml:"slug" json:"slug"`
	Title I18n     `yaml:"title" json:"title"`
	Nodes []string `yaml:"nodes" json:"nodes"`
}

// Corpus is the whole loaded content tree.
type Corpus struct {
	Root      string
	Nodes     map[string]*Node
	Works     map[string]*Work
	Media     []Media
	Paths     []Path
	Redirects map[string]string
	// LoadErrors are problems that prevented parsing a file.
	LoadErrors []LoadError
}

// LoadError is a file that could not be parsed.
type LoadError struct {
	File string
	Err  error
}

// WorkKey strips an optional "work." prefix.
func WorkKey(ref string) string {
	if len(ref) > 5 && ref[:5] == "work." {
		return ref[5:]
	}
	return ref
}
