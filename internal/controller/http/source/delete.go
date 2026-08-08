package source

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/internal/controller/http/shared"
	"github.com/khainouski/news-explorer/internal/domain"
)

// Delete removes a source - no confirmation step, deletes immediately.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := h.source.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrSourceNotFound) {
			shared.NotFound(w, r)

			return
		}

		log.Error().Err(err).Msg("delete source")
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
