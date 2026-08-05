package feed

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// Item is one entry parsed out of a feed - just the RSS2.0/Atom fields mapper.go cares about, not
// a full spec-accurate model of either format.
type Item struct {
	Title       string
	URL         string
	Summary     string
	PublishedAt time.Time

	// AtomID/GUID are the feed's own item identifiers, if the feed has one - see mapper.go's
	// ExternalID priority. RSS-specific and Atom-specific by nature (that's inherent to what they
	// are), which is exactly why they stay here in the adapter and never reach domain.Article.
	AtomID string
	GUID   string
}

// Parse reads raw (the bytes Client.Fetch returned) and extracts its entries - RSS2.0 (rss.go) or
// Atom (atom.go), picked by the document's root element, the one reliable way to tell the two
// apart without guessing from content.
func Parse(raw []byte) ([]Item, error) {
	root, err := rootElementName(raw)
	if err != nil {
		return nil, fmt.Errorf("determine feed format: %w", err)
	}

	switch root {
	case "rss":
		return parseRSS(raw)
	case "feed":
		return parseAtom(raw)
	default:
		return nil, fmt.Errorf("unrecognized feed format: root element <%s>", root)
	}
}

// rootElementName returns the document's outermost element name ("rss" or "feed") without
// decoding the rest of it - just enough to pick which format to fully unmarshal as.
func rootElementName(raw []byte) (string, error) {
	dec := xml.NewDecoder(bytes.NewReader(raw))

	for {
		tok, err := dec.Token()
		if err != nil {
			return "", fmt.Errorf("dec.Token: %w", err)
		}

		if se, ok := tok.(xml.StartElement); ok {
			return se.Name.Local, nil
		}
	}
}

// dateLayouts covers the date formats real-world feeds actually use: RFC1123Z/RFC1123 for RSS's
// pubDate (the spec says RFC822, but a numeric or named zone offset are both common), plus the
// same two with a non-zero-padded day ("Fri, 7 Aug 2026 ..." instead of "Fri, 07 Aug 2026 ...",
// which Go's reference-time parser otherwise rejects outright) - seen in the wild from Golang
// Weekly's own feed. RFC3339 and its nanosecond variant cover Atom's published/updated.
var dateLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 MST",
	time.RFC3339,
	time.RFC3339Nano,
}

// parseDate tries every layout in dateLayouts in turn, returning the zero time.Time if none match
// - an unparseable date on one item shouldn't fail the whole feed.
func parseDate(s string) time.Time {
	s = strings.TrimSpace(s)

	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}

	return time.Time{}
}
