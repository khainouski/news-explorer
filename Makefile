APP_NAME := news-explorer
IMAGE    := $(APP_NAME):local

DB_MIGRATE_URL := postgres://news_explorer:local@localhost:5432/news_explorer?sslmode=disable
MIGRATE_PATH   := ./migration/postgres

.PHONY: run build test lint tidy docker-build docker-run up down migrate-install migrate-create migrate-up migrate-down tailwind-install tailwind-build

run:
	LOG_PRETTY=true go run ./cmd/app

build:
	CGO_ENABLED=0 go build -o bin/$(APP_NAME) ./cmd/app

test:
	go test -race -cover ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

docker-build:
	docker build -t $(IMAGE) .

docker-run:
	docker run --rm -p 8080:8080 $(IMAGE)

# Postgres for local dev - the app itself still runs separately via `make run`, not inside this
# compose.
up:
	docker compose up --build --force-recreate

down:
	docker compose down

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1

migrate-create:
	@read -p "Name:" name; \
	migrate create -ext sql -dir "$(MIGRATE_PATH)" $$name

migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all

tailwind-build:
	bin/tailwindcss -i web/tailwind/input.css -o web/static/css/app.css --config web/tailwind/config.js --minify

tailwind-install:
	mkdir -p bin
	curl -sLo bin/tailwindcss "https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-$$(uname -s | tr 'A-Z' 'a-z' | sed 's/darwin/macos/')-$$(uname -m | sed 's/x86_64/x64/;s/aarch64/arm64/')"
	chmod +x bin/tailwindcss
