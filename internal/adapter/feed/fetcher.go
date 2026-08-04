package feed

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// maxBodyBytes caps how much of a feed response Fetch reads - a feed is plain text/XML, 5MB is
// already generous, and without a cap a misconfigured or malicious feed could exhaust memory.
const maxBodyBytes = 5 << 20

// Fetch downloads the raw feed body at feedURL. A non-2xx status or a declared Content-Type that
// isn't XML-shaped fails before any bytes are handed back, so Parse never has to guess whether
// it's looking at a feed or an error/redirect page a misconfigured server returned with a 200.
func (c *Client) Fetch(ctx context.Context, feedURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpClient.Do: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, feedURL)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "" && !isFeedContentType(ct) {
		return nil, fmt.Errorf("unexpected content-type %q fetching %s", ct, feedURL)
	}

	// Read one byte past the cap so a body that's exactly at the limit isn't mistaken for one
	// that's under it - ReadAll stops at maxBodyBytes+1 either way, never actually buffering more.
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes+1))
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	if len(body) > maxBodyBytes {
		return nil, fmt.Errorf("response body exceeds %d bytes fetching %s", maxBodyBytes, feedURL)
	}

	return body, nil
}

// isFeedContentType reports whether contentType (the response's Content-Type header) is one of
// the media types a feed is normally served as - the media type only, ignoring any charset
// parameter.
func isFeedContentType(contentType string) bool {
	mediaType, _, _ := strings.Cut(contentType, ";")

	switch strings.ToLower(strings.TrimSpace(mediaType)) {
	case "application/rss+xml", "application/atom+xml", "application/xml", "text/xml":
		return true
	default:
		return false
	}
}
