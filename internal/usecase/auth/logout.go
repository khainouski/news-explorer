package auth

import "context"

// Logout deletes the session behind a raw (cookie) token. A malformed token is treated the same
// as an already-gone session - nothing worth surfacing as an error either way.
func (u *UseCase) Logout(ctx context.Context, token string) error {
	hash, err := hashSessionToken(token)
	if err != nil {
		return nil
	}

	return u.sessions.Delete(ctx, hash)
}
