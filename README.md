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

## Deployment (DigitalOcean)

Repositories involved:

- `news-explorer` (this repo) — app code, Dockerfile, CI (builds and pushes the image to GHCR)
- `news-explorer-infra` — Terraform (Droplet, Firewall, DNS) + cloud-init bootstrap (k3s, Argo CD)
- `news-platform-deploy` — Helm chart + Argo CD `Application` manifests (GitOps)

Steps:

1. Push to `main` — GitHub Actions builds and pushes the image to GHCR, tagged `sha-<commit>`.
2. `news-explorer-infra`: `terraform apply` provisions the Droplet; cloud-init installs k3s and
   bootstraps Argo CD on it.
3. In `news-platform-deploy`, set `image.tag` in `argocd/applications/news-explorer.yaml` to the
   new `sha-<commit>`, commit and push.
4. Argo CD picks up the change and rolls out the new version automatically.
5. Verify: `curl http://goskills.xyz`.

## Development

```shell
make run           # go run ./cmd/app
make test          # go test -race -cover ./...
make build          # static binary in bin/
make docker-build   # build local image
make docker-run     # run local image, listens on :8080
```

Server listens on `:8080`.
