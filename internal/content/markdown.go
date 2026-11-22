package content

import (
	"regexp"
	"strings"
)

// Span is an inline piece of a block.
type Span struct {
	T    string `json:"t"`             // text|cite|link|term
	V    string `json:"v,omitempty"`   // text or link label
	B    bool   `json:"b,omitempty"`   // bold
	I    bool   `json:"i,omitempty"`   // italic
	W    string `json:"w,omitempty"`   // cite: work key
	Loc  string `json:"loc,omitempty"` // cite: locator
	To   string `json:"to,omitempty"`  // link/term: node id
	Line int    `json:"-"`
}

// Block is a block-level element of a card body.
type Block struct {
	T     string   `json:"t"`               // h|p|ul|q
	ID    string   `json:"id,omitempty"`    // h: section id
	Text  string   `json:"text,omitempty"`  // h: heading text
	S     []Span   `json:"s,omitempty"`     // p, q
	Items [][]Span `json:"items,omitempty"` // ul

	Line    int      `json:"-"`
	Section string   `json:"-"` // id of the enclosing section
	Classes []string `json:"-"` // classes of the enclosing section (e.g. "nocite")
}

// Doc is one parsed Markdown file (one language and level of a card).
type Doc struct {
	File   string
	Lang   string
	Level  string
	Blocks []Block
}

// Sections returns the section ids present in the doc, in order.
func (d *Doc) Sections() []string {
	var out []string
	for _, b := range d.Blocks {
		if b.T == "h" && b.ID != "" {
			out = append(out, b.ID)
		}
	}
	return out
}

// HasClass reports whether the block's section carries the class.
func (b *Block) HasClass(c string) bool {
	for _, x := range b.Classes {
		if x == c {
			return true
		}
	}
	return false
}

var (
	reHeading = regexp.MustCompile(`^(#{1,6})\s+(.*?)\s*(?:\{([^}]*)\})?\s*$`)
	reList    = regexp.MustCompile(`^\s*[-*]\s+(.*)$`)
)

// ParseMarkdown converts the KosmoMat Markdown dialect into blocks.
func ParseMarkdown(src string) []Block {
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")
	var (
		blocks  []Block
		section string
		classes []string
		para    []string
		paraAt  int
		kind    string // "", "p", "ul", "q"
		items   []string
		itemAt  []int
	)
	flush := func() {
		switch kind {
		case "p":
			blocks = append(blocks, Block{T: "p", S: ParseInline(strings.Join(para, " "), paraAt), Line: paraAt, Section: section, Classes: classes})
		case "q":
			blocks = append(blocks, Block{T: "q", S: ParseInline(strings.Join(para, " "), paraAt), Line: paraAt, Section: section, Classes: classes})
		case "ul":
			b := Block{T: "ul", Line: itemAt[0], Section: section, Classes: classes}
			for i, it := range items {
				b.Items = append(b.Items, ParseInline(it, itemAt[i]))
			}
			blocks = append(blocks, b)
		}
		kind, para, items, itemAt = "", nil, nil, nil
	}
	for i, raw := range lines {
		ln := i + 1
		line := strings.TrimRight(raw, " \t")
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			flush()
		case strings.HasPrefix(trimmed, "<!--"):
			// HTML comments are editorial notes, not content.
		case reHeading.MatchString(trimmed) && strings.HasPrefix(trimmed, "#"):
			flush()
			m := reHeading.FindStringSubmatch(trimmed)
			id, cls := parseAttrs(m[3])
			section, classes = id, cls
			blocks = append(blocks, Block{T: "h", ID: id, Text: m[2], Line: ln, Section: id, Classes: cls})
		case strings.HasPrefix(trimmed, ">"):
			if kind != "q" {
				flush()
				kind, paraAt = "q", ln
			}
			para = append(para, strings.TrimSpace(strings.TrimPrefix(trimmed, ">")))
		case reList.MatchString(line):
			if kind != "ul" {
				flush()
				kind = "ul"
			}
			items = append(items, reList.FindStringSubmatch(line)[1])
			itemAt = append(itemAt, ln)
		default:
			if kind == "ul" && strings.HasPrefix(line, " ") {
				items[len(items)-1] += " " + trimmed
				continue
			}
			if kind != "p" {
				flush()
				kind, paraAt = "p", ln
			}
			para = append(para, trimmed)
		}
	}
	flush()
	return blocks
}

