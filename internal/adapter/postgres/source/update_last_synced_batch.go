package source

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// UpdateLastSyncedBatch stamps last_synced_at to now for every ID in sourceIDs in one round trip.
func (r *Repo) UpdateLastSyncedBatch(ctx context.Context, sourceIDs []string) error {
	ctx, span := tracer.Start(ctx, "adapter postgres source UpdateLastSyncedBatch")
	defer span.End()

	if len(sourceIDs) == 0 {
		return nil
	}

	const q = `UPDATE sources SET last_synced_at = NOW() WHERE id = ANY($1)`

	if _, err := r.pool.Exec(ctx, q, sourceIDs); err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}
