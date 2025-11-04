# Users API Example

## Purpose

This example demonstrates how to use the `go-infra` components to build a small REST API that exposes a `GET /api/v1/users` endpoint backed by PostgreSQL (GORM), with migrations and Swagger documentation. It's intended as a minimal, runnable example you can use to learn wiring, configuration, and migrations in the `go-infra` ecosystem.

What this example shows

- Application bootstrap wiring (module, DI, config)
- Database integration (Postgres via GORM)
- Migrations (using `cmd/migration` and sql files in `migrations/`)
- HTTP server (Fiber) with route(s) and Swagger docs

## Prerequisites

- Go 1.25+ installed
- PostgreSQL 12+ (local or Docker)
- `swag` CLI if you want to regenerate Swagger locally (optional; the repository includes generated docs in `docs/`)

## Quick setup (recommended)

1. Open a terminal and change to this example directory:

```bash
cd /path/to/go-infra/examples/users-api
```

2. Primary configuration is provided by `internal/config/config.development.json`. You can use the `.env` file as a secondary/override source if you prefer environment variables.

```bash
# Primary config (already present): internal/config/config.development.json
# Optional: copy and edit .env to provide overrides via environment variables
cp .env.example .env
# Edit .env only if you want to override values from the JSON config
```

3. Install Go module dependencies:

```bash
make deps
```

4. Start or ensure PostgreSQL is running (Docker example):

```bash
docker run --name postgres-users \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=usersdb \
  -p 5432:5432 \
  -d postgres:15-alpine
```

5. Run database migrations:

```bash
make migrate
```

6. Run the server (this target generates Swagger before running):

```bash
make run
```

By default the server listens on the address configured in `internal/config/config.development.json`. You may optionally override configuration values via `.env` (copied from `.env.example`). The common default HTTP port is `:8080`.

## Makefile targets (recommended)

This example now includes a `Makefile` with these targets:

- `deps` — download go modules
- `swagger` — regenerate Swagger docs (requires `swag` installed)
- `migrate` — run migration binary (`cmd/migration`) to apply migrations
- `run` — generate swagger and run the API (`cmd/api/main.go`)
- `docker-db` — convenience target to run a postgres docker container for local testing

## Direct commands (no Makefile)

If you prefer not to use `make`, these commands perform the same actions:

- Install deps: `go mod download`
- Generate swagger: `swag init -g cmd/api/main.go -o docs`
- Run migrations: `go run cmd/migration/main.go`
- Run server: `go run cmd/api/main.go`

## API endpoints

- `GET /api/v1/users` — returns all users
- `GET /api/v1/health` — health check
- `GET /swagger/` — Swagger UI (when server runs and `docs/` exists)

## Project layout (key files)

```
examples/users-api/
├── cmd/
│   ├── api/main.go        # App entrypoint (HTTP server)
│   └── migration/main.go  # Migration runner
├── internal/
│   ├── config/            # app config
│   ├── domain/user.go     # user model
│   ├── handler/user_handler.go
│   └── repository/user_repository.go
├── migrations/            # SQL migration files
├── docs/                  # generated Swagger files (regenerate with `make swagger`)
├── .env.example
├── go.mod
└── Makefile               # convenience targets: deps, migrate, run, swagger
```

## Configuration

Primary configuration for this example is the JSON file at `internal/config/config.development.json`. That file contains the application defaults (database connection, server options, migration options, etc.).

If you prefer to provide overrides via environment variables, copy `.env.example` to `.env` and edit values there. Environment variables act as secondary/override values and will take precedence at runtime when present.

To skip automatic migrations, set one of the following depending on how you run the app:

Using environment variable (optional override):

```bash
MIGRATION_OPTIONS_SKIP_MIGRATION=true
```

Or edit `internal/config/config.development.json` to change migration options in the primary config.

## Notes & troubleshooting

- If you get a port conflict, change `FIBER_HTTP_OPTIONS_PORT` in `.env`.
- If Swagger UI shows no routes, regenerate docs with `make swagger` and restart.
- If migrations fail, inspect the SQL files under `migrations/` and ensure DB credentials in `.env` are correct.

---
