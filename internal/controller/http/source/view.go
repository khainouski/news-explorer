package source

import "github.com/khainouski/news-explorer/internal/controller/http/shared"

// SourcesView is what web/pages/sources.html renders.
type SourcesView struct {
	PageTitle   string
	Active      string
	SearchScope string // "sources" - tells topbar.html which results container the search box targets
	Sources     []SourceRow
	Total       int

	Query      string // current search term (?q=), matched against name/description/tag
	SourceHref string // "Source" header link - resets sort, keeps Query

	// SortKey/SortDir/TagID are the raw ?sort=/?dir=/?tag= values (may be ""), exposed as hidden
	// inputs the topbar's search form picks up via hx-include so searching doesn't drop the
	// current sort/tag - see web/pages/sources.html and web/components/navigation/topbar.html.
	SortKey string
	SortDir string
	TagID   string

	// TagFilters are the "Filter by tag" pills - the first is always "All" (TagID "").
	TagFilters []shared.TagPill

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

// SortHeader drives one sortable column header (see web/components/shared/sortable_th.html): Href
// already points at the next click's target (toggling direction if this column is already
// active), Dir is "asc"/"desc" if this column is the current sort, "" otherwise.
type SortHeader struct {
	Label string
	Href  string
	Dir   string
}

// SourceFormView is what web/pages/source_form.html renders: the "Add Source"/"Edit Source" form
// - both share this template, distinguished by EditingID. Title/SubmitLabel/ActionURL are
// computed server-side (see source.Handler.renderForm) rather than branched in the template, so
// the template itself doesn't need to know Add and Edit are even different flows.
type SourceFormView struct {
	PageTitle   string
	Active      string // "sources" - keeps the Sources nav item highlighted
	Tags        []TagOption
	BadgeColors []BadgeColor

	Title       string // "Add Source" or "Edit Source"
	SubmitLabel string // "Save Source" or "Save Changes"
	ActionURL   string // "/sources" (create) or "/sources/{id}" (update)

	// EditingID is "" when adding a new source, or the source's ID when editing an existing one -
	// also drives the Danger Zone/delete button's visibility (Add has nothing to delete yet).
	EditingID string

	Error  string
	Source SourceFormData

	// Query is always "" and SearchScope is always "" - this page has no results to search, so
	// topbar.html falls back to its disabled-input branch (Active stays "sources" above so the
	// Sources nav item still highlights).
	Query       string
	SearchScope string

	TopbarUser shared.TopbarUser
}

// SourceFormData is the form's own field values - "" (Status: "active") on first load, whatever
// was submitted on a failed attempt, so the user doesn't have to retype everything.
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
	Label string // for the input's accessible name, e.g. "Blue"
}
