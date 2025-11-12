---
sidebar_position: 2
---

# Migrations

Learn how to manage database migrations in go-infra applications using GORM and migration tools.

## Overview

Database migrations help you:

- **Version control** your database schema
- **Safely deploy** schema changes
- **Roll back** changes if needed
- **Collaborate** with team members

## Using GORM AutoMigrate

The simplest way to manage schema in development:

```go
package main

import (
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type User struct {
    ID    string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name  string `gorm:"not null"`
    Email string `gorm:"uniqueIndex;not null"`
}

type Product struct {
    ID          string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name        string  `gorm:"not null"`
    Description string
    Price       float64 `gorm:"not null"`
}

func main() {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // Auto migrate models
    err = db.AutoMigrate(
        &User{},
        &Product{},
    )
    if err != nil {
        log.Fatal("Failed to migrate:", err)
    }

    log.Println("Migration completed successfully")
}
```

### When to Use AutoMigrate

✅ **Good for:**

- Development environment
- Prototyping
- Simple schema changes
- Small projects

❌ **Not recommended for:**

- Production deployments
- Complex migrations
- Data transformations
- Team collaboration

## SQL Migration Files

For production, use SQL migration files:

### Directory Structure

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_products_table.up.sql
├── 000002_create_products_table.down.sql
├── 000003_add_user_role.up.sql
└── 000003_add_user_role.down.sql
```

### Example Migrations

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

```sql
-- 000001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

## Using golang-migrate

### Installation

```bash
# Install CLI tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Add to go.mod
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/postgres
go get -u github.com/golang-migrate/migrate/v4/source/file
```

### Create Migration

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq create_users_table

# This creates:
# migrations/000001_create_users_table.up.sql
# migrations/000001_create_users_table.down.sql
```

### Run Migrations

```bash
# Migrate up
migrate -path migrations \
  -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" \
  up

# Migrate down
migrate -path migrations \
  -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" \
  down

# Migrate to specific version
migrate -path migrations \
  -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" \
  goto 2

# Force version (if migration is dirty)
migrate -path migrations \
  -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" \
  force 1
```

### Programmatic Migrations

```go
// internal/migration/migrate.go
package migration

import (
    "fmt"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) error {
    m, err := migrate.New(
        "file://migrations",
        dbURL,
    )
    if err != nil {
        return err
    }
    defer m.Close()

    // Run all migrations
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    version, dirty, err := m.Version()
    if err != nil {
        return err
    }

    fmt.Printf("Migration completed. Version: %d, Dirty: %v\n", version, dirty)
    return nil
}
```

Usage in main:

```go
package main

import (
    "log"
    "myapp/internal/migration"
)

func main() {
    dbURL := "postgres://user:pass@localhost:5432/mydb?sslmode=disable"

    if err := migration.RunMigrations(dbURL); err != nil {
        log.Fatal("Migration failed:", err)
    }

    // Continue with application startup
}
```

## Using go-infra Migration Component

go-infra provides a built-in migration component:

```go
package main

import (
    "github.com/phatnt199/go-infra/pkg/component/migration"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    // Create migrator
    migrator := migration.NewMigrator(db, migration.Config{
        MigrationsPath: "./migrations",
        TableName:      "schema_migrations",
    })

    // Run migrations
    if err := migrator.Up(); err != nil {
        log.Fatal(err)
    }
}
```

## Migration Best Practices

### 1. Always Create Down Migrations

```sql
-- up.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- down.sql
ALTER TABLE users DROP COLUMN phone;
```

### 2. Make Migrations Idempotent

```sql
-- Safe to run multiple times
CREATE TABLE IF NOT EXISTS users (...);

ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

### 3. Use Transactions

```sql
BEGIN;

CREATE TABLE products (...);
CREATE INDEX idx_products_name ON products(name);

COMMIT;
```

### 4. Never Modify Existing Migrations

```bash
# ❌ Wrong - Don't modify existing migrations
# Edit: 000001_create_users_table.up.sql

# ✅ Correct - Create new migration
migrate create -ext sql -dir migrations -seq add_user_phone_column
```

### 5. Test Migrations

```go
func TestMigrations(t *testing.T) {
    // Setup test database
    db := setupTestDB()

    // Run migrations
    err := migration.RunMigrations(db)
    assert.NoError(t, err)

    // Verify schema
    assert.True(t, db.Migrator().HasTable("users"))
    assert.True(t, db.Migrator().HasColumn(&User{}, "email"))
}
```

## Data Migrations

For migrating data, create separate migrations:

```sql
-- 000004_migrate_user_roles.up.sql
BEGIN;

-- Add new column
ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'user';

-- Migrate existing data
UPDATE users SET role = 'admin' WHERE email LIKE '%@company.com';

-- Make column not null
ALTER TABLE users ALTER COLUMN role SET NOT NULL;

COMMIT;
```

## Makefile Commands

Create helpful Makefile targets:

```makefile
# Makefile
DB_URL := postgres://user:pass@localhost:5432/mydb?sslmode=disable

.PHONY: migrate-create
migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

.PHONY: migrate-reset
migrate-reset:
	migrate -path migrations -database "$(DB_URL)" drop -f
	migrate -path migrations -database "$(DB_URL)" up

.PHONY: migrate-version
migrate-version:
	migrate -path migrations -database "$(DB_URL)" version
```

Usage:

```bash
# Create migration
make migrate-create

# Run migrations
make migrate-up

# Rollback last migration
make migrate-down

# Reset database
make migrate-reset
```

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/migrate.yml
name: Database Migration

on:
  push:
    branches: [main]
    paths:
      - "migrations/**"

jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install migrate
        run: |
          curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
          sudo mv migrate /usr/local/bin/

      - name: Run migrations
        env:
          DATABASE_URL: ${{ secrets.DATABASE_URL }}
        run: |
          migrate -path migrations -database "$DATABASE_URL" up
```

## Docker Integration

```dockerfile
# Dockerfile.migrate
FROM golang:1.21-alpine

RUN apk add --no-cache curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/

WORKDIR /app
COPY migrations ./migrations

ENTRYPOINT ["migrate", "-path", "migrations", "-database"]
CMD ["up"]
```

```yaml
# docker-compose.yml
version: "3.8"
services:
  migrate:
    build:
      context: .
      dockerfile: Dockerfile.migrate
    command: ["postgres://user:pass@postgres:5432/mydb?sslmode=disable", "up"]
    depends_on:
      - postgres
```

## Troubleshooting

### Dirty Migration State

```bash
# Check current version
migrate -path migrations -database "$DB_URL" version

# Force to specific version
migrate -path migrations -database "$DB_URL" force 5
```

### Roll Back Failed Migration

```bash
# Force version before failed migration
migrate -path migrations -database "$DB_URL" force 3

# Run down migration manually
psql -d mydb -f migrations/000004_failed_migration.down.sql

# Continue migrations
migrate -path migrations -database "$DB_URL" up
```

## Next Steps

- Learn about [Advanced Database Features](./advanced)
- Explore [Repository Pattern](./repository-pattern)
- See [Testing Database Code](../testing/database-testing)
