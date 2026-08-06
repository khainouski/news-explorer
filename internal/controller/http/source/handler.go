// Package source handles the /sources pages: the list (search/sort/tag filter), add, edit,
// delete, and triggering a sync run.
package source

import (
	usecasesource "github.com/khainouski/news-explorer/internal/usecase/source"
	usecasesync "github.com/khainouski/news-explorer/internal/usecase/sync"
	usecasetag "github.com/khainouski/news-explorer/internal/usecase/tag"
)

type Handler struct {
	source *usecasesource.UseCase
	tag    *usecasetag.UseCase
	sync   *usecasesync.UseCase
}

func NewHandler(source *usecasesource.UseCase, tag *usecasetag.UseCase, sync *usecasesync.UseCase) *Handler {
	return &Handler{source: source, tag: tag, sync: sync}
}
