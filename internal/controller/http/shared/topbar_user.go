package shared

import (
	"net/http"
	"strings"
	"unicode"

	"github.com/khainouski/news-explorer/internal/controller/http/middleware"
)

// TopbarUser drives topbar.html's login button vs. account dropdown. Zero value renders the
// login button.
type TopbarUser struct {
	LoggedIn bool
	Initials string
	Name     string
}

func BuildTopbarUser(r *http.Request) TopbarUser {
	user := middleware.CurrentUser(r.Context())
	if user == nil {
		return TopbarUser{}
	}

	return TopbarUser{
		LoggedIn: true,
		Initials: loginInitials(user.Login),
		Name:     capitalize(user.Login),
	}
}

// loginInitials derives avatar initials from a login, splitting on '.', '_' and '-' the way
// multi-word logins are commonly written, e.g. "jane.doe" -> "JD", "admin" -> "AD".
func loginInitials(login string) string {
	words := strings.FieldsFunc(login, func(r rune) bool {
		return r == '.' || r == '_' || r == '-'
	})

	switch {
	case len(words) >= 2:
		return strings.ToUpper(string([]rune(words[0])[:1]) + string([]rune(words[1])[:1]))
	case len(words) == 1 && len([]rune(words[0])) >= 2:
		return strings.ToUpper(string([]rune(words[0])[:2]))
	case len(words) == 1:
		return strings.ToUpper(words[0])
	default:
		return "?"
	}
}

func capitalize(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return s
	}

	r[0] = unicode.ToUpper(r[0])

	return string(r)
}
