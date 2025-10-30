.PHONY: help run build test migrate-up migrate-down migrate-create dev-setup clean install

# Variables
BINARY_NAME=laptop-tracking
MIGRATE=migrate
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Load environment variables
include .env
export

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install Go dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/$(BINARY_NAME) cmd/web/main.go

run: ## Run the application
	go run cmd/web/main.go

dev: ## Run the application with hot reload (requires air)
	air

test: ## Run tests
	go test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

migrate-up: ## Run all database migrations
	$(MIGRATE) -path migrations -database "$(DB_URL)" up

migrate-down: ## Rollback last database migration
	$(MIGRATE) -path migrations -database "$(DB_URL)" down 1

migrate-create: ## Create a new migration file (usage: make migrate-create name=create_users_table)
	$(MIGRATE) create -ext sql -dir migrations -seq $(name)

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	$(MIGRATE) -path migrations -database "$(DB_URL)" force $(version)

db-reset: ## Reset database (drop and recreate)
	dropdb $(DB_NAME) || true
	createdb $(DB_NAME)
	$(MAKE) migrate-up

dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo ".env file created from .env.example"; fi
	@echo "Please update .env with your configuration"
	@echo "Installing dependencies..."
	@$(MAKE) install
	@echo "Setup complete!"

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME):latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

lint: ## Run linters
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	go vet ./...

check: fmt vet lint test ## Run all checks (format, vet, lint, test)

