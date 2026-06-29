package auth

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/internal/controller/http/middleware"
)

// Logout deletes the current session (both server-side and the cookie) and redirects home -
// immediate, no confirmation prompt, same philosophy as the source package's Delete.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(middleware.SessionCookieName); err == nil {
		if err = h.auth.Logout(r.Context(), cookie.Value); err != nil {
			log.Error().Err(err).Msg("logout")
		}
	}

	clearSessionCookie(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
