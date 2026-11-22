package harvest

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Arxiv reads the arXiv query API (by ID) and OAI-PMH (daily harvest).
// arXiv asks for at most one request every 3 seconds on a single connection;
// the caller's Client must be configured accordingly.
type Arxiv struct {
	C       *Client
	APIBase string // default https://export.arxiv.org/api/query
	OAIBase string // default https://oaipmh.arxiv.org/oai
}

// DefaultCategories are the arXiv categories KosmoMat follows (plan §8).
// A prefix ending in "." matches all its subcategories.
var DefaultCategories = []string{"astro-ph", "astro-ph.", "physics.plasm-ph", "physics.space-ph", "gr-qc", "physics.hist-ph"}

// DefaultSets are the OAI-PMH sets that contain DefaultCategories.
var DefaultSets = []string{"physics:astro-ph", "physics:gr-qc", "physics:physics"}

// KeepCategories reports whether any of cats is in the allowlist.
func KeepCategories(cats, allow []string) bool {
	for _, c := range cats {
		for _, a := range allow {
			if c == a || (strings.HasSuffix(a, ".") && strings.HasPrefix(c, a)) {
				return true
			}
		}
	}
	return false
}

type atomFeed struct {
	Entries []struct {
		ID        string `xml:"id"`
		Title     string `xml:"title"`
		Summary   string `xml:"summary"`
		Published string `xml:"published"`
		Updated   string `xml:"updated"`
		Authors   []struct {
			Name string `xml:"name"`
		} `xml:"author"`
		Categories []struct {
			Term string `xml:"term,attr"`
		} `xml:"category"`
		DOI        string `xml:"http://arxiv.org/schemas/atom doi"`
		JournalRef string `xml:"http://arxiv.org/schemas/atom journal_ref"`
	} `xml:"entry"`
}

// splitName turns "Benoit Famaey" into family "Famaey", given "Benoit".
// Particles (van, de, …) stay with the family name.
func splitName(full string) Author {
	parts := strings.Fields(full)
	if len(parts) < 2 {
		return Author{Family: full}
	}
	i := len(parts) - 1
	for i > 1 && isParticle(parts[i-1]) {
		i--
	}
	return Author{Family: strings.Join(parts[i:], " "), Given: strings.Join(parts[:i], " ")}
}

func isParticle(s string) bool {
	switch strings.ToLower(s) {
	case "van", "von", "de", "der", "den", "da", "di", "du", "del", "della", "la", "le", "ter", "dos":
		return true
	}
	return false
}

func year(ts string) int {
	if len(ts) < 4 {
		return 0
	}
	t, err := time.Parse("2006", ts[:4])
	if err != nil {
		return 0
	}
	return t.Year()
}

func (a *Arxiv) apiBase() string {
	if a.APIBase != "" {
		return a.APIBase
	}
	return "https://export.arxiv.org/api/query"
}

func (a *Arxiv) oaiBase() string {
	if a.OAIBase != "" {
		return a.OAIBase
	}
	return "https://oaipmh.arxiv.org/oai"
}

// ByIDs fetches records for arXiv IDs (batched by 50). Missing IDs are absent from the map.
func (a *Arxiv) ByIDs(ctx context.Context, ids []string) (map[string]Record, error) {
	out := map[string]Record{}
	for start := 0; start < len(ids); start += 50 {
		batch := ids[start:min(start+50, len(ids))]
		norm := make([]string, len(batch))
		for i, id := range batch {
			norm[i] = NormArxiv(id)
		}
		q := url.Values{}
		q.Set("id_list", strings.Join(norm, ","))
		q.Set("max_results", fmt.Sprint(len(norm)))
		body, err := a.C.Get(ctx, a.apiBase()+"?"+q.Encode())
		if err != nil {
			return out, err
		}
		var f atomFeed
		if err := xml.Unmarshal(body, &f); err != nil {
			return out, fmt.Errorf("arxiv atom: %w", err)
		}
		for _, e := range f.Entries {
			id := NormArxiv(e.ID)
			if id == "" || strings.Contains(e.ID, "api/errors") {
				continue
			}
			r := Record{Source: "arxiv", Arxiv: id, DOI: NormDOI(e.DOI), Title: CleanText(e.Title),
				Abstract: CleanText(e.Summary), Venue: e.JournalRef}
			r.Years = addYear(r.Years, year(e.Published))
			for _, au := range e.Authors {
				r.Authors = append(r.Authors, splitName(CleanText(au.Name)))
			}
			for _, c := range e.Categories {
				r.Categories = append(r.Categories, c.Term)
			}
			if len(e.Updated) >= 10 {
				r.Updated = e.Updated[:10]
			}
			out[id] = r
		}
	}
	return out, nil
}

