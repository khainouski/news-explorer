// Package shared holds what more than one of the source/auth/article handler subpackages needs:
// page rendering, the 404 page, the topbar's account-dropdown data, and view types shared across
// pages (Badge, TagPill).
package shared

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/web"
)

// Render executes the named page against data, writing HTML to w.
func Render(w http.ResponseWriter, page string, data any) {
	if err := web.Render(w, page, data); err != nil {
		log.Error().Err(err).Str("page", page).Msg("render template")
	}
}

// RenderBlock renders one block of page - for HTMX requests swapping a fragment in place instead
// of navigating to the page.
func RenderBlock(w http.ResponseWriter, page, block string, data any) {
	if err := web.RenderBlock(w, page, block, data); err != nil {
		log.Error().Err(err).Str("page", page).Str("block", block).Msg("render template block")
	}
}
