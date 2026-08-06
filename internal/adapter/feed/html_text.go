package feed

import (
	"html"
	"regexp"
	"strings"
)

var htmlTagPattern = regexp.MustCompile(`<[^>]*>`)

// htmlToText strips HTML tags and decodes entities - RSS/Atom summaries routinely contain markup,
// but this app only ever shows Title/Summary as plain text.
func htmlToText(s string) string {
	s = htmlTagPattern.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)

	return strings.Join(strings.Fields(s), " ")
}
