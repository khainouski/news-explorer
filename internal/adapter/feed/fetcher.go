package feed

import "context"

// Fetch downloads the raw feed body at feedURL.
//
// TODO: GET with ctx, follow a bounded number of redirects, cap the response size (a misconfigured
// or malicious feed shouldn't be able to exhaust memory), and reject non-2xx/non-XML responses
// before returning.
func (c *Client) Fetch(ctx context.Context, feedURL string) ([]byte, error) {
	return nil, nil
}
