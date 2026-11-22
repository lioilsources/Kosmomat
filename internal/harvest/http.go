// Package harvest fetches bibliographic metadata (arXiv, Crossref, SciX),
// stores it in a local works_raw table and checks content/works against it.
package harvest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// UserAgent is deliberately generic; a contact address is added only when
// the operator passes one (Crossref polite pool).
const UserAgent = "metadata-harvester/0.1 (Go)"

// Client is a polite HTTP client: one request at a time, a minimum gap
// between requests, and retries on 429/503 honouring Retry-After.
type Client struct {
	HTTP     *http.Client
	MinGap   time.Duration
	Retries  int
	Mailto   string            // optional, appended to the User-Agent
	Header   map[string]string // extra headers (e.g. Authorization)
	Sleep    func(time.Duration)
	mu       sync.Mutex
	lastCall time.Time
	now      func() time.Time
}

// NewClient returns a client with the given minimum gap between requests.
func NewClient(gap time.Duration) *Client {
	return &Client{
		HTTP:    &http.Client{Timeout: 60 * time.Second},
		MinGap:  gap,
		Retries: 3,
		Sleep:   time.Sleep,
		now:     time.Now,
	}
}

// Get fetches url and returns the body. 404 is reported as ErrNotFound.
func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var lastErr error
	for attempt := 0; attempt <= c.Retries; attempt++ {
		if wait := c.MinGap - c.now().Sub(c.lastCall); wait > 0 && !c.lastCall.IsZero() {
			c.Sleep(wait)
		}
		c.lastCall = c.now()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		ua := UserAgent
		if c.Mailto != "" {
			ua += " mailto:" + c.Mailto
		}
		req.Header.Set("User-Agent", ua)
		for k, v := range c.Header {
			req.Header.Set(k, v)
		}
		resp, err := c.HTTP.Do(req)
		if err != nil {
			lastErr = err
			c.Sleep(backoff(attempt, ""))
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		switch {
		case err != nil:
			lastErr = err
		case resp.StatusCode == http.StatusOK:
			return body, nil
		case resp.StatusCode == http.StatusNotFound:
			return nil, ErrNotFound
		case resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500:
			lastErr = fmt.Errorf("GET %s: %s", url, resp.Status)
			c.Sleep(backoff(attempt, resp.Header.Get("Retry-After")))
			continue
		default:
			return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
		}
	}
	return nil, lastErr
}

func backoff(attempt int, retryAfter string) time.Duration {
	if s, err := strconv.Atoi(retryAfter); err == nil && s > 0 {
		return min(time.Duration(s)*time.Second, 2*time.Minute)
	}
	return time.Duration(5<<attempt) * time.Second
}

// ErrNotFound is returned when the source has no record for an identifier.
var ErrNotFound = fmt.Errorf("not found")
