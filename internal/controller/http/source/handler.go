// Package source handles the /sources pages: the list (search/sort/tag filter), add, edit and
// delete.
package source

import (
	usecasesource "github.com/khainouski/news-explorer/internal/usecase/source"
	usecasetag "github.com/khainouski/news-explorer/internal/usecase/tag"
)

type Handler struct {
	source *usecasesource.UseCase
	tag    *usecasetag.UseCase
}

func NewHandler(source *usecasesource.UseCase, tag *usecasetag.UseCase) *Handler {
	return &Handler{source: source, tag: tag}
}
