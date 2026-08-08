// Package auth is the authentication use case - login/logout, resolving a session cookie back to
// its user, and changing a password.
package auth

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

type Users interface {
	GetByLogin(ctx context.Context, login string) (domain.User, error)
	UpdatePassword(ctx context.Context, userID int, passwordHash string) error
}

type Sessions interface {
	Create(ctx context.Context, s domain.Session) error
	GetUser(ctx context.Context, tokenHash []byte) (domain.User, error)
	Delete(ctx context.Context, tokenHash []byte) error
}

type UseCase struct {
	users    Users
	sessions Sessions
}

func New(users Users, sessions Sessions) *UseCase {
	return &UseCase{users: users, sessions: sessions}
}
