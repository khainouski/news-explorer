package sync

import "context"

// Sync fetches every active source's feed, creates any articles not already stored, and updates
// each source's LastSyncedAt regardless of whether new articles were found.
//
// TODO: Sources.ListActive; per source, Feed.Fetch -> feed.Parse -> feed.ToArticles ->
// Articles.InsertBatch (dedup is InsertBatch's ON CONFLICT, nothing to filter here); call
// Sources.UpdateLastSynced even on a fetch/parse error, so one broken feed doesn't get retried
// every run while the others succeed; one source's failure shouldn't abort the whole run.
func (u *UseCase) Sync(ctx context.Context) error {
	return nil
}
