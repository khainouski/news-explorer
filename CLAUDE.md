# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

News Explorer is a minimal Go HTTP service — the walking-skeleton first step of a deployment
pipeline (Docker → GHCR → Terraform/k3s → Helm → Argo CD, on DigitalOcean). It intentionally has
no business logic yet: its purpose is to prove the CI → registry → cluster chain end-to-end before
real functionality is added. Expect the codebase to grow well beyond the skeleton described below.

## Commands

```shell
make run           # go run ./cmd/app — server listens on :8080
make build          # static binary (CGO_ENABLED=0) at bin/news-explorer
make test           # go test -race -cover ./...
make lint           # golangci-lint run ./...
make tidy           # go mod tidy
make docker-build    # build local image news-explorer:local
make docker-run      # run local image, published on :8080
```

Run a single test: `go test -race -run ^TestName$ ./path/to/package`.

CI (`.github/workflows/ci.yml`) runs on every PR/push to `main`: `go mod tidy` (fails the build if
go.mod/go.sum drift), `go vet ./...`, `golangci-lint run` (built from source via `install-mode:
goinstall`, since prebuilt binaries lag behind fresh Go releases), then
`go test -race -shuffle=on ./...`. Match these locally before pushing. The Docker image is built
and pushed to GHCR only after a merge to `main`, not on PRs.

## Architecture

- `cmd/app/main.go` — entrypoint. Wires up a `chi` router directly in `main()`; routes and
  handlers currently live inline here rather than in separate packages.
- `pkg/logger` — thin wrapper around `zerolog`, configured via `logger.Init(logger.Config{...})`
  once at startup. Handlers log through the global `github.com/rs/zerolog/log` logger rather than
  an injected logger instance. `logger.Middleware` logs one line per request (method, path, code,
  `trace_id`) — applied only to the instrumented route group, not `/live`/`/ready`/`/metrics`.
- `pkg/otel` — OpenTelemetry tracing. `otel.Init(ctx, otel.Config{Endpoint})` sends spans via
  OTLP/gRPC to Tempo when `OTEL_ENDPOINT` is set; empty endpoint falls back to a no-op tracer
  (`otel.SilentModeInit`), which is what local `make run` gets by default.
  `otel.Middleware` starts one root span per request, named after the matched route.
- `pkg/metrics` — Prometheus HTTP metrics (`http_server_requests_total`,
  `http_server_request_duration_seconds`), registered globally via `client_golang`. Scraped at
  `/metrics` (`promhttp.Handler()`); the Helm chart annotates the pod for the plain
  `prometheus-community/prometheus` chart's annotation-based discovery, no ServiceMonitor/operator
  needed.
- `pkg/router` — small shared helpers (`WriterWrapper` to read back the status code after a
  handler runs, `ExtractPath` to get chi's matched route pattern) used by all three middlewares
  above, so metrics/spans/logs key on the route pattern, not the raw literal URL.

Mirrors the instrumentation pattern from `user-service` (`pkg/otel`, `pkg/metrics`, `pkg/router`),
trimmed to what a single-route app needs — no `pkg/metrics/process.go`/pusher (no batch jobs here,
no Pushgateway in this project either) and no `envconfig`-based config loading (still just
`os.Getenv` in `main()`).

Routes today (`/`, `/live`, `/ready`, `/metrics`) — `/live`/`/ready` are placeholders/health probes
(return 200 with no readiness logic), deliberately kept outside the otel/logger/metrics middleware
group since Kubernetes polls them every few seconds and they'd otherwise spam traces/logs/metrics.
As real endpoints are added under the instrumented group in `main()`, keep that separation.

The Dockerfile is a two-stage build: `golang:1.25.6-alpine` compiles a static binary
(`-trimpath -ldflags="-s -w"`), then an `alpine:3.22` runtime image runs it as a non-root `app`
user. Go module downloads are cached in a separate layer from the source copy.
