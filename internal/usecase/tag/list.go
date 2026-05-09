package tag

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

// List returns every available tag.
func (u *UseCase) List(ctx context.Context) ([]domain.Tag, error) {
	return u.postgres.List(ctx)
}