func parseAttrs(s string) (id string, classes []string) {
	for _, f := range strings.Fields(s) {
		switch {
		case strings.HasPrefix(f, "#"):
			id = f[1:]
		case strings.HasPrefix(f, "."):
			classes = append(classes, f[1:])
		}
	}
	return id, classes
}

// ParseInline splits a paragraph into spans: text, citations [@key loc],
// node links [[id]] / [[id|label]], **bold** and *italic*.
func ParseInline(s string, line int) []Span {
	var (
		out        []Span
		buf        strings.Builder
		bold, ital bool
	)
	emit := func() {
		if buf.Len() > 0 {
			out = append(out, Span{T: "text", V: buf.String(), B: bold, I: ital, Line: line})
			buf.Reset()
		}
	}
	for i := 0; i < len(s); {
		switch {
		case strings.HasPrefix(s[i:], "[@"):
			end := strings.IndexByte(s[i:], ']')
			if end < 0 {
				buf.WriteString(s[i:])
				i = len(s)
				continue
			}
			emit()
			// Drop the space before a citation so it renders as a superscript.
			if n := len(out); n > 0 && out[n-1].T == "text" {
				out[n-1].V = strings.TrimRight(out[n-1].V, " ")
			}
			for _, part := range strings.Split(s[i+1:i+end], ";") {
				part = strings.TrimSpace(part)
				part = strings.TrimPrefix(part, "@")
				key, loc, _ := strings.Cut(part, " ")
				out = append(out, Span{T: "cite", W: WorkKey(strings.TrimSpace(key)), Loc: strings.TrimSpace(loc), Line: line})
			}
			i += end + 1
		case strings.HasPrefix(s[i:], "[["):
			end := strings.Index(s[i:], "]]")
			if end < 0 {
				buf.WriteString(s[i:])
				i = len(s)
				continue
			}
			emit()
			inner := s[i+2 : i+end]
			to, label, _ := strings.Cut(inner, "|")
			to = strings.TrimSpace(to)
			t := "link"
			if strings.HasPrefix(to, "term.") {
				t = "term"
			}
			out = append(out, Span{T: t, To: to, V: strings.TrimSpace(label), B: bold, I: ital, Line: line})
			i += end + 2
		case strings.HasPrefix(s[i:], "**"):
			emit()
			bold = !bold
			i += 2
		case s[i] == '*' && (i+1 < len(s) && s[i+1] != ' ' || ital):
			emit()
			ital = !ital
			i++
		default:
			buf.WriteByte(s[i])
			i++
		}
	}
	emit()
	return out
}

// PlainText flattens spans into text (citations dropped, links by label or id).
func PlainText(spans []Span, titles func(id string) string) string {
	var b strings.Builder
	for _, sp := range spans {
		switch sp.T {
		case "text":
			b.WriteString(sp.V)
		case "link", "term":
			if sp.V != "" {
				b.WriteString(sp.V)
			} else if titles != nil {
				b.WriteString(titles(sp.To))
			}
		}
	}
	return b.String()
}

// Cites returns the citation spans of a span list.
func Cites(spans []Span) []Span {
	var out []Span
	for _, sp := range spans {
		if sp.T == "cite" {
			out = append(out, sp)
		}
	}
	return out
}

// Paragraphs returns every citable unit of a block: the paragraph itself,
// or each list item separately.
func (b *Block) Paragraphs() [][]Span {
	switch b.T {
	case "p", "q":
		return [][]Span{b.S}
	case "ul":
		return b.Items
	}
	return nil
}
