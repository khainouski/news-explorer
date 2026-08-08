package source

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// List returns every global source (user_id IS NULL) plus userID's own, if any.
func (r *Repo) List(ctx context.Context, userID *int) ([]domain.Source, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres source List")
	defer span.End()

	// article_count is counted live, not stored, so it can't go stale.
	const q = `
		SELECT s.id, s.user_id, s.name, s.feed_url, s.description, t.id, t.name, s.badge,
		       s.badge_color, s.status, COUNT(a.id) AS article_count, s.last_synced_at
		FROM sources s
		JOIN tags t ON t.id = s.tag_id
		LEFT JOIN articles a ON a.source_id = s.id
		WHERE s.user_id IS NULL OR s.user_id = $1
		GROUP BY s.id, t.id
		ORDER BY article_count DESC, s.name ASC`

	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var sources []domain.Source

	for rows.Next() {
		var src domain.Source

		err = rows.Scan(
			&src.ID,
			&src.UserID,
			&src.Name,
			&src.FeedURL,
			&src.Description,
			&src.Tag.ID,
			&src.Tag.Name,
			&src.Badge,
			&src.BadgeColor,
			&src.Status,
			&src.ArticleCount,
			&src.LastSyncedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		sources = append(sources, src)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return sources, nil
}
