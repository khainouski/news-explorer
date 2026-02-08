package router

import (
	"context"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ExtractPath returns the route pattern chi matched (e.g. "/users/{id}"), not the raw URL -
// keeps metric/span/log labels low-cardinality instead of one series per literal path.
func ExtractPath(ctx context.Context) string {
	chiCtx := chi.RouteContext(ctx)
	path := strings.Join(chiCtx.RoutePatterns, "")
	path = strings.ReplaceAll(path, "/*/", "/")

	return path
}
