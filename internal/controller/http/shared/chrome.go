package shared

import (
	"net/http"

	"github.com/rs/zerolog/log"

	usecasesource "github.com/khainouski/news-explorer/internal/usecase/source"
)

// BuildChrome builds the two pieces of app shell every page but login/change_password/not_found
// needs - the topbar's account state and the sidebar's last-synced time - in one call instead of
// two.
func BuildChrome(r *http.Request, source *usecasesource.UseCase) (TopbarUser, string) {
	lastSyncedAt, err := source.LastSyncedAt(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("get last synced at")
	}

	return BuildTopbarUser(r), LastSyncedAgo(lastSyncedAt)
}
