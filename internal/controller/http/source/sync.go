package source

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// syncResponse is the JSON body for the "Sync complete" toast (see web/static/js/sources.js).
type syncResponse struct {
	SourcesSynced    int `json:"sourcesSynced"`
	SourcesFailed    int `json:"sourcesFailed"`
	ArticlesInserted int `json:"articlesInserted"`
}

// Sync triggers a full sync run - admin only. Runs synchronously; a partial failure still counts
// as a completed run, logged here rather than surfaced as an HTTP error.
func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	result, err := h.sync.Sync(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("sync sources")
	}

	log.Info().
		Int("sources_synced", result.SourcesSynced).
		Int("sources_failed", result.SourcesFailed).
		Int("articles_inserted", result.ArticlesInserted).
		Msg("sync complete")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(syncResponse{
		SourcesSynced:    result.SourcesSynced,
		SourcesFailed:    result.SourcesFailed,
		ArticlesInserted: result.ArticlesInserted,
	})
}
