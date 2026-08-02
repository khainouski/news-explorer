package feed

import "time"

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

// Parse reads raw (the bytes Client.Fetch returned) and extracts its entries.
//
// TODO: detect RSS2.0 vs Atom (root element name) and decode each via encoding/xml into Item -
// see migration/postgres's seed data for the shape existing articles already have.
func Parse(raw []byte) ([]Item, error) {
	return nil, nil
}
