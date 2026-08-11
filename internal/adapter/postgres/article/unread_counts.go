package article

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// UnreadCountsBySource returns the count of unread articles per source, across the whole table -
// independent of any List filter/page, since the sidebar badges it feeds always reflect the true
// global state.
func (r *Repo) UnreadCountsBySource(ctx context.Context) (map[string]int, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres article UnreadCountsBySource")
	defer span.End()

	const q = `SELECT source_id, COUNT(*) FROM articles WHERE unread GROUP BY source_id`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)

	for rows.Next() {
		var (
			sourceID string
			count    int
		)

		if err = rows.Scan(&sourceID, &count); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		counts[sourceID] = count
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return counts, nil
}
