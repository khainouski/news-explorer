package feed

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Title     string     `xml:"title"`
	ID        string     `xml:"id"`
	Summary   string     `xml:"summary"`
	Content   string     `xml:"content"`
	Published string     `xml:"published"`
	Updated   string     `xml:"updated"`
	Links     []atomLink `xml:"link"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

func parseAtom(raw []byte) ([]Item, error) {
	var feed atomFeed
	if err := xml.Unmarshal(raw, &feed); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal atom: %w", err)
	}

	items := make([]Item, len(feed.Entries))

	for i, e := range feed.Entries {
		summary := e.Summary
		if summary == "" {
			summary = e.Content
		}

		// published is when the entry was first posted; updated is required by the Atom spec and
		// always present, published isn't - fall back to it so PublishedAt is rarely left zero.
		published := e.Published
		if published == "" {
			published = e.Updated
		}

		items[i] = Item{
			Title:       strings.TrimSpace(e.Title),
			URL:         atomLinkURL(e.Links),
			Summary:     strings.TrimSpace(summary),
			PublishedAt: parseDate(published),
			AtomID:      strings.TrimSpace(e.ID),
		}
	}

	return items, nil
}

// atomLinkURL picks the entry's own page out of its <link> elements - rel="alternate" (or no rel
// at all, which the Atom spec says defaults to "alternate") is the human-readable page; anything
// else (rel="self", "enclosure", ...) points somewhere else entirely. Falls back to the first
// link at all if none is explicitly "alternate".
func atomLinkURL(links []atomLink) string {
	for _, l := range links {
		if l.Rel == "" || l.Rel == "alternate" {
			return l.Href
		}
	}

	if len(links) > 0 {
		return links[0].Href
	}

	return ""
}
