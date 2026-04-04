package shared

// TagPill is one "Filter by tag" pill - the first is always "All". Shared between the article
// (home feed) and source (sources table) packages, which both offer the same filter.
type TagPill struct {
	Label  string
	Href   string
	Active bool
}
