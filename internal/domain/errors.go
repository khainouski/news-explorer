package domain

import "errors"

var (
	// ErrSourceExists means the generated/given source ID already exists (sources.id is a slug
	// of the name - two sources with the same name collide).
	ErrSourceExists = errors.New("a source with this name already exists")

	// ErrTagNotFound means the given tag ID doesn't exist in the tags table.
	ErrTagNotFound = errors.New("unknown tag")

	// ErrSourceNotFound means no source has the given ID.
	ErrSourceNotFound = errors.New("source not found")

	// ErrInvalidName means the name has no alphanumeric characters, so it can't be turned into a
	// non-empty slug for the source's ID.
	ErrInvalidName = errors.New("name must contain at least one letter or number")

	// ErrUserNotFound means no user has the given login.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCredentials means the login/password (or, on the change-password form, the
	// current password) didn't match - deliberately the same error for "no such user" and
	// "wrong password" so a login failure never reveals which logins exist.
	ErrInvalidCredentials = errors.New("invalid login or password")

	// ErrSessionNotFound means the session cookie's token doesn't match any live (non-expired)
	// session.
	ErrSessionNotFound = errors.New("session not found")
)
