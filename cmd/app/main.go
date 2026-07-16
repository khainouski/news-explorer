package main

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/config"
	"github.com/khainouski/news-explorer/internal/app"
	"github.com/khainouski/news-explorer/pkg/logger"
	"github.com/khainouski/news-explorer/pkg/otel"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config.New")
	}

	logger.Init(c.Logger)

	ctx := context.Background()

	if err = otel.Init(ctx, c.OTEL); err != nil {
		log.Error().Err(err).Msg("Error initializing otel")
	}
	defer otel.Close()

	handler, pgPool, err := app.New(ctx, c)
	if err != nil {
		log.Fatal().Err(err).Msg("app.New")
	}
	defer pgPool.Close()

	log.Info().Msg("Listening on port 8080")

	if err = http.ListenAndServe(":8080", handler); err != nil {
		log.Error().Err(err).Msg("Error starting server")
	}
}
