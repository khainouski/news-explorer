# syntax=docker/dockerfile:1

FROM golang:1.26.5-alpine AS build

WORKDIR /src

# Modules layer cached separately from source for faster rebuilds.
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# Standalone Tailwind CLI, pinned to v3.4.17 - must run before go build so go:embed picks up
# the real app.css.
ARG TARGETARCH
RUN apk add --no-cache curl \
    && case "$TARGETARCH" in \
         amd64) TW_ARCH=x64; TW_SHA=7d24f7fa191d2193b78cd5f5a42a6093e14409521908529f42d80b11fde1f1d4 ;; \
         arm64) TW_ARCH=arm64; TW_SHA=69b1378b8133192d7d2feb12a116fa12d035594f58db3eff215879e4ad8cf39b ;; \
         *) echo "unsupported TARGETARCH: $TARGETARCH" >&2; exit 1 ;; \
       esac \
    && curl -sLo /usr/local/bin/tailwindcss "https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-$TW_ARCH" \
    && echo "$TW_SHA  /usr/local/bin/tailwindcss" | sha256sum -c - \
    && chmod +x /usr/local/bin/tailwindcss \
    && tailwindcss -i web/tailwind/input.css -o web/static/css/app.css --config web/tailwind/config.js --minify

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/app \
    ./cmd/app

# Runs as a Kubernetes CronJob (news-platform-deploy: sync-cronjob.yaml), same image as ./app,
# just a different command.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/syncjob \
    ./cmd/syncjob

# Same version Makefile's migrate-install pins - adds ~6MB, measured. Lets the deploy pipeline's
# migrate Job (news-platform-deploy: migrate-job.yaml) reuse this image instead of a second one.
RUN CGO_ENABLED=0 GOOS=linux go install \
    -ldflags="-s -w" \
    -tags 'postgres' \
    github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1

FROM alpine:3.22 AS run

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=build /out/app ./app
COPY --from=build /out/syncjob ./syncjob
COPY --from=build /go/bin/migrate /usr/local/bin/migrate
COPY migration/ ./migration/

USER app

EXPOSE 8080

ENTRYPOINT ["./app"]
