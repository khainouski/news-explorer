package domain

// AdminUserID is the one seeded admin account's ID - the app is single-admin, so "is this user id
// 1" is the authorization check itself. No role column until there's more than one kind of account.
const AdminUserID = 1

// User is an account that can log in. PasswordHash is a bcrypt hash, never the raw password.
// Email is an optional profile field, never used for auth.
type User struct {
	ID           int
	Login        string
	Email        *string
	PasswordHash string
}

func (u User) IsAdmin() bool {
	return u.ID == AdminUserID
}
