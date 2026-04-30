package source

import (
	"errors"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/internal/controller/http/shared"
	"github.com/khainouski/news-explorer/internal/domain"
	usecasesource "github.com/khainouski/news-explorer/internal/usecase/source"
)

// badgeColorPalette are the swatches offered in the Add/Edit Source form (source_form.html) - a
// curated subset of the Tailwind colors already used across seeded sources' badges (see
// migration/postgres).
var badgeColorPalette = []BadgeColor{
	{Value: "bg-blue-500", Label: "Blue"},
	{Value: "bg-emerald-600", Label: "Green"},
	{Value: "bg-purple-600", Label: "Purple"},
	{Value: "bg-orange-500", Label: "Orange"},
	{Value: "bg-red-500", Label: "Red"},
	{Value: "bg-gray-900", Label: "Dark"},
	{Value: "bg-teal-600", Label: "Teal"},
}

// New renders the empty "Add Source" form - admin only (see the RequireAdminPage route
// middleware in router.go).
func (h *Handler) New(w http.ResponseWriter, r *http.Request) {
	h.renderForm(w, r, SourceFormView{
		Source: SourceFormData{
			SelectedColor: badgeColorPalette[0].Value,
			Status:        "active",
		},
	})
}

// Create handles the "Add Source" form submission - admin only (see the RequireAdminPage route
// middleware in router.go). On a validation or creation error, it re-renders the form with what
// the user typed still filled in, plus an error message, instead of a bare 400/500 page.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)

		return
	}

	view := parseFormView(r)

	status, ok := validateAndNormalize(&view)
	if !ok {
		h.renderForm(w, r, view)

		return
	}

	_, err := h.source.Create(r.Context(), usecasesource.CreateInput{
		Name:        view.Source.Name,
		FeedURL:     view.Source.FeedURL,
		Description: view.Source.Description,
		TagID:       view.Source.SelectedTag,
		Badge:       view.Source.Badge,
		BadgeColor:  view.Source.SelectedColor,
		Status:      status,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSourceExists):
			view.Error = "A source with this name already exists."
		case errors.Is(err, domain.ErrTagNotFound):
			view.Error = "Please choose a valid tag."
		case errors.Is(err, domain.ErrInvalidName):
			view.Error = "Name must contain at least one letter or number."
		default:
			log.Error().Err(err).Msg("create source")
			http.Error(w, "internal error", http.StatusInternalServerError)

			return
		}

		h.renderForm(w, r, view)

		return
	}

	http.Redirect(w, r, "/sources", http.StatusSeeOther)
}

// renderForm renders web/pages/source_form.html - shared by New/Create (empty/re-rendered-on-error)
// and Edit/Update (pre-filled) in edit.go. Title/SubmitLabel/ActionURL are computed here from
// EditingID rather than branched in the template.
func (h *Handler) renderForm(w http.ResponseWriter, r *http.Request, view SourceFormView) {
	tags, err := h.tag.List(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("list tags")
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	view.Title, view.SubmitLabel, view.ActionURL = "Add Source", "Save Source", "/sources"
	if view.EditingID != "" {
		view.Title, view.SubmitLabel, view.ActionURL = "Edit Source", "Save Changes", "/sources/"+view.EditingID
	}

	view.PageTitle = view.Title
	view.Active = "sources"
	view.Tags = toTagOptions(tags)
	view.BadgeColors = badgeColorPalette
	view.TopbarUser = shared.BuildTopbarUser(r)

	shared.Render(w, "source_form", view)
}

// parseFormView builds a SourceFormView from the submitted form - shared by Create and Update
// (which additionally sets EditingID after calling this).
func parseFormView(r *http.Request) SourceFormView {
	return SourceFormView{
		Source: SourceFormData{
			Name:          strings.TrimSpace(r.FormValue("name")),
			FeedURL:       strings.TrimSpace(r.FormValue("feed_url")),
			Description:   strings.TrimSpace(r.FormValue("description")),
			SelectedTag:   r.FormValue("tag"),
			Badge:         strings.TrimSpace(r.FormValue("badge")),
			SelectedColor: r.FormValue("badge_color"),
			Status:        r.FormValue("status"),
		},
	}
}

// validateAndNormalize checks the required fields (setting view.Error and returning false if any
// are missing) and fills in defaults for the rest - shared by Create and Update.
func validateAndNormalize(view *SourceFormView) (domain.SourceStatus, bool) {
	f := &view.Source

	if f.Name == "" || f.FeedURL == "" || f.SelectedTag == "" {
		view.Error = "Name, Feed URL and Tag are required."

		return "", false
	}

	if f.SelectedColor == "" {
		f.SelectedColor = badgeColorPalette[0].Value
	}

	status := domain.SourceStatusActive
	if f.Status == "inactive" {
		status = domain.SourceStatusInactive
	}

	return status, true
}

func toTagOptions(tags []domain.Tag) []TagOption {
	options := make([]TagOption, len(tags))
	for i, t := range tags {
		options[i] = TagOption{Value: t.ID, Label: t.Name}
	}

	return options
}
