package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// GetUser returns the user for a still-valid (not expired) session, or domain.ErrSessionNotFound.
func (r *Repo) GetUser(ctx context.Context, tokenHash []byte) (domain.User, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres session GetUser")
	defer span.End()

	const q = `
		SELECT u.id, u.login, u.email, u.password_hash
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > NOW()`

	var u domain.User

	err := r.pool.QueryRow(ctx, q, tokenHash).Scan(&u.ID, &u.Login, &u.Email, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrSessionNotFound
		}

		return domain.User{}, fmt.Errorf("pool.QueryRow: %w", err)
	}

	return u, nil
}
