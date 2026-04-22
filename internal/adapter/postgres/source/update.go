package source

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// Update updates every editable field of an existing source. The ID (and so its slug/URL) never
// changes once created.
func (r *Repo) Update(ctx context.Context, s domain.Source) error {
	ctx, span := tracer.Start(ctx, "adapter postgres source Update")
	defer span.End()

	const q = `
		UPDATE sources
		SET name = $2, feed_url = $3, description = $4, tag_id = $5, badge = $6,
		    badge_color = $7, status = $8, updated_at = NOW()
		WHERE id = $1`

	cmdTag, err := r.pool.Exec(ctx, q, s.ID, s.Name, s.FeedURL, s.Description, s.Tag.ID, s.Badge, s.BadgeColor, s.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrForeignKeyViolation {
			return domain.ErrTagNotFound
		}

		return fmt.Errorf("pool.Exec: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrSourceNotFound
	}

	return nil
}
