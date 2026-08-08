// Package feed adapts external RSS/Atom sources into domain-shaped data.
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

// Client is safe for concurrent use - stateless beyond its http.Client.
type Client struct {
	httpClient *http.Client
}

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
