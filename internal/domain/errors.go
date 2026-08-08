package domain

import "errors"

var (
	ErrSourceExists   = errors.New("a source with this name already exists")
	ErrTagNotFound    = errors.New("unknown tag")
	ErrSourceNotFound = errors.New("source not found")
	ErrInvalidName    = errors.New("name must contain at least one letter or number")
	ErrUserNotFound   = errors.New("user not found")

	// ErrInvalidCredentials deliberately covers both "no such user" and "wrong password" - a
	// login failure should never reveal which logins exist.
	ErrInvalidCredentials = errors.New("invalid login or password")

	ErrSessionNotFound = errors.New("session not found")
)
