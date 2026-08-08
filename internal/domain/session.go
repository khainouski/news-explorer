package domain

import "time"

// Session is a logged-in browser session. TokenHash is the SHA-256 hash of the raw token set as
// the session cookie - the raw token itself is never stored.
type Session struct {
	UserID    int
	TokenHash []byte
	ExpiresAt time.Time
}
