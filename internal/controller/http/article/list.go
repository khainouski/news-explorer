package article

import (
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/internal/controller/http/shared"
	"github.com/khainouski/news-explorer/internal/domain"
)

const (
	sortLatest = "latest"
	sortOldest = "oldest"
)

// List handles GET / - the home feed.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sortBy := r.URL.Query().Get("sort")
	sourceFilter := r.URL.Query().Get("source")
	tagFilter := r.URL.Query().Get("tag")
	q := strings.TrimSpace(r.URL.Query().Get("q"))

	// Clicking a source in the sidebar both filters the feed to it and marks its articles read -
	// before List, so the fresh list already reflects it (sidebar counts).
	if sourceFilter != "" {
		if err := h.source.MarkRead(ctx, sourceFilter); err != nil {
			log.Error().Err(err).Msg("mark source read")
		}
	}

	// Articles, sources and tags are independent queries, but fetched sequentially rather than
	// concurrently: each is cheap against this app's small tables, and the connection pool is
	// only 10 wide (pkg/postgres.maxConns) - sequential holds one connection at a time per
	// request, concurrent would hold three, tripling how many simultaneous home-page hits it
	// takes to start queueing for a connection. Not worth the trade for a saving of a few
	// sub-millisecond queries.
	articles, err := h.article.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("list articles")
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

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

	sourceByID := sourceIndex(sources)
	unreadCounts := unreadCountsBySource(articles)

	if sourceFilter != "" {
		articles = filterArticlesBySource(articles, sourceFilter)
	}

	if tagFilter != "" {
		articles = filterArticlesByTag(articles, sourceByID, tagFilter)
	}

	if q != "" {
		articles = filterArticlesByQuery(articles, sourceByID, q)
	}

	sortArticles(articles, sortBy)

	view := HomeView{
		PageTitle:        "All Articles",
		Active:           "articles",
		SearchScope:      "articles",
		Articles:         toArticleViews(articles, sourceByID),
		Sources:          toSourceViews(sources, unreadCounts, sourceFilter, sortBy, tagFilter, q),
		TotalCount:       len(articles),
		SourceCount:      len(sources),
		Query:            q,
		SortBy:           sortBy,
		SourceID:         sourceFilter,
		TagID:            tagFilter,
		SortLabel:        sortLabel(sortBy),
		SortOptions:      sortOptions(sortBy, sourceFilter, tagFilter, q),
		FilterSourceName: sourceName(sources, sourceFilter),
		ClearFilterHref:  homeURL(sortBy, "", tagFilter, q),
		TagFilters:       tagFilters(tags, sortBy, sourceFilter, tagFilter, q),
		TopbarUser:       shared.BuildTopbarUser(r),
	}

	// Same HTMX partial-render pattern as /sources: the search box, sort dropdown and source
	// sidebar all swap #home-articles in place instead of navigating.
	if r.Header.Get("HX-Request") == "true" {
		shared.RenderBlock(w, "home", "home-articles", view)

		return
	}

	shared.Render(w, "home", view)
}

func sourceIndex(sources []domain.Source) map[string]domain.Source {
	m := make(map[string]domain.Source, len(sources))
	for _, s := range sources {
		m[s.ID] = s
	}

	return m
}

func filterArticlesBySource(articles []domain.Article, sourceID string) []domain.Article {
	filtered := make([]domain.Article, 0, len(articles))

	for _, a := range articles {
		if a.SourceID == sourceID {
			filtered = append(filtered, a)
		}
	}

	return filtered
}

// filterArticlesByTag keeps articles whose source has the given tag - Article has no tag of its
// own, it's reached transitively via its source (see internal/domain.Source.Tag).
func filterArticlesByTag(articles []domain.Article, sourceByID map[string]domain.Source, tagID string) []domain.Article {
	filtered := make([]domain.Article, 0, len(articles))

	for _, a := range articles {
		if sourceByID[a.SourceID].Tag.ID == tagID {
			filtered = append(filtered, a)
		}
	}

	return filtered
}

// filterArticlesByQuery keeps articles whose title, summary or source name contains q
// (case-insensitive) - mirrors the source package's filterSources.
func filterArticlesByQuery(articles []domain.Article, sourceByID map[string]domain.Source, q string) []domain.Article {
	q = strings.ToLower(q)

	filtered := make([]domain.Article, 0, len(articles))

	for _, a := range articles {
		match := strings.Contains(strings.ToLower(a.Title), q) ||
			strings.Contains(strings.ToLower(a.Summary), q) ||
			strings.Contains(strings.ToLower(sourceByID[a.SourceID].Name), q)

		if match {
			filtered = append(filtered, a)
		}
	}

	return filtered
}

