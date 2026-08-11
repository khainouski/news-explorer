package article

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

// List returns one page of articles matching params, plus the total count of matching rows.
func (u *UseCase) List(ctx context.Context, params domain.ArticleListParams) ([]domain.Article, int, error) {
	return u.postgres.List(ctx, params)
}
