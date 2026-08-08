package shared

import "net/http"

// NotFoundView is what web/pages/not_found.html renders.
type NotFoundView struct {
	PageTitle string
}

// NotFound renders the styled 404 page - unmatched routes and missing resources both use this.
func NotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	Render(w, "not_found", NotFoundView{PageTitle: "Page Not Found"})
}
