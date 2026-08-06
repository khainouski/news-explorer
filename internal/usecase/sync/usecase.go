// Package sync is the feed-sync use case - fetches every active source's feed, turns new entries
// into articles, and stamps each source's LastSyncedAt.
package sync

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

// Sources - not usecase/source.Postgres reused: ListActive filters by status in SQL, and
// UpdateLastSyncedBatch stamps many sources in one call.
type Sources interface {
	ListActive(ctx context.Context) ([]domain.Source, error)
	UpdateLastSyncedBatch(ctx context.Context, sourceIDs []string) error
}

// Articles - dedup is Postgres's job (InsertBatch's ON CONFLICT), so no separate existence check.
type Articles interface {
	InsertBatch(ctx context.Context, articles []domain.Article) (int, error)
}

// Feed depends only on domain.Article coming out, never on adapter/feed's own types.
type Feed interface {
	FetchArticles(ctx context.Context, sourceID, feedURL string) ([]domain.Article, error)
}

type UseCase struct {
	sources  Sources
	articles Articles
	feed     Feed
}

func New(sources Sources, articles Articles, feed Feed) *UseCase {
	return &UseCase{sources: sources, articles: articles, feed: feed}
}
