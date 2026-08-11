# News Explorer

**News Explorer** ("Dev News") is a backend-focused Go project: a server-rendered news aggregator
built with clean-architecture layering (`domain → usecase → adapter → controller`), `pgx` against
Postgres with no ORM, `chi` for routing, and OpenTelemetry/Prometheus instrumentation throughout.
Feed syncing is hand-rolled end to end — no third-party RSS/Atom library — with concurrent
fetching over a worker pool and Postgres-level dedup. The frontend is deliberately thin
(`html/template` + HTMX + Tailwind CSS, no SPA, no JS framework) so the Go backend stays the
actual subject of the project, not scaffolding around a frontend. Content-wise, it aggregates
Go-developer news: articles are pulled live from a curated set of RSS/Atom feeds — Go blogs,
newsletters, podcasts, and remote-job boards — deduplicated, and stored in Postgres, served as a
fast feed with sorting, tag filtering, and search.

The whole thing doubles as an end-to-end deployment pipeline: Docker → GHCR → Terraform/k3s →
Helm → Argo CD, running on DigitalOcean with a full observability stack (Prometheus, Tempo,
Grafana). See [Cloud Deployment](#cloud-deployment) and [Architecture](#architecture) below.

## Local Setup

**Prerequisites:** Go 1.26.5, Docker (for `docker compose`), the
[`golang-migrate`](https://github.com/golang-migrate/migrate) CLI (`make migrate-install`), and
the [Tailwind CLI](https://github.com/tailwindlabs/tailwindcss/releases) (`make tailwind-install`).

```shell
cp .env.example .env    # Postgres connection vars - defaults match docker-compose.yml
make up                 # start Postgres in Docker
make migrate-up         # apply schema + seed data (Go sources, tags, the admin user)
make tailwind-build      # compile web/static/css/app.css - re-run after changing template classes
make run                # go run ./cmd/app - server listens on :8080
```

Open **http://localhost:8080**.

To stop: `Ctrl+C` the `make run` process, then `make down` to stop and remove the Postgres
container (data is not persisted across `make down`/`make up` — re-run `make migrate-up`
afterward).

### Signing in as admin

Adding, editing, or deleting a source, and triggering a sync, all require the admin account —
everything else (browsing, filtering, search) is open to anyone.

1. Go to `/login`.
2. Sign in with login **`admin`**, password **`12345`** (seeded, deliberately weak — change it
   from `/account` once logged in, though nothing forces you to).

### Running a sync

Articles aren't live until a sync has pulled them in from each source's feed:

1. While signed in as admin, go to `/sources`.
2. Click the **"Sync"** button (top right). It fetches every active source's RSS/Atom feed
   concurrently, inserts new articles, and shows a toast summarizing what happened
   (sources synced/failed, articles inserted) once it's done.

Syncs also happen automatically in the cloud deployment on a schedule — see
[Architecture](#architecture).

## Cloud Deployment

Running News Explorer in production spans three repositories:

| Repository | Role |
|---|---|
| **`news-explorer`** (this repo) | App code, Dockerfile, CI — builds and pushes the image to GHCR |
| **`news-explorer-infra`** | Terraform — provisions the DigitalOcean Droplet, firewall, and DNS; cloud-init bootstraps k3s and Argo CD on it |
| **`news-platform-deploy`** | Helm chart + Argo CD `Application` manifests (GitOps) — deploys the app and the observability stack (Prometheus, Tempo, Grafana) onto the cluster |
| **`news-explorer-client`** | Synthetic load generator — drives realistic traffic against the deployed app so Grafana dashboards have real data to show |

**Release flow:**

1. Push to `main` — GitHub Actions builds and pushes the image to GHCR, tagged `sha-<commit>`.
2. `news-explorer-infra`: `terraform apply` provisions the Droplet; cloud-init installs k3s and
   bootstraps Argo CD on it (one-time, or whenever infra changes).
3. In `news-platform-deploy`, bump `image.tag` in `argocd/applications/news-explorer.yaml` to the
   new `sha-<commit>`, commit and push.
4. Argo CD picks up the change and rolls out the new version automatically.
5. Verify: `curl http://goskills.xyz`.

| URL | What |
|---|---|
| https://goskills.xyz | the app |
| https://grafana.goskills.xyz | Grafana dashboards |
| https://argocd.goskills.xyz | Argo CD UI |

Grafana and Argo CD both need a password — read it out of the cluster (`admin`/`admin-user` is
the username in both cases):

```bash
DROPLET_IP=$(cd ../news-explorer-infra/environments/dev && terraform output -raw droplet_ipv4)

# Grafana
ssh root@$DROPLET_IP "k3s kubectl -n monitoring get secret grafana-admin-credentials -o jsonpath='{.data.admin-password}' | base64 -d; echo"

# Argo CD (argocd-initial-admin-secret is created by Argo CD's own installer, not by cloud-init)
ssh root@$DROPLET_IP "k3s kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d; echo"
```

## Architecture

### Application

```text
REQUEST PATH
────────────
Browser (HTMX + Tailwind)
   │ HTTP
   ▼
chi Router  (otel + logger + metrics middleware)
   │
   ▼
handlers: article · source · auth
   │
   ▼
usecases: article · source · tag · auth · sync
   │
   ▼
adapter/postgres repositories
   │
   ▼
Postgres


SYNC PATH — two triggers, same usecase
───────────────────────────────────────
"Sync" button (admin, POST /sources/sync)     cmd/syncjob (Kubernetes CronJob)
                    │                                       │
                    └───────────────────┬───────────────────┘
                                         ▼
                           usecase/sync (worker pool fetch)
                                         │
                                         ▼
                    adapter/feed — RSS/Atom fetch + parse
                                         │
                                         ▼
         source feeds: Go blogs · newsletters · podcasts · job boards
                                         │
                                         ▼
           new articles → Postgres (deduped by source_id + external_id)


OBSERVABILITY — cloud only (locally OTEL_ENDPOINT is unset, so this is a no-op)
────────────────────────────────────────────────────────────────────────────
Prometheus ──scrapes GET /metrics──────► chi Router                    (pull)
chi Router ──pushes spans (OTLP/gRPC)──► Tempo                         (push)
chi Router ──zerolog → stdout──► Alloy (tails pod logs) ──► Loki       (push, via Alloy)
                                                               │
                                                               ▼
                                Grafana — dashboards over Prometheus + Tempo + Loki
```

### Deployment pipeline

```text
news-explorer-infra (Terraform)
      │  terraform apply — Droplet + Firewall + DNS, one-time/on infra change
      ▼
DigitalOcean Droplet
      │  cloud-init: installs k3s, bootstraps Argo CD, applies the one
      │  "root" Application by hand (github.com/khainouski/news-platform-deploy,
      │  path argocd/applications, recurse) — everything after this is GitOps
      ▼
Argo CD "root" Application (app-of-apps)
      ▲
      │ watches & auto-syncs (prune + selfHeal) every Application committed under
      │ argocd/applications/ in news-platform-deploy
      │
git push to main (news-explorer) ──► GitHub Actions CI ──► GHCR (image, tag sha-<commit>)
                                                                 │
                    bump image.tag in argocd/applications/news-explorer.yaml, git push
                                                                 ▼
                                              news-platform-deploy (Helm chart +
                                              ~12 Argo CD Application manifests)
                                                                 │
                                                                 ▼
                                                          k3s cluster
      ┌───────────────────────────────────────────────────────────────────────────┐
      │  apps:         news-explorer — Deployment, sync CronJob (hourly),         │
      │                 migrate Job, Ingress                                      │
      │  database:     Postgres (Bitnami Helm chart)                              │
      │  monitoring:   Prometheus · Tempo · Loki · Grafana · Alloy (log shipper)  │
      │  cert-manager: cert-manager + ClusterIssuer (Let's Encrypt TLS)           │
      │  kube-system:  Traefik (k3s's built-in ingress, HTTP→HTTPS redirect)      │
      └───────────────────────────────────────────────────────────────────────────┘
      │                                          │
      │ metrics · traces · logs                  │ HTTPS (goskills.xyz)
      ▼                                          ▼
Grafana dashboards                          Browser / news-explorer-client
                                             (synthetic load generator)
```

`news-explorer-infra` provisions the box the cluster runs on and bootstraps Argo CD via
cloud-init; `news-platform-deploy` is what actually puts everything on it — the app plus the
Prometheus/Tempo/Loki/Grafana stack, Postgres, and cert-manager — as one app-of-apps root
Application that keeps every child Application in sync with the repo. `news-explorer-client`
isn't part of the release path — it's a standalone traffic generator pointed at the live
deployment so the observability stack has something non-trivial to show.
