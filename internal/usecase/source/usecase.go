// Package source is the sources use case - list/get/create/update/delete a source, and mark its
// articles read. Postgres here is a narrow interface matching internal/adapter/postgres/source.Repo
// exactly, so that Repo satisfies it directly with no adapter-side wrapper needed.
package source

import (
	"context"
	"time"

	"github.com/khainouski/news-explorer/internal/domain"
)

type Postgres interface {
	List(ctx context.Context, userID *int) ([]domain.Source, error)
	Get(ctx context.Context, id string) (domain.Source, error)
	Create(ctx context.Context, s domain.Source) error
	Update(ctx context.Context, s domain.Source) error
	Delete(ctx context.Context, id string) (string, error)
	MarkRead(ctx context.Context, sourceID string) error
	LastSyncedAt(ctx context.Context) (*time.Time, error)
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{postgres: postgres}
}
