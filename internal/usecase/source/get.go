package source

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

// Get returns one source by ID, or domain.ErrSourceNotFound.
func (u *UseCase) Get(ctx context.Context, id string) (domain.Source, error) {
	return u.postgres.Get(ctx, id)
}
