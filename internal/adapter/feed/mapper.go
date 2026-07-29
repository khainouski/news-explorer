package feed

import "github.com/khainouski/news-explorer/internal/domain"

// ToArticles converts one source's parsed feed items into domain.Article values ready to persist.
//
// TODO: derive a stable slug ID per article (same approach as adapter/postgres/source's
// name-to-slug), and let the caller (the sync usecase) dedupe against what's already stored so a
// re-sync doesn't recreate articles that already exist.
func ToArticles(sourceID string, items []Item) []domain.Article {
	return nil
}