type oaiResponse struct {
	Error *struct {
		Code string `xml:"code,attr"`
		Msg  string `xml:",chardata"`
	} `xml:"error"`
	ListRecords struct {
		Records []struct {
			Header struct {
				Status    string `xml:"status,attr"`
				Datestamp string `xml:"datestamp"`
			} `xml:"header"`
			Meta struct {
				ID      string `xml:"id"`
				Created string `xml:"created"`
				Authors []struct {
					Keyname   string `xml:"keyname"`
					Forenames string `xml:"forenames"`
				} `xml:"authors>author"`
				Title      string `xml:"title"`
				Categories string `xml:"categories"`
				JournalRef string `xml:"journal-ref"`
				DOI        string `xml:"doi"`
				License    string `xml:"license"`
				Abstract   string `xml:"abstract"`
			} `xml:"metadata>arXiv"`
		} `xml:"record"`
		Token string `xml:"resumptionToken"`
	} `xml:"ListRecords"`
}

// Harvest lists records of one OAI set changed between from and until (YYYY-MM-DD),
// keeping only records with an allowed category. "noRecordsMatch" is not an error.
func (a *Arxiv) Harvest(ctx context.Context, set, from, until string, allow []string, fn func(Record) error) error {
	q := url.Values{}
	q.Set("verb", "ListRecords")
	q.Set("metadataPrefix", "arXiv")
	q.Set("set", set)
	q.Set("from", from)
	q.Set("until", until)
	next := a.oaiBase() + "?" + q.Encode()
	for next != "" {
		body, err := a.C.Get(ctx, next)
		if err != nil {
			return err
		}
		var resp oaiResponse
		if err := xml.Unmarshal(body, &resp); err != nil {
			return fmt.Errorf("arxiv oai: %w", err)
		}
		if resp.Error != nil {
			if resp.Error.Code == "noRecordsMatch" {
				return nil
			}
			return fmt.Errorf("arxiv oai %s: %s", resp.Error.Code, strings.TrimSpace(resp.Error.Msg))
		}
		for _, rec := range resp.ListRecords.Records {
			if rec.Header.Status == "deleted" {
				continue
			}
			m := rec.Meta
			cats := strings.Fields(m.Categories)
			if !KeepCategories(cats, allow) {
				continue
			}
			r := Record{Source: "arxiv", Arxiv: NormArxiv(m.ID), DOI: NormDOI(m.DOI), Title: CleanText(m.Title),
				Abstract: CleanText(m.Abstract), Venue: CleanText(m.JournalRef), Categories: cats,
				License: m.License, Updated: rec.Header.Datestamp}
			r.Years = addYear(r.Years, year(m.Created))
			for _, au := range m.Authors {
				r.Authors = append(r.Authors, Author{Family: CleanText(au.Keyname), Given: CleanText(au.Forenames)})
			}
			if err := fn(r); err != nil {
				return err
			}
		}
		next = ""
		if t := strings.TrimSpace(resp.ListRecords.Token); t != "" {
			next = a.oaiBase() + "?verb=ListRecords&resumptionToken=" + url.QueryEscape(t)
		}
	}
	return nil
}
