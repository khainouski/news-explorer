package middleware

import (
	"net/http"
	"net/url"
)

// RequireAdminPage bounces anyone who isn't the admin (see domain.User.IsAdmin) back to /sources
// with an explanatory toast instead of letting them through - "+ Add Source" and the edit icon
// are visible to everyone (see web/pages/sources.html), so this is what actually happens when
// someone without an admin session clicks them.
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

// RequireAdminAPI is RequireAdminPage's counterpart for the JS-driven delete request (see
// web/static/js/source_form.js) - there's no page to redirect a fetch() call to, so it's a plain
// 403 instead. In practice this is unreachable through the UI: the delete button only exists on
// the Edit Source page, which RequireAdminPage already keeps non-admins off of.
func RequireAdminAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := CurrentUser(r.Context()); user != nil && user.IsAdmin() {
			next.ServeHTTP(w, r)

			return
		}

		http.Error(w, "admin only", http.StatusForbidden)
	})
}
