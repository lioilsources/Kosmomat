// Package evalset builds and scores the M2 evaluation set for the verify step
// (plan §9): sentences from the seed cards paired with the passages they cite,
// plus deliberately wrong variants.
package evalset

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/lioilsources/kosmomat/internal/content"
)

// Item is one sentence to verify against the passages of its citations.
// Passages are not stored in the repo (publisher abstracts have their own
// licences); PassageSHA pins the exact text the item was checked against.
type Item struct {
	ID         string    `yaml:"id" json:"id"`
	Node       string    `yaml:"node" json:"node"`
	Sentence   string    `yaml:"sentence" json:"sentence"`
	Cites      []string  `yaml:"cites" json:"cites"`
	PassageSHA string    `yaml:"passage_sha" json:"passage_sha"`
	Expected   string    `yaml:"expected" json:"expected"` // yes | partial | no
	Mutation   *Mutation `yaml:"mutation,omitempty" json:"mutation,omitempty"`
	VerifiedBy string    `yaml:"verified_by" json:"verified_by"` // draft | human
	Note       string    `yaml:"note,omitempty" json:"note,omitempty"`
}

// Mutation describes how a planted item differs from its genuine original.
type Mutation struct {
	Kind string `yaml:"kind" json:"kind"` // number | reversal | wrong-cite
	From string `yaml:"from" json:"from"`
	To   string `yaml:"to" json:"to"`
	Of   string `yaml:"of" json:"of"` // id of the genuine item
}

// Planted reports whether the item is a deliberately wrong variant.
func (it *Item) Planted() bool { return it.Mutation != nil }

// Segment is a run of card text and the citation group that closes it.
type Segment struct {
	Node  string
	Index int
	Text  string
	Cites []content.Span
}

// Segments splits the en.lay text of every node into cited segments. A
// segment is the text since the previous citation group up to the next one.
func Segments(c *content.Corpus, lang string) []Segment {
	var out []Segment
	for _, id := range c.NodeIDs() {
		doc := c.Nodes[id].Docs[lang+".lay"]
		if doc == nil {
			continue
		}
		n := 0
		for _, b := range doc.Blocks {
			for _, para := range b.Paragraphs() {
				var buf strings.Builder
				var group []content.Span
				flush := func() {
					if len(group) > 0 {
						if txt := tidy(buf.String()); txt != "" {
							n++
							out = append(out, Segment{Node: id, Index: n, Text: txt, Cites: group})
						}
						buf.Reset()
						group = nil
					}
				}
				for _, sp := range para {
					if sp.T == "cite" {
						group = append(group, sp)
						continue
					}
					flush()
					buf.WriteString(content.PlainText([]content.Span{sp}, func(id string) string { return c.Title(id, lang) }))
				}
				flush()
			}
		}
	}
	return out
}

var leadPunct = regexp.MustCompile(`^[\s.,;:)\]]+`)

// tidy drops punctuation left over from the previous citation and ends the sentence.
func tidy(s string) string {
	s = strings.TrimSpace(leadPunct.ReplaceAllString(s, ""))
	s = strings.Join(strings.Fields(s), " ")
	if s == "" {
		return ""
	}
	if !strings.ContainsAny(s[len(s)-1:], ".!?") {
		s += "."
	}
	return s
}

// SHA returns the hex sha256 of the joined passages, as pinned in Item.PassageSHA.
func SHA(passages []string) string {
	h := sha256.Sum256([]byte(strings.Join(passages, "\n\n")))
	return hex.EncodeToString(h[:8])
}

// Candidates turns segments whose citations all point at abstracts into
// draft items. abstracts maps work key → abstract text; segments citing a
// work without an available abstract are skipped and listed in missing.
func Candidates(segs []Segment, abstracts map[string]string) (items []Item, missing []string) {
	miss := map[string]bool{}
	for _, s := range segs {
		ok := true
		var keys, passages []string
		for _, ct := range s.Cites {
			if ct.Loc != "abstract" {
				ok = false
				break
			}
			a := abstracts[ct.W]
			if a == "" {
				miss[ct.W] = true
				ok = false
				break
			}
			keys = append(keys, ct.W)
			passages = append(passages, a)
		}
		if !ok {
			continue
		}
		items = append(items, Item{
			ID: fmt.Sprintf("%s#%d", s.Node, s.Index), Node: s.Node, Sentence: s.Text, Cites: keys,
			PassageSHA: SHA(passages), Expected: "yes", VerifiedBy: "draft",
		})
	}
	for k := range miss {
		missing = append(missing, k)
	}
	sort.Strings(missing)
	return items, missing
}
