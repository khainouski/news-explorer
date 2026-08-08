package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// newSessionToken generates a random session token and its SHA-256 hash - only the hash is
// stored; the encoded token is what goes in the cookie.
func newSessionToken() (encoded string, hash []byte, err error) {
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", nil, err
	}

	sum := sha256.Sum256(raw)

	return base64.RawURLEncoding.EncodeToString(raw), sum[:], nil
}

// hashSessionToken re-derives a token's hash from its encoded (cookie) form for lookup.
func hashSessionToken(encoded string) ([]byte, error) {
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	sum := sha256.Sum256(raw)

	return sum[:], nil
}
