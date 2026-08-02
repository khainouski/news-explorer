package feed

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/khainouski/news-explorer/internal/domain"
)

// ToArticles converts one source's parsed feed items into domain.Article values ready to persist.
// Dedup against what's already stored is PostgreSQL's job now (see
// adapter/postgres/article.InsertBatch's ON CONFLICT (source_id, external_id) DO NOTHING) - this
// just has to produce the same ExternalID for the same item on every re-sync.
//
// TODO: derive a stable slug ID (our own ID, not ExternalID) per article, same approach as
// adapter/postgres/source's name-to-slug.
func ToArticles(sourceID string, items []Item) []domain.Article {
	articles := make([]domain.Article, len(items))

	for i, item := range items {
		articles[i] = domain.Article{
			SourceID:    sourceID,
			ExternalID:  externalID(item),
			Title:       item.Title,
			Summary:     item.Summary,
			URL:         item.URL,
			PublishedAt: item.PublishedAt,
			Unread:      true,
		}
	}

	return articles
}

// externalID picks the most stable identifier available for item, in priority order:
//  1. Atom's own <id> - spec-required to be a stable, globally unique URI.
//  2. RSS's <guid> - stable by convention (most feeds don't reuse or change one), not guaranteed.
//  3. The article's URL - feeds without either of the above still almost always have a link.
//  4. A SHA-256 fingerprint of title+URL+published time, deterministic so the same item hashes
//     the same way on every re-sync - the last resort for a feed with none of the above.
func externalID(item Item) string {
	switch {
	case item.AtomID != "":
		return item.AtomID
	case item.GUID != "":
		return item.GUID
	case item.URL != "":
		return item.URL
	default:
		sum := sha256.Sum256([]byte(item.Title + "|" + item.URL + "|" + item.PublishedAt.String()))

		return hex.EncodeToString(sum[:])
	}
}
