package harvest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Scix reads the SciX / NASA ADS search API. It needs a token
// (env ADS_API_TOKEN); without one the caller should skip it.
type Scix struct {
	C    *Client // must carry "Authorization: Bearer <token>"
	Base string  // default https://api.adsabs.harvard.edu/v1
}

const scixFields = "bibcode,title,author,year,doi,identifier,pub,abstract,property"

type scixDoc struct {
	Bibcode    string   `json:"bibcode"`
	Title      []string `json:"title"`
	Author     []string `json:"author"`
	Year       string   `json:"year"`
	DOI        []string `json:"doi"`
	Identifier []string `json:"identifier"`
	Pub        string   `json:"pub"`
	Abstract   string   `json:"abstract"`
	Property   []string `json:"property"`
}

func (d *scixDoc) record() Record {
	r := Record{Source: "scix", Bibcode: d.Bibcode, Venue: d.Pub, Abstract: CleanText(d.Abstract)}
	if len(d.Title) > 0 {
		r.Title = CleanText(d.Title[0])
	}
	if len(d.DOI) > 0 {
		r.DOI = NormDOI(d.DOI[0])
	}
	for _, id := range d.Identifier {
		if strings.HasPrefix(strings.ToLower(id), "arxiv:") {
			r.Arxiv = NormArxiv(id)
		}
	}
	for _, a := range d.Author {
		fam, given, _ := strings.Cut(a, ",")
		r.Authors = append(r.Authors, Author{Family: strings.TrimSpace(fam), Given: strings.TrimSpace(given)})
	}
	if y, err := strconv.Atoi(d.Year); err == nil {
		r.Years = addYear(r.Years, y)
	}
	for _, p := range d.Property {
		if p == "RETRACTION" || p == "RETRACTED" {
			r.Retracted = true
		}
	}
	return r
}

func (s *Scix) base() string {
	if s.Base != "" {
		return strings.TrimRight(s.Base, "/")
	}
	return "https://api.adsabs.harvard.edu/v1"
}

// Query runs one search and pages through all results.
func (s *Scix) Query(ctx context.Context, q string, fn func(Record) error) error {
	for start := 0; ; {
		v := url.Values{}
		v.Set("q", q)
		v.Set("fl", scixFields)
		v.Set("rows", "200")
		v.Set("start", strconv.Itoa(start))
		body, err := s.C.Get(ctx, s.base()+"/search/query?"+v.Encode())
		if err != nil {
			return err
		}
		var env struct {
			Response struct {
				NumFound int       `json:"numFound"`
				Docs     []scixDoc `json:"docs"`
			} `json:"response"`
		}
		if err := json.Unmarshal(body, &env); err != nil {
			return fmt.Errorf("scix: %w", err)
		}
		for i := range env.Response.Docs {
			if err := fn(env.Response.Docs[i].record()); err != nil {
				return err
			}
		}
		start += len(env.Response.Docs)
		if len(env.Response.Docs) == 0 || start >= env.Response.NumFound {
			return nil
		}
	}
}

// ByBibcode fetches one bibcode.
func (s *Scix) ByBibcode(ctx context.Context, bibcode string) (*Record, error) {
	var got *Record
	err := s.Query(ctx, fmt.Sprintf("bibcode:%q", bibcode), func(r Record) error {
		if got == nil {
			got = &r
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if got == nil {
		return nil, ErrNotFound
	}
	return got, nil
}

// Refereed lists refereed astronomy records entered between from and until (YYYY-MM-DD).
func (s *Scix) Refereed(ctx context.Context, from, until string, fn func(Record) error) error {
	return s.Query(ctx, fmt.Sprintf("collection:astronomy property:refereed entdate:[%s TO %s]", from, until), fn)
}
