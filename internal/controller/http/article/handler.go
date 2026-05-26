// Package article handles the home feed (GET /) - the article list, its sort/search/tag filters,
// and the sources sidebar alongside it.
package article

import (
	usecasearticle "github.com/khainouski/news-explorer/internal/usecase/article"
	usecasesource "github.com/khainouski/news-explorer/internal/usecase/source"
	usecasetag "github.com/khainouski/news-explorer/internal/usecase/tag"
)

type Handler struct {
	article *usecasearticle.UseCase
	source  *usecasesource.UseCase
	tag     *usecasetag.UseCase
}

func NewHandler(article *usecasearticle.UseCase, source *usecasesource.UseCase, tag *usecasetag.UseCase) *Handler {
	return &Handler{article: article, source: source, tag: tag}
}
