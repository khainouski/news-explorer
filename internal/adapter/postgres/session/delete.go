package session

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// Delete removes a session by its token hash - a no-op if it's already gone.
func (r *Repo) Delete(ctx context.Context, tokenHash []byte) error {
	ctx, span := tracer.Start(ctx, "adapter postgres session Delete")
	defer span.End()

	const q = `DELETE FROM sessions WHERE token_hash = $1`

	if _, err := r.pool.Exec(ctx, q, tokenHash); err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}
