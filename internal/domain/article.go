package domain

import "time"

// Article is a single aggregated news item, pulled in from a Source.
type Article struct {
	ID         string // our own ID - a slug, like Source.ID (see adapter/postgres/source's name-to-slug)
	SourceID   string
	ExternalID string // stable identifier from the source's own feed - see adapter/feed.ToArticles

	Title       string
	Summary     string
	URL         string
	PublishedAt time.Time
	Unread      bool
}
