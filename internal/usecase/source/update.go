package source

import (
	"context"
	"fmt"

	"github.com/khainouski/news-explorer/internal/domain"
)

// UpdateInput is what the "Edit Source" form submits. Unlike CreateInput it carries an ID - the
// source being edited, which never changes even if Name does.
type UpdateInput struct {
	ID          string
	Name        string
	FeedURL     string
	Description string
	TagID       string
	Badge       string
	BadgeColor  string
	Status      domain.SourceStatus
}

// Update updates every editable field of an existing source.
func (u *UseCase) Update(ctx context.Context, input UpdateInput) (domain.Source, error) {
	badge := input.Badge
	if badge == "" {
		badge = defaultBadge(input.Name)
	}

	s := domain.Source{
		ID:          input.ID,
		Name:        input.Name,
		FeedURL:     input.FeedURL,
		Description: input.Description,
		Tag:         domain.Tag{ID: input.TagID},
		Badge:       badge,
		BadgeColor:  input.BadgeColor,
		Status:      input.Status,
	}

	if err := u.postgres.Update(ctx, s); err != nil {
		return domain.Source{}, fmt.Errorf("postgres.Update: %w", err)
	}

	return s, nil
}
