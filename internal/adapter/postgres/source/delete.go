package source

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// Delete removes a source and, via ON DELETE CASCADE, every article that came from it. Returns
// the deleted source's name (e.g. for a confirmation toast).
func (r *Repo) Delete(ctx context.Context, id string) (string, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres source Delete")
	defer span.End()

	var name string

	err := r.pool.QueryRow(ctx, `DELETE FROM sources WHERE id = $1 RETURNING name`, id).Scan(&name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrSourceNotFound
		}

		return "", fmt.Errorf("pool.QueryRow: %w", err)
	}

	return name, nil
}
