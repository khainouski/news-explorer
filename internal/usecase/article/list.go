package article

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

// List returns every article, newest first.
func (u *UseCase) List(ctx context.Context) ([]domain.Article, error) {
	return u.postgres.List(ctx)
}
