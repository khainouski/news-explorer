package auth

// LoginView is what web/pages/login.html renders.
type LoginView struct {
	PageTitle string
	Login     string // whatever was submitted on a failed attempt - password never round-trips
	Error     string
}

// ChangePasswordView is what web/pages/change_password.html renders.
type ChangePasswordView struct {
	PageTitle string
	Error     string
}