func sourceName(sources []domain.Source, sourceID string) string {
	if sourceID == "" {
		return ""
	}

	for _, s := range sources {
		if s.ID == sourceID {
			return s.Name
		}
	}

	return ""
}

// sortArticles sorts in place. An empty/unknown sortBy leaves the existing (DB) order untouched
// - already newest-first, i.e. sortLatest.
func sortArticles(articles []domain.Article, sortBy string) {
	switch sortBy {
	case sortOldest:
		sort.SliceStable(articles, func(i, j int) bool {
			return articles[i].PublishedAt.Before(articles[j].PublishedAt)
		})
	}
}

func sortLabel(sortBy string) string {
	switch sortBy {
	case sortOldest:
		return "Oldest first"
	default:
		return "Latest first"
	}
}

func sortOptions(activeSort, sourceFilter, tagFilter, q string) []SortOption {
	if activeSort == "" {
		activeSort = sortLatest
	}

	options := []struct {
		key   string
		label string
	}{
		{sortLatest, "Latest first"},
		{sortOldest, "Oldest first"},
	}

	views := make([]SortOption, 0, len(options))
	for _, o := range options {
		views = append(views, SortOption{
			Label:  o.label,
			Href:   homeURL(o.key, sourceFilter, tagFilter, q),
			Active: o.key == activeSort,
		})
	}

	return views
}

func toSourceViews(sources []domain.Source, unreadCounts map[string]int, activeSourceID, sortBy, tagFilter, q string) []SourceView {
	views := make([]SourceView, 0, len(sources))

	for _, s := range sources {
		active := s.ID == activeSourceID

		href := homeURL(sortBy, s.ID, tagFilter, q)
		if active {
			href = homeURL(sortBy, "", tagFilter, q) // click the active source again to clear the filter
		}

		views = append(views, SourceView{
			Name:        s.Name,
			Badge:       shared.Badge{Text: s.Badge, Color: s.BadgeColor},
			Count:       s.ArticleCount,
			UnreadCount: unreadCounts[s.ID],
			Href:        href,
			Active:      active,
		})
	}

	return views
}

// tagFilters builds the "Filter by tag" pills - "All" first (clears the tag filter), then one
// per domain.Tag. Clicking a tag keeps the current sort/source/query, same as clicking a source
// or sort option.
func tagFilters(tags []domain.Tag, sortBy, sourceFilter, activeTag, q string) []shared.TagPill {
	pills := make([]shared.TagPill, 0, len(tags)+1)

	pills = append(pills, shared.TagPill{
		Label:  "All",
		Href:   homeURL(sortBy, sourceFilter, "", q),
		Active: activeTag == "",
	})

	for _, t := range tags {
		pills = append(pills, shared.TagPill{
			Label:  t.Name,
			Href:   homeURL(sortBy, sourceFilter, t.ID, q),
			Active: t.ID == activeTag,
		})
	}

	return pills
}

func unreadCountsBySource(articles []domain.Article) map[string]int {
	counts := make(map[string]int, len(articles))

	for _, a := range articles {
		if a.Unread {
			counts[a.SourceID]++
		}
	}

	return counts
}

// homeURL builds a "/" URL carrying sortBy/sourceID/tagID/q - any of them can be empty to omit
// it. sortLatest is the default, so it's never written out either.
func homeURL(sortBy, sourceID, tagID, q string) string {
	v := url.Values{}
	if sortBy != "" && sortBy != sortLatest {
		v.Set("sort", sortBy)
	}

	if sourceID != "" {
		v.Set("source", sourceID)
	}

	if tagID != "" {
		v.Set("tag", tagID)
	}

	if q != "" {
		v.Set("q", q)
	}

	if len(v) == 0 {
		return "/"
	}

	return "/?" + v.Encode()
}

func toArticleViews(articles []domain.Article, sourceByID map[string]domain.Source) []ArticleView {
	views := make([]ArticleView, 0, len(articles))
	for _, a := range articles {
		s := sourceByID[a.SourceID]

		views = append(views, ArticleView{
			Title:      a.Title,
			Summary:    a.Summary,
			URL:        a.URL,
			SourceName: s.Name,
			Badge:      shared.Badge{Text: s.Badge, Color: s.BadgeColor},
			TimeAgo:    shared.TimeAgo(a.PublishedAt),
		})
	}

	return views
}
