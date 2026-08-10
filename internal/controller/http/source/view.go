package source

import "github.com/khainouski/news-explorer/internal/controller/http/shared"

// SourcesView is what web/pages/sources.html renders.
type SourcesView struct {
	PageTitle     string
	Active        string
	LastSyncedAgo string
	SearchScope   string
	Sources       []SourceRow
	Total         int

	Query      string // current search term (?q=), matched against name/description/tag
	SourceHref string // "Source" header link - resets sort, keeps Query

	// SortKey/SortDir/TagID are the raw ?sort=/?dir=/?tag= values, picked up by the topbar's
	// search form via hx-include so searching doesn't drop them.
	SortKey string
	SortDir string
	TagID   string

	TagFilters []shared.TagPill // "All" (TagID "") first, then one per tag

	TagHeader      SortHeader
	StatusHeader   SortHeader
	ArticlesHeader SortHeader
	SyncedHeader   SortHeader

	TopbarUser shared.TopbarUser
}

// SourceRow is one row of web/pages/sources.html's table.
type SourceRow struct {
	ID           string
	Name         string
	FeedURL      string
	Description  string
	Tag          string
	Badge        shared.Badge
	Active       bool
	StatusLabel  string
	ArticleCount string
	SyncedAgo    string
	EditHref     string
}

// SortHeader drives one sortable column header. Dir is "asc"/"desc" if this column is the
// current sort, "" otherwise.
type SortHeader struct {
	Label string
	Href  string
	Dir   string
}

// SourceFormView is what web/pages/source_form.html renders - the Add/Edit Source form, both
// share this template. EditingID is "" for Add.
type SourceFormView struct {
	PageTitle     string
	Active        string
	LastSyncedAgo string
	Tags          []TagOption
	BadgeColors   []BadgeColor

	Title       string // "Add Source" or "Edit Source"
	SubmitLabel string // "Save Source" or "Save Changes"
	ActionURL   string // "/sources" (create) or "/sources/{id}" (update)
	EditingID   string // "" for Add; also drives the Danger Zone button's visibility

	Error  string
	Source SourceFormData

	Query       string // always "" - this page has no results to search
	SearchScope string // always "" - topbar.html falls back to its disabled-input branch

	TopbarUser shared.TopbarUser
}

type SourceFormData struct {
	Name          string
	FeedURL       string
	Description   string
	SelectedTag   string
	Badge         string
	SelectedColor string
	Status        string // "active" | "inactive"
}

// TagOption is one <option> in the "Add Source" form's tag dropdown.
type TagOption struct {
	Value string // domain.Tag.ID (the slug) - what actually gets submitted/stored
	Label string
}

// BadgeColor is one swatch in the "Add Source" form's badge color picker.
type BadgeColor struct {
	Value string // Tailwind background class, e.g. "bg-blue-500" - what actually gets stored
	Label string
}
