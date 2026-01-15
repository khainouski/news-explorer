# syntax=docker/dockerfile:1

FROM golang:1.25.6-alpine AS build

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

FROM alpine:3.22 AS run

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=build /out/app ./app

USER app

EXPOSE 8080

ENTRYPOINT ["./app"]
