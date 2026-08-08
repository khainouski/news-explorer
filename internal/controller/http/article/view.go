package article

import "github.com/khainouski/news-explorer/internal/controller/http/shared"

// HomeView is what web/pages/home.html renders: the article feed plus the sources sidebar.
type HomeView struct {
	PageTitle   string
	Active      string
	SearchScope string // "articles" - tells topbar.html which results container the search box targets
	Articles    []ArticleView
	Sources     []SourceView
	TotalCount  int
	SourceCount int

	Query string // current search term (?q=), matched against title/summary/source name

	// SortBy/SourceID/TagID are the raw ?sort=/?source=/?tag= values, picked up by the topbar's
	// search form via hx-include so searching doesn't drop them.
	SortBy   string
	SourceID string
	TagID    string

	SortLabel   string
	SortOptions []SortOption

	FilterSourceName string // "" if the feed isn't filtered to one source
	ClearFilterHref  string // href to remove the source filter (keeps the current sort/tag)

	TagFilters []shared.TagPill // "All" (TagID "") first, then one per tag

	TopbarUser shared.TopbarUser
}

// ArticleView is what web/components/article/row.html renders.
type ArticleView struct {
	Title      string
	Summary    string
	URL        string
	SourceName string
	Badge      shared.Badge
	TimeAgo    string
}

// SourceView is what web/components/source/item.html renders - the home feed's sources sidebar.
// Distinct from the source package's SourceRow, which is the fuller /sources table row.
type SourceView struct {
	Name        string
	Badge       shared.Badge
	Count       int    // total articles (the source's own denormalized count, not derived from rows)
	UnreadCount int    // unread articles actually in the DB for this source - 0 hides the badge
	Href        string // filters the home feed to this source; clears the filter if already active
	Active      bool
}

// SortOption is one entry in the sort dropdown menu - see web/pages/home.html.
type SortOption struct {
	Label  string
	Href   string
	Active bool
}
