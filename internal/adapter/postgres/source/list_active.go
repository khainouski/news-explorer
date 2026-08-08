package source

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// ListActive returns every active source for the sync usecase - unscoped (a sync run has no
// viewer), and Source.ArticleCount is always 0 here, not a real count.
func (r *Repo) ListActive(ctx context.Context) ([]domain.Source, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres source ListActive")
	defer span.End()

	const q = `
		SELECT s.id, s.user_id, s.name, s.feed_url, s.description, t.id, t.name, s.badge,
		       s.badge_color, s.status, s.last_synced_at
		FROM sources s
		JOIN tags t ON t.id = s.tag_id
		WHERE s.status = 'active'
		ORDER BY s.name ASC`

	rows, err := r.pool.Query(ctx, q)
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
