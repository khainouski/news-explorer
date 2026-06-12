package user

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// UpdatePassword sets a new (already-hashed) password for a user.
func (r *Repo) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	ctx, span := tracer.Start(ctx, "adapter postgres user UpdatePassword")
	defer span.End()

	const q = `UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`

	if _, err := r.pool.Exec(ctx, q, userID, passwordHash); err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}
