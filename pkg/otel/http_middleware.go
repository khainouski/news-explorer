package otel

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"

	"github.com/khainouski/news-explorer/pkg/otel/tracer"
	"github.com/khainouski/news-explorer/pkg/router"
)

// Middleware starts one root span per request (or continues one propagated in via headers),
// named after the matched route so spans for "/users/{id}" don't fragment into one series per ID.
func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		ctx, span := tracer.Start(ctx, "", trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		ww := router.WriterWrapper(w)

		next.ServeHTTP(ww, r.WithContext(ctx))

		span.SetName("http " + r.Method + " " + router.ExtractPath(ctx))

		span.SetAttributes(
			semconv.HTTPResponseStatusCode(ww.Code()),
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.HTTPRoute(r.URL.Path),
		)

		if ww.Code() >= http.StatusBadRequest {
			span.SetStatus(codes.Error, http.StatusText(ww.Code()))
			span.AddEvent("error", trace.WithAttributes(
				attribute.String("error.message", http.StatusText(ww.Code())),
			))
		}
	}

	return http.HandlerFunc(fn)
}
