.PHONY: run build test clean help dev install tidy fmt vet lint

# Variables
APP_NAME=usermes-backend
MAIN_PATH=./cmd/main.go
BUILD_DIR=./build
BINARY_NAME=$(APP_NAME)

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Colors for terminal output
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

## help: Show this help message
help:
	@echo '$(GREEN)Usage:$(NC)'
	@echo '  make <target>'
	@echo ''
	@echo '$(GREEN)Targets:$(NC)'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run: Run the application
run:
	@echo "$(GREEN)Running $(APP_NAME)...$(NC)"
	@$(GORUN) $(MAIN_PATH)

## dev: Run the application with auto-reload (requires air)
dev:
	@echo "$(YELLOW)Starting development server...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(YELLOW)Air not installed. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		$(GORUN) $(MAIN_PATH); \
	fi

## build: Build the application
build:
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## build-linux: Build for Linux
build-linux:
	@echo "$(GREEN)Building $(APP_NAME) for Linux...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux$(NC)"

## build-windows: Build for Windows
build-windows:
	@echo "$(GREEN)Building $(APP_NAME) for Windows...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME).exe$(NC)"

## build-mac: Build for macOS
build-mac:
	@echo "$(GREEN)Building $(APP_NAME) for macOS...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-mac $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)-mac$(NC)"

## build-all: Build for all platforms
build-all: build-linux build-windows build-mac
	@echo "$(GREEN)All builds complete!$(NC)"

## test: Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	@$(GOTEST) -v -race ./...
	@echo "$(GREEN)Tests complete!$(NC)"

## test-coverage: Run tests with coverage report and check threshold
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@./scripts/coverage.sh

## coverage: Alias for test-coverage
coverage: test-coverage

## coverage-html: Generate HTML coverage report
coverage-html:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@$(GOTEST) -coverprofile=coverage.out ./...
	@echo "$(GREEN)Generating HTML report...$(NC)"
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

## coverage-func: Show coverage by function
coverage-func:
	@$(GOTEST) -coverprofile=coverage.out ./... > /dev/null 2>&1
	@$(GOCMD) tool cover -func=coverage.out

## bench: Run benchmarks
bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@$(GOTEST) -bench=. -benchmem ./...

## install: Install dependencies
install:
	@echo "$(GREEN)Installing dependencies...$(NC)"
	@$(GOMOD) download
	@echo "$(GREEN)Dependencies installed!$(NC)"

## tidy: Tidy up go.mod and go.sum
tidy:
	@echo "$(GREEN)Tidying up dependencies...$(NC)"
	@$(GOMOD) tidy
	@echo "$(GREEN)Dependencies tidied!$(NC)"

## fmt: Format Go code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	@$(GOFMT) ./...
	@echo "$(GREEN)Code formatted!$(NC)"

## vet: Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	@$(GOVET) ./...
	@echo "$(GREEN)Vet complete!$(NC)"

## lint: Run golangci-lint (requires golangci-lint)
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Install from https://golangci-lint.run/$(NC)"; \
	fi

## check: Run fmt, vet, and lint
check: fmt vet lint
	@echo "$(GREEN)All checks passed!$(NC)"

## clean: Clean build artifacts and cache
clean:
	@echo "$(GREEN)Cleaning...$(NC)"
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean complete!$(NC)"

## docker-build: Build Docker image
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t $(APP_NAME):latest .
	@echo "$(GREEN)Docker image built: $(APP_NAME):latest$(NC)"

## docker-run: Run Docker container
docker-run:
	@echo "$(GREEN)Running Docker container...$(NC)"
	@docker run -p 3000:3000 $(APP_NAME):latest

## up: Start all services with docker-compose
up:
	@echo "$(GREEN)Starting services...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Services started!$(NC)"

## down: Stop all services
down:
	@echo "$(GREEN)Stopping services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Services stopped!$(NC)"

## logs: Show logs from docker-compose
logs:
	@docker-compose logs -f

## migrate-up: Run database migrations up
migrate-up:
	@echo "$(GREEN)Running migrations up...$(NC)"
	@# Add your migration command here
	@echo "$(YELLOW)Not implemented yet$(NC)"

## migrate-down: Run database migrations down
migrate-down:
	@echo "$(GREEN)Running migrations down...$(NC)"
	@# Add your migration command here
	@echo "$(YELLOW)Not implemented yet$(NC)"

## swagger: Generate Swagger documentation
swagger:
	@echo "$(GREEN)Generating Swagger docs...$(NC)"
	@if command -v swag > /dev/null; then \
		swag init -g cmd/main.go; \
	else \
		echo "$(YELLOW)swag not installed. Install with: go install github.com/swaggo/swag/cmd/swag@latest$(NC)"; \
	fi

## air-init: Initialize Air for hot reload
air-init:
	@echo "$(GREEN)Initializing Air...$(NC)"
	@if ! command -v air > /dev/null; then \
		echo "$(YELLOW)Installing Air...$(NC)"; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@air init
	@echo "$(GREEN)Air initialized! Run 'make dev' to start$(NC)"

.DEFAULT_GOAL := help
