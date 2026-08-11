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

// ArticleListParams filters/sorts/paginates List - zero value means no filter, Limit 0 means no
// page (every matching row).
type ArticleListParams struct {
	SourceID string
	TagID    string
	Query    string // matched against title/summary/source name
	Oldest   bool   // false = newest first

	Limit  int
	Offset int
}
