package source

import (
	"context"
	"fmt"
	"strings"

	"github.com/khainouski/news-explorer/internal/domain"
)

// CreateInput is what the "Add Source" form submits.
type CreateInput struct {
	Name        string
	FeedURL     string
	Description string
	TagID       string
	Badge       string
	BadgeColor  string
	Status      domain.SourceStatus
}

// Create adds a new global source (see domain.Source.UserID - no auth-scoped creation flow yet,
// so every source created this way is visible to everyone). The ID is a slug of the name; two
// sources with the same name collide (domain.ErrSourceExists) rather than silently overwriting
// one another.
func (u *UseCase) Create(ctx context.Context, input CreateInput) (domain.Source, error) {
	id := slugify(input.Name)
	if id == "" {
		return domain.Source{}, domain.ErrInvalidName
	}

	badge := input.Badge
	if badge == "" {
		badge = defaultBadge(input.Name)
	}

	s := domain.Source{
		ID:          id,
		Name:        input.Name,
		FeedURL:     input.FeedURL,
		Description: input.Description,
		Tag:         domain.Tag{ID: input.TagID},
		Badge:       badge,
		BadgeColor:  input.BadgeColor,
		Status:      input.Status,
	}

	if err := u.postgres.Create(ctx, s); err != nil {
		return domain.Source{}, fmt.Errorf("postgres.Create: %w", err)
	}

	return s, nil
}

// slugify turns a name into a URL/ID-friendly slug: lowercase, non-alphanumeric runs collapsed
// to a single hyphen, matching the style of the migration-seeded IDs (e.g. "go-blog").
func slugify(name string) string {
	var b strings.Builder

	prevHyphen := true // true so a leading run of non-alphanumerics doesn't start with a hyphen

	for _, r := range strings.ToLower(name) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)

			prevHyphen = false
		case !prevHyphen:
			b.WriteByte('-')

			prevHyphen = true
		}
	}

	return strings.TrimSuffix(b.String(), "-")
}

// defaultBadge derives a short badge label when the form leaves it blank: initials of the first
// two words, or the first two letters of a single-word name.
func defaultBadge(name string) string {
	words := strings.Fields(name)

	switch {
	case len(words) >= 2:
		return strings.ToUpper(string([]rune(words[0])[:1]) + string([]rune(words[1])[:1]))
	case len(words) == 1 && len([]rune(words[0])) >= 2:
		return strings.ToUpper(string([]rune(words[0])[:2]))
	case len(words) == 1:
		return strings.ToUpper(words[0])
	default:
		return "?"
	}
}
