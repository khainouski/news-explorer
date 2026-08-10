# syntax=docker/dockerfile:1

FROM golang:1.26.5-alpine AS build

WORKDIR /src

# Modules layer cached separately from source for faster rebuilds.
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

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
