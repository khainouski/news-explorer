// Package article is the articles use case - currently just listing them (the home feed). No
// single-article view exists (article rows link straight to the source's own URL), so there's no
// Get here - added if/when that changes, not before.
package article

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

type Postgres interface {
	List(ctx context.Context) ([]domain.Article, error)
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{postgres: postgres}
}
