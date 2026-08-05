package feed

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

func parseRSS(raw []byte) ([]Item, error) {
	var feed rssFeed
	if err := xml.Unmarshal(raw, &feed); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal rss: %w", err)
	}

	items := make([]Item, len(feed.Channel.Items))

	for i, ri := range feed.Channel.Items {
		items[i] = Item{
			Title:       strings.TrimSpace(ri.Title),
			URL:         strings.TrimSpace(ri.Link),
			Summary:     strings.TrimSpace(ri.Description),
			PublishedAt: parseDate(ri.PubDate),
			GUID:        strings.TrimSpace(ri.GUID),
		}
	}

	return items, nil
}
