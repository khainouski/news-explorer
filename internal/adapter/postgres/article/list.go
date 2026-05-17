package article

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// List returns every article, newest first.
func (r *Repo) List(ctx context.Context) ([]domain.Article, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres article List")
	defer span.End()

	const q = `
		SELECT id, source_id, title, summary, url, published_at, unread
		FROM articles
		ORDER BY published_at DESC`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var articles []domain.Article

	for rows.Next() {
		var a domain.Article

		err = rows.Scan(&a.ID, &a.SourceID, &a.Title, &a.Summary, &a.URL, &a.PublishedAt, &a.Unread)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		articles = append(articles, a)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return articles, nil
}
