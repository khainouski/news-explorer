package auth

// LoginView is what web/pages/login.html renders - a standalone page (no sidebar/topbar, like
// shared.NotFoundView) since there's nothing meaningful in the app shell before you're signed in.
type LoginView struct {
	PageTitle string

	// Login is whatever was submitted on a failed attempt, so the user doesn't have to retype it
	// (the password never round-trips back into the form).
	Login string

	Error string
}

// ChangePasswordView is what web/pages/change_password.html renders, reached via "Change
// Password" in the topbar's account dropdown (see web/components/navigation/topbar.html).
// Standalone page - no sidebar/topbar, same reasoning as LoginView: nothing there is relevant
// while you're mid-task on this one form. "Back to Sources" is the way out.
type ChangePasswordView struct {
	PageTitle string

	Error string
}
