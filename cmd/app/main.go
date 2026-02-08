package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/pkg/logger"
	"github.com/khainouski/news-explorer/pkg/metrics"
	"github.com/khainouski/news-explorer/pkg/otel"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello from News Explorer!"))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response")
	}
}

func probe(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Plain JSON by default (what Grafana Alloy/Loki expect from stdout in the cluster).
	// LOG_PRETTY=true switches to the human-readable console writer for local development.
	pretty := os.Getenv("LOG_PRETTY") == "true"
	logger.Init(logger.Config{Level: "info", PrettyConsole: pretty})

	// Empty OTEL_ENDPOINT (the default for local `make run`) falls back to a no-op tracer.
	if err := otel.Init(context.Background(), otel.Config{Endpoint: os.Getenv("OTEL_ENDPOINT")}); err != nil {
		log.Error().Err(err).Msg("Error initializing otel")
	}
	defer otel.Close()

	m := metrics.NewHTTPServer()

	router := chi.NewRouter()

	// Health probes and /metrics stay outside the otel/logger/metrics stack below - Kubernetes
	// polls them every few seconds and they'd otherwise spam traces/logs/metrics with noise.
	router.Get("/live", probe)
	router.Get("/ready", probe)
	router.Handle("/metrics", promhttp.Handler())

	router.Group(func(r chi.Router) {
		r.Use(otel.Middleware)
		r.Use(logger.Middleware)
		r.Use(metrics.NewMiddleware(m))

		r.Get("/", hello)
	})

	log.Info().Msgf("Listening on port 8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Error().Err(err).Msg("Error starting server")
	}
}
