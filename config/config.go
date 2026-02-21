// Package config assembles every component's Config from the environment in one place.
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/khainouski/news-explorer/pkg/logger"
	"github.com/khainouski/news-explorer/pkg/otel"
	"github.com/khainouski/news-explorer/pkg/postgres"
)

type Config struct {
	Logger   logger.Config
	OTEL     otel.Config
	Postgres postgres.Config
}

// New loads .env if present (not required - variables can also be set directly in the
// environment) and fills Config from the environment.
func New() (Config, error) {
	var c Config

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return c, fmt.Errorf("godotenv.Load: %w", err)
	}

	if err := envconfig.Process("", &c); err != nil {
		return c, fmt.Errorf("envconfig.Process: %w", err)
	}

	return c, nil
}
