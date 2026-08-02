// Package sync is the feed-sync use case - fetches every active source's feed, turns new entries
// into articles, and stamps each source's LastSyncedAt. Needs three narrow interfaces from two
// different adapters (Sources and Articles from adapter/postgres, Feed from adapter/feed) - same
// multi-interface-per-usecase shape as usecase/auth's Users/Sessions, not one combined interface
// forcing adapter-side glue code.
package sync

import (
	"context"

	"github.com/khainouski/news-explorer/internal/domain"
)

// Sources deliberately isn't usecase/source.Postgres reused - List there returns every source
// regardless of status (the /sources page shows both), but a sync run only ever wants the active
// ones. ListActive filters by status in SQL instead of fetching everything and discarding
// inactive sources in Go.
type Sources interface {
	ListActive(ctx context.Context) ([]domain.Source, error)
	UpdateLastSynced(ctx context.Context, sourceID string) error
}

// Articles has just the one method - dedup against what's already stored is PostgreSQL's job
// (InsertBatch's ON CONFLICT (source_id, external_id) DO NOTHING), not an app-side existence
// check, so there's no separate "what already exists" query here.
type Articles interface {
	InsertBatch(ctx context.Context, articles []domain.Article) (int, error)
}

// Feed only covers Fetch - feed.Parse/feed.ToArticles are plain functions, not part of this
// interface, since they're pure and don't touch the network (nothing to fake in a usecase test).
type Feed interface {
	Fetch(ctx context.Context, feedURL string) ([]byte, error)
}

type UseCase struct {
	sources  Sources
	articles Articles
	feed     Feed
}

func New(sources Sources, articles Articles, feed Feed) *UseCase {
	return &UseCase{sources: sources, articles: articles, feed: feed}
}
