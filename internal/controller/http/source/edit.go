package source

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/internal/controller/http/shared"
	"github.com/khainouski/news-explorer/internal/domain"
	usecasesource "github.com/khainouski/news-explorer/internal/usecase/source"
)

// Edit renders the "Edit Source" form (same template as "Add Source" - see
// SourceFormView.EditingID), pre-filled with the source's current values - admin only (see the
// RequireAdminPage route middleware in router.go).
func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	s, err := h.source.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSourceNotFound) {
			shared.NotFound(w, r)

			return
		}

		log.Error().Err(err).Msg("get source")
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	status := "active"
	if s.Status == domain.SourceStatusInactive {
		status = "inactive"
	}

	h.renderForm(w, r, SourceFormView{
		EditingID: s.ID,
		Source: SourceFormData{
			Name:          s.Name,
			FeedURL:       s.FeedURL,
			Description:   s.Description,
			SelectedTag:   s.Tag.ID,
			Badge:         s.Badge,
			SelectedColor: s.BadgeColor,
			Status:        status,
		},
	})
}

// Update handles the "Edit Source" form submission - admin only (see the RequireAdminPage route
// middleware in router.go), same validation/re-render-on-error shape as Create.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)

		return
	}

	view := parseFormView(r)
	view.EditingID = id

	status, ok := validateAndNormalize(&view)
	if !ok {
		h.renderForm(w, r, view)

		return
	}

	_, err := h.source.Update(r.Context(), usecasesource.UpdateInput{
		ID:          id,
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
		case errors.Is(err, domain.ErrTagNotFound):
			view.Error = "Please choose a valid tag."
		case errors.Is(err, domain.ErrSourceNotFound):
			shared.NotFound(w, r)

			return
		default:
			log.Error().Err(err).Msg("update source")
			http.Error(w, "internal error", http.StatusInternalServerError)

			return
		}

		h.renderForm(w, r, view)

		return
	}

	http.Redirect(w, r, "/sources", http.StatusSeeOther)
}
