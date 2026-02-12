package domain

import "time"

// Article is a single aggregated news item, pulled in from a Source.
type Article struct {
	ID          string
	Title       string
	Summary     string
	URL         string
	SourceID    string
	PublishedAt time.Time
	Unread      bool
}
