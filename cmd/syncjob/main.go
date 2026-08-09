// Command syncjob runs one feed sync and exits - meant to be run as a Kubernetes CronJob, not as
// a long-running process (see cmd/app for that). A non-zero exit code means at least one source
// failed, so the CronJob's own success/failure tracking reflects the sync's real outcome.
package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/config"
	"github.com/khainouski/news-explorer/internal/adapter/feed"
	pgarticle "github.com/khainouski/news-explorer/internal/adapter/postgres/article"
	pgsource "github.com/khainouski/news-explorer/internal/adapter/postgres/source"
	usecasesync "github.com/khainouski/news-explorer/internal/usecase/sync"
	"github.com/khainouski/news-explorer/pkg/logger"
	"github.com/khainouski/news-explorer/pkg/otel"
	"github.com/khainouski/news-explorer/pkg/postgres"
)

func main() {
	if run() != nil {
		os.Exit(1)
	}
}

// run returns errors instead of calling log.Fatal/os.Exit directly, so the deferred Close calls
// below always run - os.Exit skips deferred functions.
func run() error {
	c, err := config.New()
	if err != nil {
		log.Error().Err(err).Msg("config.New")

		return err
	}

	logger.Init(c.Logger)

	ctx := context.Background()

	if err = otel.Init(ctx, c.OTEL); err != nil {
		log.Error().Err(err).Msg("Error initializing otel")
	}
	defer otel.Close()

	pgPool, err := postgres.New(ctx, c.Postgres)
	if err != nil {
		log.Error().Err(err).Msg("postgres.New")

		return err
	}
	defer pgPool.Close()

	sync := usecasesync.New(pgsource.NewRepo(pgPool), pgarticle.NewRepo(pgPool), feed.New())

	result, err := sync.Sync(ctx)

	log.Info().
		Int("sources_synced", result.SourcesSynced).
		Int("sources_failed", result.SourcesFailed).
		Int("articles_inserted", result.ArticlesInserted).
		Msg("sync complete")

	if err != nil {
		log.Error().Err(err).Msg("sync")
	}

	return err
}
