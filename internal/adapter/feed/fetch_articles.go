package feed

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/internal/domain"
)

// FetchArticles fetches, parses, and maps one source's feed into domain.Article values.
func (c *Client) FetchArticles(ctx context.Context, sourceID, feedURL string) ([]domain.Article, error) {
	raw, err := c.Fetch(ctx, feedURL)
	if err != nil {
		return nil, fmt.Errorf("Fetch: %w", err)
	}

	items, err := Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("Parse: %w", err)
	}

	return ToArticles(sourceID, items), nil
}
