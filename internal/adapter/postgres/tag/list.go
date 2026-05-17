package tag

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// List returns every available tag, e.g. to populate the "Add Source" form's dropdown.
func (r *Repo) List(ctx context.Context) ([]domain.Tag, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres tag List")
	defer span.End()

	const q = `SELECT id, name FROM tags ORDER BY name`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var tags []domain.Tag

	for rows.Next() {
		var t domain.Tag

		if err = rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		tags = append(tags, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return tags, nil
}
