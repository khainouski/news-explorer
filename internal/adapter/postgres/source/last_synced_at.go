package source

import (
	"context"
	"fmt"
	"time"

	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// LastSyncedAt returns the most recent sync across every source, or nil if none has synced yet.
func (r *Repo) LastSyncedAt(ctx context.Context) (*time.Time, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres source LastSyncedAt")
	defer span.End()

	var t *time.Time

	if err := r.pool.QueryRow(ctx, `SELECT MAX(last_synced_at) FROM sources`).Scan(&t); err != nil {
		return nil, fmt.Errorf("pool.QueryRow: %w", err)
	}

	return t, nil
}
