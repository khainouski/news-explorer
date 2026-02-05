package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/khainouski/news-explorer/pkg/logger"
)

func hello(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello from News Explorer!"))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response")

		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remote_addr", r.RemoteAddr).
		Str("user_agent", r.UserAgent()).
		Msg("page visited")
}

func probe(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	logger.Init(logger.Config{Level: "info", PrettyConsole: true})

	router := chi.NewRouter()
	router.Get("/", hello)
	router.Get("/live", probe)
	router.Get("/ready", probe)

	log.Info().Msgf("Listening on port 8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Error().Err(err).Msg("Error starting server")
	}
}
