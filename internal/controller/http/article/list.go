package article

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/internal/controller/http/middleware"
	"github.com/khainouski/news-explorer/internal/controller/http/shared"
	"github.com/khainouski/news-explorer/internal/domain"
)

const (
	sortLatest = "latest"
	sortOldest = "oldest"

	pageSize = 20
)

// List handles GET / - the home feed.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sortBy := r.URL.Query().Get("sort")
	sourceFilter := r.URL.Query().Get("source")
	tagFilter := r.URL.Query().Get("tag")
	q := strings.TrimSpace(r.URL.Query().Get("q"))

	// Clicking a source marks its articles read before List, so the fresh list reflects it.
	// Admin-only: unread is a single shared flag (no per-user tracking), so an anonymous visitor
	// clicking a source must not clear it for everyone else.
	if user := middleware.CurrentUser(ctx); sourceFilter != "" && user != nil && user.IsAdmin() {
		if err := h.source.MarkRead(ctx, sourceFilter); err != nil {
			log.Error().Err(err).Msg("mark source read")
		}
	}

	page := currentPage(r)
	isHXRequest := r.Header.Get("HX-Request") == "true"
	isLoadMore := page > 1 && isHXRequest

	params := domain.ArticleListParams{
		SourceID: sourceFilter,
		TagID:    tagFilter,
		Query:    q,
		Oldest:   sortBy == sortOldest,
		Limit:    pageSize,
		Offset:   (page - 1) * pageSize,
	}

	// Sequential, not concurrent: cheap queries, and the connection pool is only 10 wide
	// (pkg/postgres.maxConns) - not worth 3x the connections for microseconds saved.
	articles, totalCount, err := h.article.List(ctx, params)
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

	sourceByID := sourceIndex(sources)
	hasMore := params.Offset+len(articles) < totalCount

	var nextPageHref string
	if hasMore {
		nextPageHref = homeURL(sortBy, sourceFilter, tagFilter, q, page+1)
	}

	if isLoadMore {
		shared.RenderBlock(w, "home", "article-page", HomeView{
			Articles:     toArticleViews(articles, sourceByID),
			NextPageHref: nextPageHref,
		})

		return
	}

	tags, err := h.tag.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("list tags")
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	unreadCounts, err := h.article.UnreadCountsBySource(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unread counts")
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	view := HomeView{
		PageTitle:        "All Articles",
		Active:           "articles",
		SearchScope:      "articles",
		Articles:         toArticleViews(articles, sourceByID),
		Sources:          toSourceViews(sources, unreadCounts, sourceFilter, sortBy, tagFilter, q),
		TotalCount:       totalCount,
		SourceCount:      len(sources),
		Query:            q,
		SortBy:           sortBy,
		SourceID:         sourceFilter,
		TagID:            tagFilter,
		SortLabel:        sortLabel(sortBy),
		SortOptions:      sortOptions(sortBy, sourceFilter, tagFilter, q),
		FilterSourceName: sourceName(sources, sourceFilter),
		ClearFilterHref:  homeURL(sortBy, "", tagFilter, q, 0),
		TagFilters:       tagFilters(tags, sortBy, sourceFilter, tagFilter, q),
		NextPageHref:     nextPageHref,
	}

	if isHXRequest {
		shared.RenderBlock(w, "home", "home-articles", view)

		return
	}

	view.TopbarUser, view.LastSyncedAgo = shared.BuildChrome(r, h.source)

	shared.Render(w, "home", view)
}

func currentPage(r *http.Request) int {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		return 1
	}

	return page
}

func sourceIndex(sources []domain.Source) map[string]domain.Source {
	m := make(map[string]domain.Source, len(sources))
	for _, s := range sources {
		m[s.ID] = s
	}

	return m
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
			Href:   homeURL(o.key, sourceFilter, tagFilter, q, 0),
			Active: o.key == activeSort,
		})
	}

	return views
}

func toSourceViews(sources []domain.Source, unreadCounts map[string]int, activeSourceID, sortBy, tagFilter, q string) []SourceView {
	views := make([]SourceView, 0, len(sources))

	for _, s := range sources {
		active := s.ID == activeSourceID

		href := homeURL(sortBy, s.ID, tagFilter, q, 0)
		if active {
			href = homeURL(sortBy, "", tagFilter, q, 0)
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

func tagFilters(tags []domain.Tag, sortBy, sourceFilter, activeTag, q string) []shared.TagPill {
	pills := make([]shared.TagPill, 0, len(tags)+1)

	pills = append(pills, shared.TagPill{
		Label:  "All",
		Href:   homeURL(sortBy, sourceFilter, "", q, 0),
		Active: activeTag == "",
	})

	for _, t := range tags {
		pills = append(pills, shared.TagPill{
			Label:  t.Name,
			Href:   homeURL(sortBy, sourceFilter, t.ID, q, 0),
			Active: t.ID == activeTag,
		})
	}

	return pills
}

// homeURL builds a "/" URL carrying sortBy/sourceID/tagID/q/page - any can be zero-valued to
// omit it. sortLatest and page<=1 are the defaults, so neither is ever written out.
func homeURL(sortBy, sourceID, tagID, q string, page int) string {
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

	if page > 1 {
		v.Set("page", strconv.Itoa(page))
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
