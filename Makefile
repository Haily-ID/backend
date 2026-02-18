.PHONY: help run-api run-worker dev-api dev-worker build test test-coverage lint docker-up docker-down clean

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run-api: ## Run API server
	go run cmd/api/main.go

run-worker: ## Run worker
	go run cmd/worker/main.go

dev-api: ## Run API server with hot reload
	air -c .air.toml

dev-worker: ## Run worker with hot reload
	air -c .air.worker.toml

build: ## Build binaries
	@echo "Building API..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api cmd/api/main.go
	@echo "Building Worker..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/worker cmd/worker/main.go
	@echo "Build complete!"

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

docker-up: ## Start Docker services
	docker-compose up -d

docker-down: ## Stop Docker services
	docker-compose down

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf tmp/
