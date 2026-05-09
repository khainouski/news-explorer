// Package tag is the tags use case - just listing them, e.g. for the "Add Source" form's
// dropdown and the "Filter by tag" pills.
package tag

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

type Postgres interface {
	List(ctx context.Context) ([]domain.Tag, error)
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{postgres: postgres}
}
