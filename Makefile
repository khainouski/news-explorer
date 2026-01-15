APP_NAME := news-explorer
IMAGE    := $(APP_NAME):local

.PHONY: run build test lint tidy docker-build docker-run

run:
	go run ./cmd/app

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
