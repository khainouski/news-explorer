package source

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

func (r *Repo) MarkRead(ctx context.Context, sourceID string) error {
	ctx, span := tracer.Start(ctx, "adapter postgres source MarkRead")
	defer span.End()

	const q = `UPDATE articles SET unread = FALSE, updated_at = NOW() WHERE source_id = $1 AND unread`

	if _, err := r.pool.Exec(ctx, q, sourceID); err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}
