package middleware

import (
	"net/http"
	"net/url"
)

// RequireAdminPage bounces anyone who isn't the admin back to /sources with an explanatory toast.
func RequireAdminPage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := CurrentUser(r.Context()); user != nil && user.IsAdmin() {
			next.ServeHTTP(w, r)

			return
		}

		params := url.Values{
			"toast":   {"Admin only"},
			"message": {"This action is only available for the admin account."},
			"variant": {"warning"},
		}
		http.Redirect(w, r, "/sources?"+params.Encode(), http.StatusSeeOther)
	})
}

// RequireAdminAPI is RequireAdminPage's counterpart for fetch()-driven requests - a plain 403
// instead of a redirect, since there's no page to redirect to.
func RequireAdminAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := CurrentUser(r.Context()); user != nil && user.IsAdmin() {
			next.ServeHTTP(w, r)

			return
		}

		http.Error(w, "admin only", http.StatusForbidden)
	})
}
