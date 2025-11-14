---
sidebar_position: 1
---

# Installation

## Prerequisites

Before you begin, ensure you have:

- **Go 1.25+** installed
- **PostgreSQL** (if using database features)
- **Git** for version control

Check your Go version:

```bash
go version
```

## Install go-infra

### Option 1: Add to Existing Project

Add go-infra to your `go.mod`:

```bash
go get github.com/phatnt199/go-infra
```

### Option 2: Local Development

For local development or customization:

```bash
# Clone the repository
git clone https://github.com/phatnt199/go-infra.git

# In your project's go.mod, add:
replace github.com/phatnt199/go-infra => /path/to/go-infra
```

## Verify Installation

Create a simple test file:

```go title="main.go"
package main

import (
    "github.com/phatnt199/go-infra/pkg/logger"
)

func main() {
    logger.Info("go-infra installed successfully!")
}
```

Run it:

```bash
go run main.go
```

You should see a log message indicating successful installation.

## Database Setup (Optional)

If you plan to use database features:

### PostgreSQL with Docker

```bash
docker run --name postgres-dev \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=myapp \
  -p 5432:5432 \
  -d postgres:15-alpine
```

### PostgreSQL Installation

**macOS:**

```bash
brew install postgresql@15
brew services start postgresql@15
```

**Ubuntu/Debian:**

```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
```

**Windows:**

Download and install from [PostgreSQL Official Site](https://www.postgresql.org/download/windows/)

## Additional Tools

### Swagger CLI (for API documentation)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Migration Tools

**Goose (recommended):**

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

**Go-Migrate:**

```bash
# macOS
brew install golang-migrate

# Ubuntu/Debian
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/bin/migrate

# Windows
choco install migrate
```

### Air (for hot reload during development)

```bash
go install github.com/cosmtrek/air@latest
```

## Project Structure

Create a new project with this recommended structure:

```
myapp/
├── cmd/
│   ├── api/
│   │   └── main.go          # HTTP server entry point
│   └── migration/
│       └── main.go          # Migration runner
├── internal/
│   ├── domain/              # Domain models
│   ├── repository/          # Data access layer
│   ├── handler/             # HTTP handlers
│   └── service/             # Business logic
├── migrations/              # Database migrations
├── config/                  # Configuration files
├── .env.example             # Environment template
├── go.mod
└── go.sum
```

Create this structure:

```bash
mkdir -p myapp/{cmd/{api,migration},internal/{domain,repository,handler,service},migrations,config}
cd myapp
go mod init myapp
```

## Environment Configuration

Create `.env.example`:

```bash title=".env.example"
# Application
APP_ENV=development
SERVICE_NAME=myapp

# HTTP Server
HTTP_PORT=8080
HTTP_HOST=localhost

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=myapp
POSTGRES_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_ISSUER=myapp
JWT_AUDIENCE=myapp-users
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d
```

Copy to `.env`:

```bash
cp .env.example .env
```

:::tip
Add `.env` to your `.gitignore` to avoid committing secrets:

```bash
echo ".env" >> .gitignore
```

:::

## Verify Complete Setup

Create a complete test application:

```go title="cmd/api/main.go"
package main

import (
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
    "github.com/phatnt199/go-infra/pkg/logger"
)

func main() {
    logger.Info("Starting application...")

    app := fxapp.NewApplicationBuilder().
        ProvideModule(fiber_adapter.Module).
        ProvideModule(gorm.Module).
        Build()

    app.Run()
}
```

Run it:

```bash
go run cmd/api/main.go
```

You should see:

- Server starting on configured port
- Database connection established
- Application ready to accept requests

Visit `http://localhost:8080/health` to verify it's running.

## Troubleshooting

### "Cannot find package"

Make sure you've run:

```bash
go mod download
go mod tidy
```

### Database Connection Fails

Verify PostgreSQL is running:

```bash
# macOS/Linux
psql -U postgres -h localhost

# Check .env has correct credentials
```

### Port Already in Use

Change the port in `.env`:

```bash
HTTP_PORT=8081
```

## Next Steps

Now that you have go-infra installed, let's build something!

- **[Quick Start](./quick-start)** - Build your first API
- **[Architecture](../core-concepts/architecture)** - Understand the framework
- **[Examples](../examples/users-api)** - See real implementations
