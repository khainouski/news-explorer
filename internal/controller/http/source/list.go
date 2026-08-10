package source

import (
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/internal/controller/http/shared"
	"github.com/khainouski/news-explorer/internal/domain"
)

// List handles GET /sources - the sources table with search/sort/tag filter.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sources, err := h.source.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("list sources")
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	tags, err := h.tag.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("list tags")
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		sources = filterSources(sources, q)
	}

	tagFilter := r.URL.Query().Get("tag")
	if tagFilter != "" {
		sources = filterSourcesByTag(sources, tagFilter)
	}

	sortBy := r.URL.Query().Get("sort")
	desc := r.URL.Query().Get("dir") == "desc"
	sortSources(sources, sortBy, desc)

	dir := "asc"
	if desc {
		dir = "desc"
	}

	view := SourcesView{
		PageTitle:      "Sources",
		Active:         "sources",
		SearchScope:    "sources",
		Sources:        toSourceRows(sources),
		Total:          len(sources),
		Query:          q,
		SourceHref:     sourcesURL(q, "", "", tagFilter),
		SortKey:        sortBy,
		SortDir:        dir,
		TagID:          tagFilter,
		TagFilters:     sourceTagFilters(tags, q, sortBy, dir, tagFilter),
		TagHeader:      sortHeader("Tag", "tag", sortBy, dir, q, tagFilter),
		StatusHeader:   sortHeader("Status", "status", sortBy, dir, q, tagFilter),
		ArticlesHeader: sortHeader("Articles", "articles", sortBy, dir, q, tagFilter),
		SyncedHeader:   sortHeader("Last Synced", "synced", sortBy, dir, q, tagFilter),
	}

	if r.Header.Get("HX-Request") == "true" {
		shared.RenderBlock(w, "sources", "sources-results", view)

		return
	}

	// Only the full page renders the topbar/sidebar, so this is skipped above for the
	// "sources-results" HTMX partial (search/sort/filter) - no point building it there.
	view.TopbarUser, view.LastSyncedAgo = shared.BuildChrome(r, h.source)

	shared.Render(w, "sources", view)
}

// filterSources keeps sources whose name, description or tag contains q (case-insensitive).
func filterSources(sources []domain.Source, q string) []domain.Source {
	q = strings.ToLower(q)

	filtered := make([]domain.Source, 0, len(sources))

	for _, s := range sources {
		match := strings.Contains(strings.ToLower(s.Name), q) ||
			strings.Contains(strings.ToLower(s.Description), q) ||
			strings.Contains(strings.ToLower(s.Tag.Name), q)

		if match {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

func filterSourcesByTag(sources []domain.Source, tagID string) []domain.Source {
	filtered := make([]domain.Source, 0, len(sources))

	for _, s := range sources {
		if s.Tag.ID == tagID {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

func sourceTagFilters(tags []domain.Tag, q, sortBy, dir, activeTag string) []shared.TagPill {
	pills := make([]shared.TagPill, 0, len(tags)+1)

	pills = append(pills, shared.TagPill{
		Label:  "All",
		Href:   sourcesURL(q, sortBy, dir, ""),
		Active: activeTag == "",
	})

	for _, t := range tags {
		pills = append(pills, shared.TagPill{
			Label:  t.Name,
			Href:   sourcesURL(q, sortBy, dir, t.ID),
			Active: t.ID == activeTag,
		})
	}

	return pills
}

func sortSources(sources []domain.Source, by string, desc bool) {
	var less func(i, j int) bool

	switch by {
	case "tag":
		less = func(i, j int) bool {
			if sources[i].Tag.Name != sources[j].Tag.Name {
				return lessStr(sources[i].Tag.Name, sources[j].Tag.Name, desc)
			}

			return sources[i].Name < sources[j].Name
		}
	case "status":
		less = func(i, j int) bool {
			if sources[i].Status != sources[j].Status {
				return lessStr(string(sources[i].Status), string(sources[j].Status), desc)
			}

			return sources[i].Name < sources[j].Name
		}
	case "articles":
		less = func(i, j int) bool {
			if sources[i].ArticleCount != sources[j].ArticleCount {
				return lessInt(sources[i].ArticleCount, sources[j].ArticleCount, desc)
			}

			return sources[i].Name < sources[j].Name
		}
	case "synced":
		less = func(i, j int) bool {
			ti, tj := syncedAt(sources[i]), syncedAt(sources[j])
			if !ti.Equal(tj) {
				if desc {
					return ti.After(tj)
				}

				return ti.Before(tj)
			}

			return sources[i].Name < sources[j].Name
		}
	default:
		return
	}

	sort.SliceStable(sources, less)
}

func syncedAt(s domain.Source) time.Time {
	if s.LastSyncedAt == nil {
		return time.Time{}
	}

	return *s.LastSyncedAt
}

func lessStr(a, b string, desc bool) bool {
	if desc {
		return a > b
	}

	return a < b
}

func lessInt(a, b int, desc bool) bool {
	if desc {
		return a > b
	}

	return a < b
}

// sortHeader builds one SortHeader: Href toggles direction if key is already the active sort,
// otherwise starts it ascending. q/tag are carried along so search/tag survive a sort click.
func sortHeader(label, key, activeSort, activeDir, q, tag string) SortHeader {
	dir := "asc"

	var currentDir string

	if activeSort == key {
		currentDir = activeDir

		if activeDir == "asc" {
			dir = "desc"
		}
	}

	return SortHeader{
		Label: label,
		Href:  sourcesURL(q, key, dir, tag),
		Dir:   currentDir,
	}
}

func sourcesURL(q, sortBy, dir, tag string) string {
	v := url.Values{}
	if q != "" {
		v.Set("q", q)
	}

	if sortBy != "" {
		v.Set("sort", sortBy)
		v.Set("dir", dir)
	}

	if tag != "" {
		v.Set("tag", tag)
	}

	if len(v) == 0 {
		return "/sources"
	}

	return "/sources?" + v.Encode()
}

func toSourceRows(sources []domain.Source) []SourceRow {
	rows := make([]SourceRow, 0, len(sources))

	for _, s := range sources {
		active := s.Status == domain.SourceStatusActive

		row := SourceRow{
			ID:           s.ID,
			Name:         s.Name,
			FeedURL:      s.FeedURL,
			Description:  s.Description,
			Tag:          s.Tag.Name,
			Badge:        shared.Badge{Text: s.Badge, Color: s.BadgeColor},
			Active:       active,
			StatusLabel:  "Inactive",
			ArticleCount: "—",
			SyncedAgo:    "—",
			EditHref:     "/sources/" + s.ID + "/edit",
		}

		if active {
			row.StatusLabel = "Active"
			row.ArticleCount = formatCount(s.ArticleCount)
		}

		if s.LastSyncedAt != nil {
			row.SyncedAgo = shared.TimeAgo(*s.LastSyncedAt)
		}

		rows = append(rows, row)
	}

	return rows
}
