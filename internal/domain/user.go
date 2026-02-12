package domain

// AdminUserID is the one seeded admin account's ID (see migration/postgres/20260507171000_seed.up.sql)
// - the app is single-admin for now, so "is this user id 1" is the authorization check itself
// (see internal/controller/http/middleware.RequireAdminPage/RequireAdminAPI). No role column:
// with only one possible value it would carry no information beyond the row existing at all -
// reintroduce one properly if/when there's more than one kind of account.
const AdminUserID = 1

// User is an account that can log in (see Session). PasswordHash is a bcrypt hash, never the raw
// password. Login is the sign-in identifier; Email is an optional profile field, never used for
// auth - nil for the seeded admin.
type User struct {
	ID           int
	Login        string
	Email        *string
	PasswordHash string
}

// IsAdmin reports whether this is the one admin account - see AdminUserID.
func (u User) IsAdmin() bool {
	return u.ID == AdminUserID
}
