package domain

import "time"

// Session is a logged-in browser session. TokenHash is the SHA-256 hash of the random token set
// as the session cookie - the raw token is never stored, only ever handed to the browser once at
// login (see internal/usecase/auth.Login).
type Session struct {
	UserID    int
	TokenHash []byte
	ExpiresAt time.Time
}
