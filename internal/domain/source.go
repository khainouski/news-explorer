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

	UserID *int // nil = global source, visible to everyone

	FeedURL     string
	Description string
	Tag         Tag

	Badge      string // short label for the source's square icon, e.g. "Go", "TC"
	BadgeColor string // Tailwind background color class for the badge, e.g. "bg-blue-500"

	Status SourceStatus

	ArticleCount int        // computed live, not a stored column
	LastSyncedAt *time.Time // nil if never synced yet
}
