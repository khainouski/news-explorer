package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

// GetByLogin returns one user by login, or domain.ErrUserNotFound.
func (r *Repo) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres user GetByLogin")
	defer span.End()

	const q = `SELECT id, login, email, password_hash FROM users WHERE login = $1`

	var u domain.User

	err := r.pool.QueryRow(ctx, q, login).Scan(&u.ID, &u.Login, &u.Email, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, fmt.Errorf("pool.QueryRow: %w", err)
	}

	return u, nil
}
