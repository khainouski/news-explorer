package domain

import "time"

// SourceStatus is whether a source's feed is currently being synced.
type SourceStatus string

const (
	SourceStatusActive   SourceStatus = "active"
	SourceStatusInactive SourceStatus = "inactive"
)

// Source is a feed news articles get aggregated from (a blog, an RSS feed, ...).
type Source struct {
	ID   string
	Name string

	// UserID is nil for a global/default source (visible to every user - what every seeded
	// source is today, and the only kind any creation flow makes). Once source ownership exists,
	// a non-nil UserID (a real FK to users.id) would scope a source to whoever added it.
	UserID *int

	FeedURL     string
	Description string
	Tag         Tag

	Badge      string // short label for the source's square icon, e.g. "Go", "TC"
	BadgeColor string // Tailwind background color class for the badge, e.g. "bg-blue-500"

	Status SourceStatus

	// ArticleCount is COUNT(*) of this source's rows in the articles table, computed live by the
	// adapter query - not a stored column, so it can't go stale.
	ArticleCount int
	LastSyncedAt *time.Time // nil if never synced yet
}
