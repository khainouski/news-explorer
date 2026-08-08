package logger

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"github.com/khainouski/news-explorer/pkg/router"
)

// Middleware logs one line per request after the handler runs, carrying the OTEL trace_id so
// Grafana can jump from a log line straight to its trace.
func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := router.WriterWrapper(w)
		next.ServeHTTP(ww, r)

		log.Info().
			Int("code", ww.Code()).
			Str("method", fmt.Sprintf("%s %s", r.Method, router.ExtractPath(r.Context()))).
			Str("remote_addr", router.ClientIP(r)).
			Str("user_agent", r.UserAgent()).
			Str("trace_id", trace.SpanContextFromContext(r.Context()).TraceID().String()).
			Msg("request handled")
	}

	return http.HandlerFunc(fn)
}
