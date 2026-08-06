package feed

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/khainouski/news-explorer/internal/domain"
)

// ToArticles converts one source's parsed feed items into domain.Article values ready to persist.
// Dedup is Postgres's job (see adapter/postgres/article.InsertBatch's ON CONFLICT). Items with no
// title or URL are dropped - Postgres requires both, and one bad item would fail the whole batch.
func ToArticles(sourceID string, items []Item) []domain.Article {
	articles := make([]domain.Article, 0, len(items))

	for _, item := range items {
		title := htmlToText(item.Title)
		url := strings.TrimSpace(item.URL)

		if title == "" || url == "" {
			continue
		}

		extID := externalID(item)

		articles = append(articles, domain.Article{
			ID:          articleID(sourceID, extID),
			SourceID:    sourceID,
			ExternalID:  extID,
			Title:       title,
			Summary:     htmlToText(item.Summary),
			URL:         url,
			PublishedAt: item.PublishedAt,
			Unread:      true,
		})
	}

	return articles
}

// articleID is a SHA-256 fingerprint of (sourceID, externalID) - deterministic, and never shown
// to anyone, so no need to slugify a messy title like Source.ID does.
func articleID(sourceID, externalID string) string {
	sum := sha256.Sum256([]byte(sourceID + "|" + externalID))

	return hex.EncodeToString(sum[:])
}

// externalID: Atom's <id>, then RSS's <guid>, then the URL (always present here - ToArticles
// already dropped anything without one).
func externalID(item Item) string {
	switch {
	case item.AtomID != "":
		return item.AtomID
	case item.GUID != "":
		return item.GUID
	default:
		return item.URL
	}
}
