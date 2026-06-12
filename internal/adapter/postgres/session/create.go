package session

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

func (r *Repo) Create(ctx context.Context, s domain.Session) error {
	ctx, span := tracer.Start(ctx, "adapter postgres session Create")
	defer span.End()

	const q = `INSERT INTO sessions (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`

	if _, err := r.pool.Exec(ctx, q, s.UserID, s.TokenHash, s.ExpiresAt); err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}
