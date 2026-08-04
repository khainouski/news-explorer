// Package feed adapts external RSS/Atom sources into domain-shaped data - the HTTP-fetching
// counterpart to adapter/postgres, used by usecase/sync to pull live articles instead of relying
// on migration-seeded data.
package feed

import (
	"errors"
	"net/http"
	"time"
)

const (
	fetchTimeout = 10 * time.Second
	maxRedirects = 5
)

// Client fetches and parses one source's feed at a time - stateless beyond its http.Client, safe
// for concurrent use across sources.
type Client struct {
	httpClient *http.Client
}

// New builds a Client with a bounded per-request timeout and a bounded number of redirects, so
// one slow, unreachable, or redirect-looping feed can't hang a sync run indefinitely.
func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout:       fetchTimeout,
			CheckRedirect: limitRedirects,
		},
	}
}

func limitRedirects(_ *http.Request, via []*http.Request) error {
	if len(via) >= maxRedirects {
		return errors.New("stopped after too many redirects")
	}

	return nil
}
