package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	tracer_noop "go.opentelemetry.io/otel/trace/noop"

	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

type Config struct {
	// Endpoint is the Tempo OTLP/gRPC address (e.g. "tempo.monitoring.svc.cluster.local:4317").
	// Empty disables tracing (SilentModeInit) - the default for local `make run`.
	Endpoint string `envconfig:"OTEL_ENDPOINT"`
}

var shutdownTracing func(ctx context.Context) error

// SilentModeInit wires up a no-op tracer so span-producing code paths work identically whether
// or not a collector is configured, instead of branching on "is tracing enabled" everywhere.
func SilentModeInit() {
	otel.SetTracerProvider(tracer_noop.NewTracerProvider())
	tracer.Init(otel.Tracer(""))

	log.Info().Msg("otel: tracer is disabled")
}

func Init(ctx context.Context, c Config) error {
	if c.Endpoint == "" {
		SilentModeInit()

		return nil
	}

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(prop)

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(c.Endpoint), otlptracegrpc.WithInsecure())
	if err != nil {
		// Without this, tracer.tracer stays nil and every tracer.Start call panics.
		SilentModeInit()

		return fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("news-explorer"),
		)),
	)

	shutdownTracing = traceProvider.Shutdown

	otel.SetTracerProvider(traceProvider)
	tracer.Init(otel.Tracer(""))

	return nil
}

func Close() {
	if shutdownTracing == nil {
		return
	}

	err := shutdownTracing(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("otel: failed to shutdown tracing")
	}

	log.Info().Msg("otel: closed")
}
