package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/khainouski/news-explorer/internal/domain"
)

const sessionTTL = 7 * 24 * time.Hour

// dummyPasswordHash is compared against on an unknown login (see Login) purely to spend the same
// bcrypt time as a real comparison would - a fixed, valid bcrypt hash of an arbitrary password
// nobody is trying to guess, not a real credential.
const dummyPasswordHash = "$2a$10$d4L1N2X0CBo82t3YR6.s/uiDwWywv4InRre1qsfPAVle6Hy5B0rCK"

// Login verifies login/password and creates a new session, returning the token to set as the
// session cookie (see internal/controller/http/auth) - only its hash is ever stored.
// domain.ErrInvalidCredentials covers both "no such user" and "wrong password" - a login failure
// should never reveal which logins exist, including via response timing: an unknown login still
// runs a bcrypt comparison (against dummyPasswordHash) before returning, so it takes the same
// time as a known login with a wrong password instead of returning early.
func (u *UseCase) Login(ctx context.Context, login, password string) (string, error) {
	user, err := u.users.GetByLogin(ctx, strings.TrimSpace(login))
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return "", fmt.Errorf("users.GetByLogin: %w", err)
		}

		bcrypt.CompareHashAndPassword([]byte(dummyPasswordHash), []byte(password)) //nolint:errcheck // deliberately discarded, see dummyPasswordHash

		return "", domain.ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, hash, err := newSessionToken()
	if err != nil {
		return "", fmt.Errorf("newSessionToken: %w", err)
	}

	err = u.sessions.Create(ctx, domain.Session{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(sessionTTL),
	})
	if err != nil {
		return "", fmt.Errorf("sessions.Create: %w", err)
	}

	return token, nil
}
