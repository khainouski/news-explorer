package source

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

// List returns every source visible to the caller. userID is always nil for now: nothing scopes
// a source to a user yet, so every request just sees the global source list.
func (u *UseCase) List(ctx context.Context) ([]domain.Source, error) {
	return u.postgres.List(ctx, nil)
}
