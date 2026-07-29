// Package feed adapts external RSS/Atom sources into domain-shaped data - the HTTP-fetching
// counterpart to adapter/postgres, meant to be used by an upcoming sync usecase to pull live
// articles instead of relying on migration-seeded data.
package feed

import (
	"net/http"
	"time"
)

const fetchTimeout = 10 * time.Second

// Client fetches and parses one source's feed at a time - stateless beyond its http.Client, safe
// for concurrent use across sources.
type Client struct {
	httpClient *http.Client
}

// New builds a Client with a bounded per-request timeout, so one slow or unreachable feed can't
// hang a sync run indefinitely.
func New() *Client {
	return &Client{httpClient: &http.Client{Timeout: fetchTimeout}}
}
