package source

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// Get returns one source by ID, or domain.ErrSourceNotFound.
func (r *Repo) Get(ctx context.Context, id string) (domain.Source, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres source Get")
	defer span.End()

	// article_count isn't a stored column - see List.
	const q = `
		SELECT s.id, s.user_id, s.name, s.feed_url, s.description, t.id, t.name, s.badge,
		       s.badge_color, s.status, COUNT(a.id) AS article_count, s.last_synced_at
		FROM sources s
		JOIN tags t ON t.id = s.tag_id
		LEFT JOIN articles a ON a.source_id = s.id
		WHERE s.id = $1
		GROUP BY s.id, t.id`

	var src domain.Source

	err := r.pool.QueryRow(ctx, q, id).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Source{}, domain.ErrSourceNotFound
		}

		return domain.Source{}, fmt.Errorf("pool.QueryRow: %w", err)
	}

	return src, nil
}
