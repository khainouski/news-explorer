// Package app wires the application together: config, adapters, usecases, controllers, router.
// Kept separate from cmd/app/main.go so main() stays a thin process-lifecycle shell.
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/khainouski/news-explorer/config"
	"github.com/khainouski/news-explorer/internal/adapter/feed"
	pgarticle "github.com/khainouski/news-explorer/internal/adapter/postgres/article"
	pgsession "github.com/khainouski/news-explorer/internal/adapter/postgres/session"
	pgsource "github.com/khainouski/news-explorer/internal/adapter/postgres/source"
	pgtag "github.com/khainouski/news-explorer/internal/adapter/postgres/tag"
	pguser "github.com/khainouski/news-explorer/internal/adapter/postgres/user"
	httpcontroller "github.com/khainouski/news-explorer/internal/controller/http"
	usecasearticle "github.com/khainouski/news-explorer/internal/usecase/article"
	usecaseauth "github.com/khainouski/news-explorer/internal/usecase/auth"
	usecasesource "github.com/khainouski/news-explorer/internal/usecase/source"
	usecasesync "github.com/khainouski/news-explorer/internal/usecase/sync"
	usecasetag "github.com/khainouski/news-explorer/internal/usecase/tag"
	"github.com/khainouski/news-explorer/pkg/metrics"
	"github.com/khainouski/news-explorer/pkg/postgres"
)

// New builds the full HTTP handler for the application, plus the Postgres pool backing it - the
// caller owns the pool's lifecycle (Close it on shutdown).
func New(ctx context.Context, c config.Config) (http.Handler, *postgres.Pool, error) {
	pgPool, err := postgres.New(ctx, c.Postgres)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres.New: %w", err)
	}

	m := metrics.NewHTTPServer()

	// Shared between usecase/source and usecase/sync - one repo per table, not per usecase.
	sourceRepo := pgsource.NewRepo(pgPool)
	articleRepo := pgarticle.NewRepo(pgPool)

	deps := httpcontroller.Dependencies{
		Article: usecasearticle.New(articleRepo),
		Source:  usecasesource.New(sourceRepo),
		Tag:     usecasetag.New(pgtag.NewRepo(pgPool)),
		Auth:    usecaseauth.New(pguser.NewRepo(pgPool), pgsession.NewRepo(pgPool)),
		Sync:    usecasesync.New(sourceRepo, articleRepo, feed.New()),
		Metrics: m,
	}

	handler := httpcontroller.NewRouter(deps)

	return handler, pgPool, nil
}
