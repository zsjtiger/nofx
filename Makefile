# Crux Makefile for testing and development

.PHONY: help test test-backend test-frontend test-coverage clean

# Default target
help:
	@echo "Crux Testing & Development Commands"
	@echo ""
	@echo "Testing:"
	@echo "  make test                 - Run all tests (backend + frontend)"
	@echo "  make test-backend         - Run backend tests only"
	@echo "  make test-frontend        - Run frontend tests only"
	@echo "  make test-coverage        - Generate backend coverage report"
	@echo ""
	@echo "Build:"
	@echo "  make build                - Build backend binary"
	@echo "  make build-frontend       - Build frontend"
	@echo ""
	@echo "Clean:"
	@echo "  make clean                - Clean build artifacts and test cache"

# =============================================================================
# Testing
# =============================================================================

# Run all tests
test:
	@echo "🧪 Running backend tests..."
	cd server && go test -v ./...
	@echo ""
	@echo "🧪 Running frontend tests..."
	cd web && npm run test
	@echo "✅ All tests completed"

# Backend tests only
test-backend:
	@echo "🧪 Running backend tests..."
	cd server && go test -v ./...

# Frontend tests only
test-frontend:
	@echo "🧪 Running frontend tests..."
	cd web && npm run test

# Coverage report
test-coverage:
	@echo "📊 Generating coverage..."
	cd server && go test -coverprofile=coverage.out ./...
	cd server && go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Backend coverage: coverage.html"

# =============================================================================
# Build
# =============================================================================

# Build backend binary
build:
	@echo "🔨 Building backend..."
	cd server && go build -o crux
	@echo "✅ Backend built: ./server/crux"

# Build frontend
build-frontend:
	@echo "🔨 Building frontend..."
	cd web && npm run build
	@echo "✅ Frontend built: ./web/dist"

# =============================================================================
# Development
# =============================================================================

# Run backend in development mode
run:
	@echo "🚀 Starting backend..."
	cd server && go run main.go

# Run frontend in development mode
run-frontend:
	@echo "🚀 Starting frontend dev server..."
	cd web && npm run dev

# Format Go code
fmt:
	@echo "🎨 Formatting Go code..."
	cd server && go fmt ./...
	@echo "✅ Code formatted"

# Lint Go code (requires golangci-lint)
lint:
	@echo "🔍 Linting Go code..."
	cd server && golangci-lint run
	@echo "✅ Linting completed"

# =============================================================================
# Clean
# =============================================================================

clean:
	@echo "🧹 Cleaning..."
	rm -f server/nofx
	rm -f server/coverage.out server/coverage.html
	rm -rf web/dist
	cd server && go clean -testcache
	@echo "✅ Cleaned"

# =============================================================================
# Docker
# =============================================================================

# Build Docker images
docker-build:
	@echo "🐳 Building Docker images..."
	docker compose build
	@echo "✅ Docker images built"

# Run Docker containers
docker-up:
	@echo "🐳 Starting Docker containers..."
	docker compose up -d
	@echo "✅ Docker containers started"

# Stop Docker containers
docker-down:
	@echo "🐳 Stopping Docker containers..."
	docker compose down
	@echo "✅ Docker containers stopped"

# View Docker logs
docker-logs:
	docker compose logs -f

# =============================================================================
# Dependencies
# =============================================================================

# Download Go dependencies
deps:
	@echo "📦 Downloading Go dependencies..."
	cd server && go mod download
	@echo "✅ Dependencies downloaded"

# Update Go dependencies
deps-update:
	@echo "📦 Updating Go dependencies..."
	cd server && go get -u ./...
	cd server && go mod tidy
	@echo "✅ Dependencies updated"

# Install frontend dependencies
deps-frontend:
	@echo "📦 Installing frontend dependencies..."
	cd web && npm install
	@echo "✅ Frontend dependencies installed"
