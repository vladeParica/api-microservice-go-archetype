BINARY      := api
BIN_DIR     := bin
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X main.version=$(VERSION)

DATABASE_URL    ?= postgres://app:app@localhost:5432/app?sslmode=disable
MIGRATE_IMAGE   := migrate/migrate:v4.17.1
MIGRATE         := docker run --rm --network host -v $(CURDIR)/db/migrations:/migrations:ro $(MIGRATE_IMAGE) -path=/migrations -database "$(DATABASE_URL)"

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## ---- Development ----

.PHONY: run
run: ## Run the API locally
	go run ./cmd/api

.PHONY: build
build: ## Build the binary into ./bin
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) ./cmd/api

.PHONY: test
test: ## Run unit tests with the race detector
	go test -race -count=1 ./...

.PHONY: cover
cover: ## Run tests and open the coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: fmt
fmt: ## Format the code
	gofmt -s -w .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

## ---- Database ----

.PHONY: sqlc
sqlc: ## Generate Go code from SQL queries
	sqlc generate

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## Roll back the last migration
	$(MIGRATE) down 1

.PHONY: migrate-new
migrate-new: ## Create a new migration: make migrate-new name=add_orders
	docker run --rm -v $(CURDIR)/db/migrations:/migrations $(MIGRATE_IMAGE) create -ext sql -dir /migrations -seq $(name)

## ---- Docker ----

.PHONY: up
up: ## Start the full stack with docker compose
	docker compose up --build -d

.PHONY: down
down: ## Stop the docker compose stack
	docker compose down

.PHONY: docker-build
docker-build: ## Build the container image
	docker build --build-arg VERSION=$(VERSION) -t $(BINARY):$(VERSION) .
