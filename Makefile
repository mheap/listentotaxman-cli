.PHONY: test test-verbose test-coverage test-race test-update-golden test-ci test-pkg help

# Default target
.DEFAULT_GOAL := help

# Help target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run all tests (quiet mode)
	@go test ./...

test-verbose: ## Run tests with verbose output
	@go test -v ./...

test-race: ## Run tests with race detector
	@go test -race ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out -covermode=atomic ./...
	@echo "\nFiltering out testutil from coverage..."
	@grep -v "internal/testutil" coverage.out > coverage-filtered.out || true
	@mv coverage-filtered.out coverage.out
	@echo "\nCoverage Summary (excluding testutil):"
	@go tool cover -func=coverage.out | grep total
	@echo "\nGenerating HTML report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in your browser to view detailed report"

test-update-golden: ## Update golden files
	@echo "Updating golden files..."
	@go test ./internal/display/... -update-golden
	@echo "Golden files updated. Please review changes before committing."

test-ci: ## Simulate CI environment (race + coverage check)
	@echo "Running CI tests..."
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "\nFiltering out testutil from coverage..."
	@grep -v "internal/testutil" coverage.out > coverage-filtered.out || true
	@mv coverage-filtered.out coverage.out
	@echo "\nChecking coverage threshold..."
	@./scripts/check-coverage.sh 90
	@echo "\n✓ All CI checks passed!"

test-pkg: ## Run tests for a specific package (use: make test-pkg PKG=cmd)
	@go test -v ./$(PKG)

build: ## Build the binary
	@go build -o listentotaxman .

clean: ## Clean build artifacts and test cache
	@rm -f listentotaxman coverage.out coverage.html
	@go clean -testcache

lint: ## Run golangci-lint
	@golangci-lint run

fmt: ## Format code with gofmt
	@gofmt -s -w .

vet: ## Run go vet
	@go vet ./...

tidy: ## Tidy go modules
	@go mod tidy

all: tidy fmt vet lint test ## Run all checks (tidy, fmt, vet, lint, test)
