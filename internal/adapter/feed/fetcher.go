package feed

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const maxBodyBytes = 5 << 20 // cap a feed response can't exceed

const userAgent = "NewsExplorer/1.0 (+https://goskills.xyz)" // some feeds 403 on Go's default UA

// Fetch downloads the raw feed body at feedURL. Content-Type isn't checked - real feeds get it
// wrong often enough that Parse's own root-element check is the more reliable gate.
func (c *Client) Fetch(ctx context.Context, feedURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpClient.Do: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, feedURL)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes+1))
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	if len(body) > maxBodyBytes {
		return nil, fmt.Errorf("response body exceeds %d bytes fetching %s", maxBodyBytes, feedURL)
	}

	return body, nil
}
