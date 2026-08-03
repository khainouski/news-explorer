package source

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// UpdateLastSynced stamps a source's last_synced_at to now - called once per source at the end of
// a sync run, whether or not it found any new articles (see usecase/sync.Sync). Like MarkRead, a
// no-op if sourceID doesn't match anything - by the time this runs, ListActive already fetched
// this exact source, so it should always exist.
func (r *Repo) UpdateLastSynced(ctx context.Context, sourceID string) error {
	ctx, span := tracer.Start(ctx, "adapter postgres source UpdateLastSynced")
	defer span.End()

	const q = `UPDATE sources SET last_synced_at = NOW() WHERE id = $1`

	if _, err := r.pool.Exec(ctx, q, sourceID); err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}
