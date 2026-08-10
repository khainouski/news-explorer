package shared

import "time"

// LastSyncedAgo renders lastSyncedAt for the sidebar's "Last synced" row - "Never" if no source
// has synced yet.
func LastSyncedAgo(lastSyncedAt *time.Time) string {
	if lastSyncedAt == nil {
		return "Never"
	}

	return TimeAgo(*lastSyncedAt)
}
