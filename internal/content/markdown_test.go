package content

import "testing"

func TestParseInline(t *testing.T) {
	spans := ParseInline("Plochá křivka [@rubin-1980 obr. 4; @work.milgrom-1983 §2] a [[term.sigma|5σ]] i [[obs.cmb]] **tučně**.", 3)
	var cites, terms, links, bold int
	for _, s := range spans {
		switch s.T {
		case "cite":
			cites++
			if s.Loc == "" {
				t.Errorf("cite %s without locator", s.W)
			}
		case "term":
			terms++
			if s.To != "term.sigma" || s.V != "5σ" {
				t.Errorf("term span = %+v", s)
			}
		case "link":
			links++
		case "text":
			if s.B {
				bold++
			}
		}
	}
	if cites != 2 || terms != 1 || links != 1 || bold != 1 {
		t.Fatalf("cites=%d terms=%d links=%d bold=%d: %+v", cites, terms, links, bold, spans)
	}
	if spans[1].W != "rubin-1980" || spans[2].W != "milgrom-1983" || spans[2].Loc != "§2" {
		t.Fatalf("cite parse: %+v %+v", spans[1], spans[2])
	}
	if spans[0].V != "Plochá křivka" {
		t.Fatalf("space before citation not trimmed: %q", spans[0].V)
	}
}

func TestParseMarkdownSections(t *testing.T) {
	src := "## Kroky {#steps .nocite}\n\nUdělej tohle.\n\n- jedna [@a s. 1]\n- dva\n  pokračuje\n\n> citát [@b s. 2]\n\n## Dál {#next}\n\nText\nna dva řádky.\n"
	bl := ParseMarkdown(src)
	if len(bl) != 6 {
		t.Fatalf("want 6 blocks, got %d: %+v", len(bl), bl)
	}
	if bl[0].ID != "steps" || !bl[0].HasClass("nocite") {
		t.Fatalf("heading attrs: %+v", bl[0])
	}
	if bl[2].T != "ul" || len(bl[2].Items) != 2 || bl[2].Section != "steps" {
		t.Fatalf("list: %+v", bl[2])
	}
	if got := PlainText(bl[2].Items[1], nil); got != "dva pokračuje" {
		t.Fatalf("continuation: %q", got)
	}
	if bl[3].T != "q" || bl[5].Section != "next" || PlainText(bl[5].S, nil) != "Text na dva řádky." {
		t.Fatalf("quote/para: %+v %+v", bl[3], bl[5])
	}
}
