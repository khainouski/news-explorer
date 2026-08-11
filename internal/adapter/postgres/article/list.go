package article

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// whereClause is shared by List and count - $1 source_id, $2 tag_id, $3 search text, all
// empty-string-means-no-filter.
const whereClause = `
	WHERE ($1 = '' OR a.source_id = $1)
	  AND ($2 = '' OR s.tag_id = $2)
	  AND ($3 = '' OR a.title ILIKE '%' || $3 || '%' OR a.summary ILIKE '%' || $3 || '%' OR s.name ILIKE '%' || $3 || '%')`

// List returns one page of articles matching params, plus the total count of matching rows -
// a separate count query rather than COUNT(*) OVER(), which would report 0 whenever Offset lands
// past the last matching row (Postgres has no row left to carry the window value on).
func (r *Repo) List(ctx context.Context, params domain.ArticleListParams) ([]domain.Article, int, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres article List")
	defer span.End()

	total, err := r.count(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	if total == 0 {
		return nil, 0, nil
	}

	order := "DESC"
	if params.Oldest {
		order = "ASC"
	}

	// order is one of the two literals above, never user input - safe to interpolate. Every
	// other value stays a placeholder.
	q := fmt.Sprintf(`
		SELECT a.id, a.source_id, a.title, a.summary, a.url, a.published_at, a.unread
		FROM articles a
		JOIN sources s ON s.id = a.source_id
		%s
		ORDER BY a.published_at %s
		LIMIT NULLIF($4::int, 0) OFFSET $5`, whereClause, order)

	rows, err := r.pool.Query(ctx, q, params.SourceID, params.TagID, params.Query, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var articles []domain.Article

	for rows.Next() {
		var a domain.Article

		err = rows.Scan(&a.ID, &a.SourceID, &a.Title, &a.Summary, &a.URL, &a.PublishedAt, &a.Unread)
		if err != nil {
			return nil, 0, fmt.Errorf("rows.Scan: %w", err)
		}

		articles = append(articles, a)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows.Err: %w", err)
	}

	return articles, total, nil
}

func (r *Repo) count(ctx context.Context, params domain.ArticleListParams) (int, error) {
	q := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM articles a
		JOIN sources s ON s.id = a.source_id
		%s`, whereClause)

	var total int

	err := r.pool.QueryRow(ctx, q, params.SourceID, params.TagID, params.Query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("pool.QueryRow: %w", err)
	}

	return total, nil
}
