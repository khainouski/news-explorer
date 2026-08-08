package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/internal/controller/http/middleware"
	"github.com/khainouski/news-explorer/internal/controller/http/shared"
	"github.com/khainouski/news-explorer/internal/domain"
)

const sessionCookieTTL = 7 * 24 * time.Hour

// Login renders the sign-in page. Already-logged-in visitors are sent to the home feed - there's
// nothing to sign in for again.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if middleware.CurrentUser(r.Context()) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	shared.Render(w, "login", LoginView{PageTitle: "Log in"})
}

// LoginSubmit handles the sign-in form. On failure it re-renders the form with an error instead
// of a bare 401 page.
func (h *Handler) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)

		return
	}

	login := strings.TrimSpace(r.FormValue("login"))
	password := r.FormValue("password")

	token, err := h.auth.Login(r.Context(), login, password)
	if err != nil {
		view := LoginView{PageTitle: "Log in", Login: login}

		if errors.Is(err, domain.ErrInvalidCredentials) {
			view.Error = "Incorrect login or password."
		} else {
			log.Error().Err(err).Msg("login")
			view.Error = "Something went wrong. Please try again."
		}

		shared.Render(w, "login", view)

		return
	}

	setSessionCookie(w, r, token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func setSessionCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(sessionCookieTTL),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   isSecureRequest(r),
	})
}

func clearSessionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   isSecureRequest(r),
	})
}

// isSecureRequest checks X-Forwarded-Proto too, since Traefik terminates TLS in front of the app.
func isSecureRequest(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}
