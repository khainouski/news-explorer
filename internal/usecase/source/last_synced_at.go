package source

import (
	"context"
	"time"
)

// LastSyncedAt returns the most recent sync across every source - drives the sidebar's
// "Last synced" row, shown on every page.
func (u *UseCase) LastSyncedAt(ctx context.Context) (*time.Time, error) {
	return u.postgres.LastSyncedAt(ctx)
}
