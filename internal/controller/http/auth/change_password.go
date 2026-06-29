package auth

import (
	"errors"
	"net/http"
	"net/url"
	"unicode"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/internal/controller/http/middleware"
	"github.com/khainouski/news-explorer/internal/controller/http/shared"
	"github.com/khainouski/news-explorer/internal/domain"
)

const minPasswordLength = 8

// Account renders the change-password page - reached via "Change Password" in the topbar's
// account dropdown. Visitors without a session are sent to /login.
func (h *Handler) Account(w http.ResponseWriter, r *http.Request) {
	if middleware.CurrentUser(r.Context()) == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	shared.Render(w, "change_password", ChangePasswordView{PageTitle: "Change Password"})
}

// ChangePassword handles the change-password form. On failure it re-renders the form with an
// error instead of a bare 400 page - same pattern as the source package's Create.
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)

		return
	}

	current := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirm := r.FormValue("confirm_password")

	view := ChangePasswordView{PageTitle: "Change Password"}

	switch {
	case !isStrongPassword(newPassword):
		view.Error = "New password must be at least 8 characters and include uppercase and lowercase letters, a number and a special character."
		shared.Render(w, "change_password", view)

		return
	case newPassword != confirm:
		view.Error = "New password and confirmation don't match."
		shared.Render(w, "change_password", view)

		return
	}

	err := h.auth.ChangePassword(r.Context(), user.ID, user.PasswordHash, current, newPassword)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			view.Error = "Current password is incorrect."
		} else {
			log.Error().Err(err).Msg("change password")
			view.Error = "Something went wrong. Please try again."
		}

		shared.Render(w, "change_password", view)

		return
	}

	params := url.Values{
		"toast":   {"Password changed"},
		"message": {"Your password was updated successfully."},
	}
	http.Redirect(w, r, "/account?"+params.Encode(), http.StatusSeeOther)
}

// isStrongPassword mirrors the checklist shown on web/pages/change_password.html - at least 8
// characters, upper and lower case, a digit, and a special (neither letter nor digit) character.
func isStrongPassword(pw string) bool {
	if len(pw) < minPasswordLength {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, r := range pw {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
