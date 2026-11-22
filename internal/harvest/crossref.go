package harvest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Crossref reads the Crossref REST API.
type Crossref struct {
	C    *Client
	Base string // default https://api.crossref.org
}

type crDate struct {
	DateParts [][]int `json:"date-parts"`
}

func (d *crDate) year() int {
	if d == nil || len(d.DateParts) == 0 || len(d.DateParts[0]) == 0 {
		return 0
	}
	return d.DateParts[0][0]
}

type crWork struct {
	DOI            string   `json:"DOI"`
	Title          []string `json:"title"`
	Subtitle       []string `json:"subtitle"`
	ContainerTitle []string `json:"container-title"`
	Author         []struct {
		Given  string `json:"given"`
		Family string `json:"family"`
		Name   string `json:"name"`
	} `json:"author"`
	Issued          *crDate `json:"issued"`
	PublishedPrint  *crDate `json:"published-print"`
	PublishedOnline *crDate `json:"published-online"`
	Abstract        string  `json:"abstract"`
	Indexed         *struct {
		DateTime string `json:"date-time"`
	} `json:"indexed"`
	License []struct {
		URL string `json:"URL"`
	} `json:"license"`
	UpdatedBy []struct {
		Type string `json:"type"`
	} `json:"updated-by"`
}

func (w *crWork) record() Record {
	r := Record{Source: "crossref", DOI: NormDOI(w.DOI), Abstract: CleanText(w.Abstract)}
	if len(w.Title) > 0 {
		r.Title = CleanText(w.Title[0])
		if len(w.Subtitle) > 0 && w.Subtitle[0] != "" {
			r.Title += ": " + CleanText(w.Subtitle[0])
		}
	}
	if len(w.ContainerTitle) > 0 {
		r.Venue = w.ContainerTitle[0]
	}
	for _, a := range w.Author {
		fam := a.Family
		if fam == "" {
			fam = a.Name // collaborations
		}
		r.Authors = append(r.Authors, Author{Family: fam, Given: a.Given})
	}
	r.Years = addYear(r.Years, w.Issued.year())
	r.Years = addYear(r.Years, w.PublishedPrint.year())
	r.Years = addYear(r.Years, w.PublishedOnline.year())
	if len(w.License) > 0 {
		r.License = w.License[0].URL
	}
	if w.Indexed != nil && len(w.Indexed.DateTime) >= 10 {
		r.Updated = w.Indexed.DateTime[:10]
	}
	for _, u := range w.UpdatedBy {
		if t := strings.ToLower(u.Type); t == "retraction" || t == "withdrawal" || t == "removal" {
			r.Retracted = true
		}
	}
	return r
}

func (c *Crossref) base() string {
	if c.Base != "" {
		return strings.TrimRight(c.Base, "/")
	}
	return "https://api.crossref.org"
}

// Work fetches one DOI.
func (c *Crossref) Work(ctx context.Context, doi string) (*Record, error) {
	body, err := c.C.Get(ctx, c.base()+"/works/"+strings.ReplaceAll(url.PathEscape(NormDOI(doi)), "%2F", "/"))
	if err != nil {
		return nil, err
	}
	var env struct {
		Message crWork `json:"message"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("crossref %s: %w", doi, err)
	}
	r := env.Message.record()
	return &r, nil
}

// ByISSN lists works of one journal indexed between from and until (YYYY-MM-DD),
// following Crossref deep-paging cursors.
func (c *Crossref) ByISSN(ctx context.Context, issn, from, until string, fn func(Record) error) error {
	cursor := "*"
	for {
		q := url.Values{}
		q.Set("filter", fmt.Sprintf("issn:%s,from-index-date:%s,until-index-date:%s", issn, from, until))
		q.Set("rows", "500")
		q.Set("cursor", cursor)
		body, err := c.C.Get(ctx, c.base()+"/works?"+q.Encode())
		if err != nil {
			return err
		}
		var env struct {
			Message struct {
				NextCursor string   `json:"next-cursor"`
				Items      []crWork `json:"items"`
			} `json:"message"`
		}
		if err := json.Unmarshal(body, &env); err != nil {
			return fmt.Errorf("crossref issn %s: %w", issn, err)
		}
		for i := range env.Message.Items {
			if err := fn(env.Message.Items[i].record()); err != nil {
				return err
			}
		}
		if len(env.Message.Items) == 0 || env.Message.NextCursor == "" || env.Message.NextCursor == cursor {
			return nil
		}
		cursor = env.Message.NextCursor
	}
}
