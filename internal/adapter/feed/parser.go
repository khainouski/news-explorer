package feed

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// Item is one entry parsed out of a feed - just the RSS2.0/Atom fields mapper.go cares about.
type Item struct {
	Title       string
	URL         string
	Summary     string
	PublishedAt time.Time
	AtomID      string
	GUID        string
}

// Parse picks RSS2.0 (rss.go) or Atom (atom.go) by the document's root element.
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

// dateLayouts includes non-zero-padded-day variants ("Fri, 7 Aug ..." not "Fri, 07 Aug ...") -
// Go's reference-time parser otherwise rejects them outright, and real feeds send them.
var dateLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 MST",
	time.RFC3339,
	time.RFC3339Nano,
}

func parseDate(s string) time.Time {
	s = strings.TrimSpace(s)

	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}

	return time.Time{}
}
