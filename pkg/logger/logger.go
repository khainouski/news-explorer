package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Level         string
	PrettyConsole bool
}

func Init(c Config) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if level, err := zerolog.ParseLevel(c.Level); err == nil && level != zerolog.NoLevel {
		zerolog.SetGlobalLevel(level)
	}

	log.Logger = log.With().Logger()

	if c.PrettyConsole {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	}

	log.Info().Msg("logger initialized")
}
