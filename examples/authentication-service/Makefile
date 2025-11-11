.PHONY: help run migrate swagger test clean build dev

# Default target
help:
	@echo "Authentication Service - Development Commands"
	@echo ""
	@echo "Available commands:"
	@echo "  make run        - Run the application"
	@echo "  make migrate    - Run database migrations"
	@echo "  make swagger    - Generate swagger documentation"
	@echo "  make test       - Run API tests"
	@echo "  make dev        - Run migrations, swagger, and start server"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make build      - Build the application"
	@echo "  make docker-up  - Start PostgreSQL with Docker"
	@echo "  make docker-down- Stop PostgreSQL Docker container"

# Run the application
run:
	@echo "Starting authentication service..."
	go run cmd/app/main.go

# Run database migrations
migrate:
	@echo "Running database migrations..."
	go run cmd/migration/main.go

# Generate swagger documentation
swagger:
	@echo "Generating swagger documentation..."
	@if ! command -v swag &> /dev/null; then \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init --parseDependency --parseInternal -g cmd/app/main.go -o docs
	@echo "Swagger documentation generated successfully!"

# Run API tests
test:
	@echo "Running API tests..."
	@chmod +x test-apis.sh
	./test-apis.sh

# Development workflow: migrate + swagger + run
dev: migrate swagger
	@echo "Starting development server..."
	@echo "API available at: http://localhost:8081/api/v1"
	@echo "Swagger UI at: http://localhost:8081/swagger/index.html"
	@echo "Health check at: http://localhost:8081/health"
	go run cmd/app/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f authentication-service
	@echo "Clean completed!"

# Build the application
build:
	@echo "Building authentication service..."
	@mkdir -p bin
	go build -o bin/authentication-service cmd/app/main.go
	@echo "Build completed! Binary: bin/authentication-service"

# Start PostgreSQL with Docker
docker-up:
	@echo "Starting PostgreSQL with Docker..."
	docker run --name auth-postgres -d \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=auth_service_db \
		-p 5432:5432 \
		postgres:13-alpine
	@echo "PostgreSQL started on port 5432"
	@echo "Wait 5 seconds for PostgreSQL to be ready..."
	@sleep 5

# Stop PostgreSQL Docker container
docker-down:
	@echo "Stopping PostgreSQL..."
	docker stop auth-postgres || true
	docker rm auth-postgres || true
	@echo "PostgreSQL stopped"

# Full setup for new developers
setup: docker-up migrate swagger
	@echo ""
	@echo "Setup completed! You can now run:"
	@echo "  make run"
	@echo ""
	@echo "Or for auto-reload development:"
	@echo "  make dev"

# Check dependencies
check-deps:
	@echo "Checking dependencies..."
	@command -v go >/dev/null 2>&1 || { echo "Go is not installed!"; exit 1; }
	@command -v psql >/dev/null 2>&1 || { echo "Warning: PostgreSQL client not installed"; }
	@command -v docker >/dev/null 2>&1 || { echo "Warning: Docker not installed"; }
	@echo "Dependencies check completed!"

# Database operations
db-create:
	@echo "Creating database..."
	createdb auth_service_db

db-drop:
	@echo "Dropping database..."
	dropdb auth_service_db

db-reset: db-drop db-create migrate
	@echo "Database reset completed!"

# View logs
logs:
	@echo "Viewing application logs..."
	@tail -f logs/*.log || echo "No log files found"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatting completed!"

# Run linter
lint:
	@echo "Running linter..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed!"

# Show environment variables
env:
	@echo "Environment Configuration:"
	@echo "=========================="
	@grep -v '^#' .env 2>/dev/null || echo ".env file not found"
