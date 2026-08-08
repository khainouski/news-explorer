// Package http wires the chi router - health probes and static assets outside the instrumented
// group, then the source/auth/article handler subpackages inside it.
package http

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/khainouski/news-explorer/internal/controller/http/article"
	"github.com/khainouski/news-explorer/internal/controller/http/auth"
	"github.com/khainouski/news-explorer/internal/controller/http/middleware"
	"github.com/khainouski/news-explorer/internal/controller/http/shared"
	"github.com/khainouski/news-explorer/internal/controller/http/source"
	usecasearticle "github.com/khainouski/news-explorer/internal/usecase/article"
	usecaseauth "github.com/khainouski/news-explorer/internal/usecase/auth"
	usecasesource "github.com/khainouski/news-explorer/internal/usecase/source"
	usecasesync "github.com/khainouski/news-explorer/internal/usecase/sync"
	usecasetag "github.com/khainouski/news-explorer/internal/usecase/tag"
	applogger "github.com/khainouski/news-explorer/pkg/logger"
	"github.com/khainouski/news-explorer/pkg/metrics"
	appotel "github.com/khainouski/news-explorer/pkg/otel"
	"github.com/khainouski/news-explorer/web"
)

// Dependencies are every use case plus Metrics, built once in internal/app/app.go.
type Dependencies struct {
	Article *usecasearticle.UseCase
	Source  *usecasesource.UseCase
	Tag     *usecasetag.UseCase
	Auth    *usecaseauth.UseCase
	Sync    *usecasesync.UseCase
	Metrics *metrics.HTTPServer
}

func NewRouter(d Dependencies) *chi.Mux {
	r := chi.NewRouter()

	// Health probes and /metrics stay outside the otel/logger/metrics stack below - Kubernetes
	// polls them every few seconds and they'd otherwise spam traces/logs/metrics with noise.
	r.Get("/live", probe)
	r.Get("/ready", probe)
	r.Handle("/metrics", promhttp.Handler())

	// Embedded (web/templates.go), not read from disk - works the same locally and in the
	// container, which ships only the compiled binary, not the web/ source tree.
	staticFS, err := fs.Sub(web.StaticFS, "static")
	if err != nil {
		panic(err) // static assets are embedded at build time - can only fail if the build itself is broken
	}

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	r.NotFound(shared.NotFound) // unmatched routes get the styled 404 page, not chi's plain-text default

	articleHandler := article.NewHandler(d.Article, d.Source, d.Tag)
	sourceHandler := source.NewHandler(d.Source, d.Tag, d.Sync)
	authHandler := auth.NewHandler(d.Auth)

	r.Group(func(r chi.Router) {
		r.Use(appotel.Middleware)
		r.Use(applogger.Middleware)
		r.Use(metrics.NewMiddleware(d.Metrics))
		r.Use(middleware.Auth(d.Auth))

		r.Get("/", articleHandler.List)

		r.Get("/sources", sourceHandler.List)
		r.With(middleware.RequireAdminPage).Get("/sources/new", sourceHandler.New)
		r.With(middleware.RequireAdminPage).Post("/sources", sourceHandler.Create)
		r.With(middleware.RequireAdminPage).Get("/sources/{id}/edit", sourceHandler.Edit)
		r.With(middleware.RequireAdminPage).Post("/sources/{id}", sourceHandler.Update)
		r.With(middleware.RequireAdminAPI).Delete("/sources/{id}", sourceHandler.Delete)
		r.With(middleware.RequireAdminAPI).Post("/sources/sync", sourceHandler.Sync)

		r.Get("/search", search)

		r.Get("/login", authHandler.Login)
		r.Post("/login", authHandler.LoginSubmit)
		r.Post("/logout", authHandler.Logout)
		r.Get("/account", authHandler.Account)
		r.Post("/account/password", authHandler.ChangePassword)
	})

	return r
}

// comingSoonView is what web/pages/coming_soon.html renders. Currently just /search.
type comingSoonView struct {
	PageTitle   string
	Active      string
	SearchScope string // always "" - no results container on this page, topbar shows a disabled input
	Description string
	TopbarUser  shared.TopbarUser
}

func search(w http.ResponseWriter, r *http.Request) {
	shared.Render(w, "coming_soon", comingSoonView{
		PageTitle:   "Search",
		Active:      "search",
		Description: "Search across every article and source.",
		TopbarUser:  shared.BuildTopbarUser(r),
	})
}

func probe(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
