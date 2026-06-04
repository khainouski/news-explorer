package auth

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

// CurrentUser resolves a raw (cookie) session token to its user, or domain.ErrSessionNotFound -
// see internal/controller/http/middleware.Auth, which calls this on every request.
func (u *UseCase) CurrentUser(ctx context.Context, token string) (domain.User, error) {
	hash, err := hashSessionToken(token)
	if err != nil {
		return domain.User{}, domain.ErrSessionNotFound
	}

	return u.sessions.GetUser(ctx, hash)
}
