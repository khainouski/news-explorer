package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// newSessionToken generates a random session token and its SHA-256 hash - the hash is what gets
// stored (see internal/domain.Session.TokenHash), the encoded token is what gets handed to the
// browser as the cookie value (see internal/controller/http/auth).
func newSessionToken() (encoded string, hash []byte, err error) {
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", nil, err
	}

	sum := sha256.Sum256(raw)

	return base64.RawURLEncoding.EncodeToString(raw), sum[:], nil
}

// hashSessionToken re-derives a session token's hash from its encoded (cookie) form, so it can be
// looked up - errors on a malformed/tampered cookie value.
func hashSessionToken(encoded string) ([]byte, error) {
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	sum := sha256.Sum256(raw)

	return sum[:], nil
}
