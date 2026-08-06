package sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/khainouski/news-explorer/internal/domain"
)

const workerCount = 8 // sources fetched concurrently - fetching only, not DB writes

const maxArticleAge = 90 * 24 * time.Hour // don't import a feed's entire back-catalog

func isRecent(publishedAt time.Time) bool {
	return time.Since(publishedAt) <= maxArticleAge
}

// Result summarizes one Sync run. Sync itself doesn't log - the caller decides what to do with it.
type Result struct {
	SourcesSynced    int
	SourcesFailed    int
	ArticlesInserted int
}

type fetchResult struct {
	source   domain.Source
	articles []domain.Article
	err      error
}

// Sync fetches every active source's feed concurrently (workerCount fetch workers) while a
// single writer drains the results and batches writes into Postgres, instead of each fetch worker
// competing for the connection pool. Every source's LastSyncedAt gets stamped, even on a failed
// fetch. One source's failure never stops the run - every failure is joined into the returned
// error.
func (u *UseCase) Sync(ctx context.Context) (Result, error) {
	sources, err := u.sources.ListActive(ctx)
	if err != nil {
		return Result{}, fmt.Errorf("sources.ListActive: %w", err)
	}

	// Buffered to fit every source at once, so neither the producer nor a worker ever blocks on a
	// channel op - cancellation just needs to reach feed.FetchArticles, not the channels too.
	jobs := make(chan domain.Source, len(sources))
	results := make(chan fetchResult, len(sources))

	for _, s := range sources {
		jobs <- s
	}
	close(jobs)

	var fetchWG sync.WaitGroup

	for range workerCount {
		fetchWG.Add(1)

		go func() {
			defer fetchWG.Done()
			u.fetchWorker(ctx, jobs, results)
		}()
	}

	go func() {
		fetchWG.Wait()
		close(results)
	}()

	res, errs := u.writeResults(ctx, results)

	return res, errors.Join(errs...)
}

func (u *UseCase) fetchWorker(ctx context.Context, jobs <-chan domain.Source, results chan<- fetchResult) {
	for s := range jobs {
		articles, err := u.feed.FetchArticles(ctx, s.ID, s.FeedURL)
		results <- fetchResult{source: s, articles: articles, err: err}
	}
}

// writeResults is the single DB-writing goroutine - it batches whatever's arrived since the last
// write instead of one write per source.
func (u *UseCase) writeResults(ctx context.Context, results <-chan fetchResult) (Result, []error) {
	var (
		total Result
		errs  []error
	)

	for r := range results {
		batch := []fetchResult{r}

		for range len(results) { // drain what's already queued, results is buffered
			batch = append(batch, <-results)
		}

		res, err := u.writeBatch(ctx, batch)

		total.SourcesSynced += res.SourcesSynced
		total.SourcesFailed += res.SourcesFailed
		total.ArticlesInserted += res.ArticlesInserted

		if err != nil {
			errs = append(errs, err)
		}
	}

	return total, errs
}

func (u *UseCase) writeBatch(ctx context.Context, batch []fetchResult) (Result, error) {
	var (
		articles  []domain.Article
		syncedIDs []string
		errs      []error
		res       Result
	)

	for _, r := range batch {
		syncedIDs = append(syncedIDs, r.source.ID)
		res.SourcesSynced++

		if r.err != nil {
			res.SourcesFailed++
			errs = append(errs, fmt.Errorf("source %s: feed.FetchArticles: %w", r.source.ID, r.err))

			continue
		}

		for _, a := range r.articles {
			if isRecent(a.PublishedAt) {
				articles = append(articles, a)
			}
		}
	}

	if len(articles) > 0 {
		inserted, err := u.articles.InsertBatch(ctx, articles)
		if err != nil {
			errs = append(errs, fmt.Errorf("articles.InsertBatch: %w", err))
		} else {
			res.ArticlesInserted = inserted
		}
	}

	if err := u.sources.UpdateLastSyncedBatch(ctx, syncedIDs); err != nil {
		errs = append(errs, fmt.Errorf("sources.UpdateLastSyncedBatch: %w", err))
	}

	return res, errors.Join(errs...)
}
