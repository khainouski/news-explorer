// Package middleware holds chi middleware for the main route group - session resolution and the
// admin-only gate built on top of it. A leaf package (only depends on domain/usecase) so both the
// top-level router and the source/auth/article handler subpackages can import it without a cycle.
package middleware

import (
	"context"
	"net/http"

	"github.com/khainouski/news-explorer/internal/domain"
	usecaseauth "github.com/khainouski/news-explorer/internal/usecase/auth"
)

type ctxKey int

const userCtxKey ctxKey = iota

// SessionCookieName is the cookie both Auth (reading it) and the auth subpackage's login/logout
// handlers (setting/clearing it) agree on.
const SessionCookieName = "session"

// Auth attaches the logged-in user (if any) to the request context, resolved from the session
// cookie. It never blocks a request - an absent/invalid/expired cookie just leaves
// CurrentUser(ctx) nil, same as a first-time visitor; see the auth subpackage for the handlers
// that actually require a session, and RequireAdminPage/RequireAdminAPI for the ones that require
// the admin specifically.
func Auth(uc *usecaseauth.UseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				next.ServeHTTP(w, r)

				return
			}

			user, err := uc.CurrentUser(r.Context(), cookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)

				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userCtxKey, &user)))
		})
	}
}

// CurrentUser returns the logged-in user attached by Auth, or nil if none.
func CurrentUser(ctx context.Context) *domain.User {
	u, _ := ctx.Value(userCtxKey).(*domain.User)

	return u
}
