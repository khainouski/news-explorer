package source

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// Create inserts a new source. ArticleCount starts at 0 implicitly (it's counted live from
// articles, not stored - see List); LastSyncedAt starts NULL (never synced).
func (r *Repo) Create(ctx context.Context, s domain.Source) error {
	ctx, span := tracer.Start(ctx, "adapter postgres source Create")
	defer span.End()

	const q = `
		INSERT INTO sources (id, name, feed_url, description, tag_id, badge, badge_color, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.pool.Exec(ctx, q, s.ID, s.Name, s.FeedURL, s.Description, s.Tag.ID, s.Badge, s.BadgeColor, s.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgErrUniqueViolation:
				return domain.ErrSourceExists
			case pgErrForeignKeyViolation:
				return domain.ErrTagNotFound
			}
		}

		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}
