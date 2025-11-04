APP_CMD := ./cmd/api/main.go
MIGRATE_CMD := ./cmd/migration/main.go

.PHONY: help deps swagger migrate migrate-up migrate-down migrate-status run docker-db stop-db

help:
	@echo "Makefile targets for examples/users-api"
	@echo "  deps           - download go module dependencies"
	@echo "  swagger        - generate swagger docs (requires swag CLI)"
	@echo "  migrate        - run all pending database migrations (UP)"
	@echo "  migrate-up     - run all pending database migrations (UP)"
	@echo "  migrate-down   - rollback the last database migration (DOWN)"
	@echo "  migrate-status - show current migration version and status"
	@echo "  run            - regenerate swagger and run the API server"
	@echo "  docker-db      - run a local postgres docker instance for testing"
	@echo "  stop-db        - stop and remove the local postgres docker instance"

deps:
	@echo "Downloading go modules..."
	go mod download

swagger:
	@echo "Generating swagger docs (swag must be installed)..."
	@swag init -g cmd/api/main.go -o docs || echo "swag not found or failed; ensure swag is installed if you want to regenerate docs"

migrate: migrate-up

migrate-up:
	@echo "Running migrations UP..."
	go run $(MIGRATE_CMD) -direction=up

migrate-down:
	@echo "Running migrations DOWN (last version)..."
	go run $(MIGRATE_CMD) -direction=down

migrate-status:
	@echo "Getting migration status..."
	go run $(MIGRATE_CMD) -direction=status

run: swagger
	@echo "Starting API server (go run $(APP_CMD))..."
	go run $(APP_CMD)

docker-db:
	@echo "Starting local postgres container 'postgres-users'..."
	docker run --name postgres-users \
	  -e POSTGRES_USER=admin \
	  -e POSTGRES_PASSWORD=123456 \
	  -e POSTGRES_DB=go-sandbox \
	  -p 54100:5432 \
	  -d postgres:15-alpine

stop-db:
	@echo "Stopping and removing postgres-users container (if exists)..."
	docker rm -f postgres-users || true
