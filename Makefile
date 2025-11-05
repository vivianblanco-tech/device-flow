.PHONY: help run build test migrate-up migrate-down migrate-create dev-setup clean install

# Variables
BINARY_NAME=laptop-tracking
MIGRATE=migrate
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
TEST_DB_URL=postgres://postgres:password@localhost:5432/laptop_tracking_test?sslmode=disable

# Load environment variables
include .env
export

# Export test database URL for all test targets
export TEST_DATABASE_URL=$(TEST_DB_URL)

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

test: test-all ## Run all tests (alias for test-all)

test-all: ## Run all tests with test database (sequential, reliable)
	@echo "Running all tests with test database..."
	@echo "Database: $(TEST_DB_URL)"
	go test ./... -p=1 -v -race -cover

test-parallel: ## Run all tests in parallel (faster but may have conflicts)
	@echo "Running tests in parallel..."
	@echo "Note: May have database conflicts. Use 'make test-all' for reliability."
	go test ./... -v -race -cover

test-unit: ## Run unit tests only (no database required)
	@echo "Running unit tests (short mode)..."
	go test -v -race -short ./...

test-integration: ## Run integration tests only (requires database)
	@echo "Running integration tests..."
	go test ./... -p=1 -v -run Integration

test-package: ## Run tests for specific package (usage: make test-package PKG=internal/auth)
	@echo "Running tests for package: $(PKG)..."
	go test ./$(PKG) -v -race

test-verbose: ## Run all tests with verbose output
	@echo "Running all tests with verbose output..."
	go test ./... -p=1 -v -race -cover

test-coverage: ## Run tests with coverage report (HTML)
	@echo "Running tests with coverage analysis..."
	go test ./... -p=1 -race -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-coverage-summary: ## Run tests and show coverage summary
	@echo "Running tests with coverage summary..."
	go test ./... -p=1 -race -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out | tail -n 1

test-watch: ## Run tests in watch mode (requires gotestsum)
	@echo "Running tests in watch mode..."
	gotestsum --watch -- -p=1 ./...

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

test-db-setup: ## Set up test database (Docker)
	@echo "Setting up test database in Docker..."
	docker exec laptop-tracking-db psql -U postgres -c "DROP DATABASE IF EXISTS laptop_tracking_test;" || true
	docker exec laptop-tracking-db psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"
	$(MIGRATE) -path migrations -database "$(TEST_DB_URL)" up
	@echo "✓ Test database ready!"
	@echo "  Database: laptop_tracking_test"
	@echo "  URL: $(TEST_DB_URL)"

test-db-reset: ## Reset test database (clean slate)
	@echo "Resetting test database..."
	docker exec laptop-tracking-db psql -U postgres -c "DROP DATABASE IF EXISTS laptop_tracking_test;" || true
	docker exec laptop-tracking-db psql -U postgres -c "CREATE DATABASE laptop_tracking_test;"
	$(MIGRATE) -path migrations -database "$(TEST_DB_URL)" up
	@echo "✓ Test database reset complete!"

test-db-clean: ## Clean test database (remove all data, keep schema)
	@echo "Cleaning test database..."
	docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_test -c "TRUNCATE TABLE audit_logs, notification_logs, magic_links, sessions, delivery_forms, reception_reports, pickup_forms, shipment_laptops, shipments, laptops, software_engineers, client_companies, users CASCADE;"
	@echo "✓ Test database cleaned!"

test-db-verify: ## Verify test database setup
	@echo "Verifying test database..."
	@docker exec laptop-tracking-db psql -U postgres -l | grep laptop_tracking_test || (echo "✗ Test database not found!" && exit 1)
	@docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_test -c "\dt" | grep -q "users" || (echo "✗ Tables not found!" && exit 1)
	@echo "✓ Test database verified!"
	@docker exec laptop-tracking-db psql -U postgres -d laptop_tracking_test -c "SELECT COUNT(*) as table_count FROM information_schema.tables WHERE table_schema = 'public';"

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

test-quick: ## Quick test run (unit tests only, no race detection)
	@echo "Running quick tests..."
	go test -short ./...

test-ci: ## Run tests in CI mode (with coverage and sequential execution)
	@echo "Running tests in CI mode..."
	go test ./... -p=1 -race -coverprofile=coverage.out -covermode=atomic -v
	go tool cover -func=coverage.out

test-help: ## Show detailed test command help
	@echo "Test Commands:"
	@echo ""
	@echo "Basic Testing:"
	@echo "  make test-all              - Run all tests (sequential, reliable) [RECOMMENDED]"
	@echo "  make test-parallel         - Run all tests (parallel, faster but may conflict)"
	@echo "  make test-unit             - Run only unit tests (no database)"
	@echo "  make test-quick            - Quick test run (unit tests, no race detection)"
	@echo ""
	@echo "Specific Testing:"
	@echo "  make test-package PKG=path - Run tests for specific package"
	@echo "  make test-integration      - Run integration tests only"
	@echo "  make test-verbose          - Run all tests with verbose output"
	@echo ""
	@echo "Coverage:"
	@echo "  make test-coverage         - Run tests with HTML coverage report"
	@echo "  make test-coverage-summary - Show coverage summary"
	@echo ""
	@echo "Database Management:"
	@echo "  make test-db-setup         - Set up test database from scratch"
	@echo "  make test-db-reset         - Reset test database (drop and recreate)"
	@echo "  make test-db-clean         - Clean test data (keep schema)"
	@echo "  make test-db-verify        - Verify test database is set up correctly"
	@echo ""
	@echo "Environment:"
	@echo "  TEST_DB_URL: $(TEST_DB_URL)"
	@echo ""
	@echo "Examples:"
	@echo "  make test-all                              # Run all tests"
	@echo "  make test-package PKG=internal/auth        # Test specific package"
	@echo "  make test-db-reset test-all                # Reset DB and run tests"

