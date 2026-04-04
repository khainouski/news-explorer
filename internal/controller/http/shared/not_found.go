package shared

import "net/http"

// NotFoundView is what web/pages/not_found.html renders.
type NotFoundView struct {
	PageTitle string
}

// NotFound renders the styled 404 page - both for unmatched routes (chi's router-level NotFound,
// see router.go) and for routes that matched but whose resource doesn't exist (e.g. editing a
// source that's already been deleted - see the source subpackage).
func NotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	Render(w, "not_found", NotFoundView{PageTitle: "Page Not Found"})
}
