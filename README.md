# News Explorer

Minimal Go HTTP service — the walking-skeleton first step of the deployment pipeline (Docker →
GHCR → Terraform/k3s → Helm → Argo CD, on DigitalOcean). No business logic yet on purpose: this
exists to prove the CI → registry → cluster chain end-to-end before any real functionality is
added.

## Endpoints

| Method | Path     | Purpose                   |
|--------|----------|----------------------------|
| GET    | `/`      | placeholder root response  |
| GET    | `/live`  | liveness probe             |
| GET    | `/ready` | readiness probe            |

## Development

```shell
make run           # go run ./cmd/app
make test          # go test -race -cover ./...
make build          # static binary in bin/
make docker-build   # build local image
make docker-run     # run local image, listens on :8080
```

Server listens on `:8080`.
