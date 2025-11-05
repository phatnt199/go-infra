# User Service - Database Migration Guide

This guide explains how to manage database migrations, schema changes, and data seeding in the User Service microservice.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Directory Structure](#directory-structure)
4. [Configuration](#configuration)
5. [How Migrations Work](#how-migrations-work)
6. [Migration File Format](#migration-file-format)
7. [Common Tasks](#common-tasks)
8. [Best Practices](#best-practices)
9. [Troubleshooting](#troubleshooting)

---

## Overview

The User Service uses **Goose** as its database migration tool. Goose is a simple and fast SQL migration tool written in Go that allows you to:

- Track database schema changes with version control
- Rollback changes if needed
- Apply migrations automatically on service startup
- Support both forward and backward migrations

### Key Features

- ✅ **Automatic on Startup**: Migrations run automatically when the service starts
- ✅ **Version Tracking**: Goose maintains a `goose_db_version` table to track applied migrations
- ✅ **Reversible**: Each migration has an "Up" (apply) and "Down" (rollback) section
- ✅ **Type-Safe**: Uses SQL directly for database operations
- ✅ **Seeding Support**: Data seeding is handled separately via GORM

---

## Architecture

### Migration Flow

```
Application Startup
        ↓
[UsersServiceConfigurator.ConfigureUsersService()]
        ↓
[UsersServiceConfigurator.migrateUsers()]
        ↓
[PostgresMigrationRunner.Up()] (Goose)
        ↓
Read all .sql files from db/migrations/goose-migrate/
        ↓
Execute pending migrations in version order
        ↓
Update goose_db_version table
        ↓
Application Ready
```

### Components

| Component                   | File                                                                   | Purpose                                                |
| --------------------------- | ---------------------------------------------------------------------- | ------------------------------------------------------ |
| **Migration Configuration** | `config.development.json`                                              | Specifies migrations directory and database connection |
| **Migration Handler**       | `internal/shared/configurations/users/users_configurator_migration.go` | Triggers migration execution                           |
| **Migration Files**         | `db/migrations/goose-migrate/*.sql`                                    | Actual schema change SQL files                         |
| **Data Seeding**            | `internal/shared/configurations/users/users_configurator_seed.go`      | Populates initial data                                 |
| **Models**                  | `internal/users/data/datamodels/user_data_model.go`                    | GORM models reflecting database schema                 |

---

## Directory Structure

```
userservice/
├── db/
│   └── migrations/
│       └── goose-migrate/
│           ├── 00001_create_users_table.sql
│           ├── 00002_add_profile_column.sql
│           └── ...
├── config.development.json
├── config.production.json
├── internal/
│   ├── users/
│   │   ├── data/
│   │   │   └── datamodels/
│   │   │       └── user_data_model.go
│   │   └── configurations/
│   ├── shared/
│   │   └── configurations/
│   │       └── users/
│   │           ├── users_configurator_migration.go
│   │           └── users_configurator_seed.go
│   └── ...
└── cmd/
    └── app/
        └── main.go
```

---

## Configuration

### Migration Settings

Edit `config.development.json` (or your environment config):

```json
{
	"migrationOptions": {
		"host": "localhost",
		"port": 54100,
		"user": "admin",
		"password": "123456",
		"dbName": "go-sandbox",
		"sslMode": false,
		"migrationsDir": "db/migrations/goose-migrate",
		"skipMigration": false
	}
}
```

**Key Configuration Fields:**

| Field           | Description                  | Default                       |
| --------------- | ---------------------------- | ----------------------------- |
| `host`          | Database server hostname     | `localhost`                   |
| `port`          | PostgreSQL port              | `54100`                       |
| `user`          | Database user                | `admin`                       |
| `password`      | Database password            | `123456`                      |
| `dbName`        | Database name                | `go-sandbox`                  |
| `sslMode`       | Enable SSL for DB connection | `false`                       |
| `migrationsDir` | Path to migration files      | `db/migrations/goose-migrate` |
| `skipMigration` | Skip migrations on startup   | `false`                       |

### Disabling Migrations

To skip migrations during startup:

```json
{
	"migrationOptions": {
		"skipMigration": true
	}
}
```

---

## How Migrations Work

### Startup Process

When the User Service starts:

1. **Service Initialization**: The application calls `UsersServiceConfigurator.ConfigureUsersService()`
2. **Migration Trigger**: This calls `migrateUsers()` which invokes `migrateGoose()`
3. **Database Connection**: Goose connects to PostgreSQL using credentials from `config.development.json`
4. **Version Check**: Goose reads the `goose_db_version` table to determine which migrations have been applied
5. **Pending Migrations**: Identifies all `.sql` files with version numbers higher than the current version
6. **Sequential Execution**: Executes pending migrations in version order (ascending)
7. **Version Update**: After each successful migration, the version is recorded in `goose_db_version`
8. **Seeding**: After migrations complete, `seedUsers()` runs to insert initial data (idempotent)

### Version Tracking

Goose maintains a table called `goose_db_version`:

```sql
SELECT * FROM goose_db_version;
```

Output example:

```
 version_id | is_applied |            tstamp
------------+------------+-----------------------------
          1 | t          | 2025-11-05 14:35:56.456713
```

---

## Migration File Format

### File Naming Convention

```
NNNNN_description.sql
```

- **NNNNN**: 5-digit version number (e.g., `00001`, `00002`, `00003`)
- **description**: Human-readable name in snake_case
- Must be unique and sequential

**Examples:**

```
00001_create_users_table.sql
00002_add_profile_column.sql
00003_create_orders_table.sql
00004_add_indexes.sql
```

### Migration Structure

Every migration file must have `-- +goose Up` and `-- +goose Down` sections:

```sql
-- +goose Up
-- +goose StatementBegin
-- Your forward migration SQL here
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Your rollback SQL here
-- +goose StatementEnd
```

### Complete Example: Initial Users Table

**File: `00001_create_users_table.sql`**

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
```

---

## Common Tasks

### Task 1: Create a New Table

**Scenario:** Add a new `user_profiles` table to store additional user information.

#### Step 1: Create Migration File

Create `db/migrations/goose-migrate/00002_create_user_profiles_table.sql`:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    bio TEXT,
    avatar_url VARCHAR(500),
    phone VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_profiles_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_profiles_user_id;
DROP TABLE IF EXISTS user_profiles;
-- +goose StatementEnd
```

#### Step 2: Create GORM Model

Create `internal/users/data/datamodels/user_profile_data_model.go`:

```go
package datamodels

import (
    "time"
    uuid "github.com/satori/go.uuid"
    "gorm.io/gorm"
)

type UserProfileDataModel struct {
    ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
    UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
    Bio       string    `gorm:"column:bio;type:text"`
    AvatarURL string    `gorm:"column:avatar_url;type:varchar(500)"`
    Phone     string    `gorm:"column:phone;type:varchar(20)"`
    CreatedAt time.Time `gorm:"column:created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u *UserProfileDataModel) TableName() string {
    return "user_profiles"
}
```

#### Step 3: Start the Service

The migration runs automatically on startup:

```bash
cd examples/microservices/userservice
go run ./cmd/app/main.go
```

Expected output:

```
2025/11/05 14:35:56 OK   00002_create_user_profiles_table.sql (2.3ms)
2025/11/05 14:35:56 goose: successfully migrated database to version: 2
```

---

### Task 2: Add a New Column

**Scenario:** Add a `status` column to the `users` table.

**File: `db/migrations/goose-migrate/00003_add_status_to_users.sql`**

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'active';
ALTER TABLE users ADD COLUMN IF NOT EXISTS status_updated_at TIMESTAMP WITH TIME ZONE;

-- Create index for status-based queries
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_status;
ALTER TABLE users DROP COLUMN IF EXISTS status_updated_at;
ALTER TABLE users DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
```

**Important Notes:**

- Use `IF NOT EXISTS` / `IF EXISTS` for idempotency
- Order alterations logically (add columns, then indexes)
- Always reverse the order in the Down section

---

### Task 3: Create an Index

**Scenario:** Add an index for faster queries on created_at.

**File: `db/migrations/goose-migrate/00004_add_indexes.sql`**

```sql
-- +goose Up
-- +goose StatementBegin
-- Create composite index for common queries
CREATE INDEX IF NOT EXISTS idx_users_created_at
    ON users(created_at DESC) WHERE deleted_at IS NULL;

-- Create index for email search (case-insensitive)
CREATE INDEX IF NOT EXISTS idx_users_email_lower
    ON users(LOWER(email));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email_lower;
DROP INDEX IF EXISTS idx_users_created_at;
-- +goose StatementEnd
```

---

### Task 4: Modify a Column

**Scenario:** Change the `name` column from `VARCHAR(255)` to `VARCHAR(500)`.

**File: `db/migrations/goose-migrate/00005_increase_name_length.sql`**

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN name TYPE VARCHAR(500);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN name TYPE VARCHAR(255);
-- +goose StatementEnd
```

---

### Task 5: Add a Constraint

**Scenario:** Add a UNIQUE constraint on email addresses.

**File: `db/migrations/goose-migrate/00006_add_unique_email.sql`**

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD CONSTRAINT unique_email_active
    UNIQUE (email) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT IF EXISTS unique_email_active;
-- +goose StatementEnd
```

---

### Task 6: Seed Initial Data

Data seeding is handled in `internal/shared/configurations/users/users_configurator_seed.go`.

#### Current Seeding Logic

```go
func seedDataManually(db *gorm.DB) error {
    var count int64

    // Check if data already exists (idempotent)
    db.Model(&datamodels.UserDataModel{}).Count(&count)
    if count > 0 {
        return nil // Skip if data exists
    }

    // Create seed data
    users := []datamodels.UserDataModel{
        {
            ID:    uuid.NewV4(),
            Name:  "John Doe",
            Email: "john.doe@example.com",
        },
        {
            ID:    uuid.NewV4(),
            Name:  "Jane Smith",
            Email: "jane.smith@example.com",
        },
    }

    // Insert in batches
    err := db.CreateInBatches(users, len(users)).Error
    if err != nil {
        return errors.Wrap(err, "Error seeding users data")
    }

    return nil
}
```

#### How to Add More Seed Data

Edit `users_configurator_seed.go`:

```go
func seedDataManually(db *gorm.DB) error {
    var count int64

    db.Model(&datamodels.UserDataModel{}).Count(&count)
    if count > 0 {
        return nil
    }

    users := []datamodels.UserDataModel{
        {
            ID:    uuid.NewV4(),
            Name:  "John Doe",
            Email: "john.doe@example.com",
        },
        {
            ID:    uuid.NewV4(),
            Name:  "Jane Smith",
            Email: "jane.smith@example.com",
        },
        // ADD NEW SEED DATA HERE
        {
            ID:    uuid.NewV4(),
            Name:  "Bob Johnson",
            Email: "bob.johnson@example.com",
        },
    }

    err := db.CreateInBatches(users, len(users)).Error
    if err != nil {
        return errors.Wrap(err, "Error seeding users data")
    }

    return nil
}
```

**Important:** Seeding is **idempotent** - it checks if data exists before inserting. To reseed:

```sql
-- Delete all records
DELETE FROM users;

-- Then restart the service
go run ./cmd/app/main.go
```

---

### Task 7: Rollback a Migration

**Scenario:** You applied migration 00003 but need to revert it.

#### Approach 1: Manual Database Rollback

Connect to the database and manually execute the "Down" SQL from the migration file:

```bash
psql -h localhost -p 54100 -U admin -d go-sandbox
```

```sql
-- Get current version
SELECT * FROM goose_db_version;

-- Manually run the Down section of your migration
-- (copy the SQL from -- +goose Down section)
DROP INDEX IF EXISTS idx_users_status;
ALTER TABLE users DROP COLUMN IF EXISTS status_updated_at;
ALTER TABLE users DROP COLUMN IF EXISTS status;

-- Update version (set to previous version number)
UPDATE goose_db_version SET version_id = 2;
```

#### Approach 2: Delete Migration and Restart

1. Delete or rename the migration file you want to undo
2. Manually update the `goose_db_version` table to the previous version
3. Restart the service

**Example:**

```bash
# Rename the migration file
mv db/migrations/goose-migrate/00003_add_status_to_users.sql \
   db/migrations/goose-migrate/00003_add_status_to_users.sql.bak

# Connect to database and downgrade version
psql -h localhost -p 54100 -U admin -d go-sandbox -c \
    "UPDATE goose_db_version SET version_id = 2;"

# Restart service
go run ./cmd/app/main.go
```

---

## Best Practices

### 1. ✅ Always Use IF EXISTS / IF NOT EXISTS

```sql
-- Good
ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20);
DROP INDEX IF EXISTS idx_users_email;
CREATE TABLE IF NOT EXISTS users (...)

-- Bad (will fail if column already exists)
ALTER TABLE users ADD COLUMN status VARCHAR(20);
```

### 2. ✅ Make Migrations Reversible

Every migration needs a complete Down section:

```sql
-- +goose Up
CREATE TABLE new_table (...);

-- +goose Down
DROP TABLE new_table;
```

### 3. ✅ Create Indexes on Foreign Keys

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
```

### 4. ✅ Use Meaningful Names

```
❌ 00001_change.sql
❌ 00002_fix.sql
✅ 00001_create_users_table.sql
✅ 00002_add_email_index.sql
```

### 5. ✅ Add Comments Explaining Changes

```sql
-- +goose Up
-- Adding soft delete support: users can be "deleted" without removing data
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
-- +goose StatementEnd
```

### 6. ✅ Test Migrations Locally

```bash
# Run the service locally to test migrations
go run ./cmd/app/main.go

# Check that tables were created
psql -h localhost -p 54100 -U admin -d go-sandbox

# Inside psql:
\dt  -- List all tables
\di  -- List all indexes
SELECT * FROM goose_db_version;
```

### 7. ✅ Seed Data with Conditions

Check before seeding to avoid duplicates:

```go
func seedDataManually(db *gorm.DB) error {
    var count int64

    // This check makes seeding idempotent
    db.Model(&datamodels.UserDataModel{}).Count(&count)
    if count > 0 {
        return nil  // Skip if data exists
    }

    // ... create and insert seed data
}
```

---

## Troubleshooting

### Issue 1: "Migration file not found"

**Symptom:**

```
panic: no migration files found
```

**Solution:**

1. Check that `db/migrations/goose-migrate/` directory exists
2. Ensure migration files exist in the directory:
   ```bash
   ls -la db/migrations/goose-migrate/
   ```
3. Verify the path in `config.development.json`:
   ```json
   "migrationsDir": "db/migrations/goose-migrate"
   ```

---

### Issue 2: "Duplicate version detected"

**Symptom:**

```
panic: goose: duplicate version 1 detected:
    00001_create_users_table.sql.up
    00001_create_users_table.sql.down
```

**Solution:**

- Rename migration files to use single `.sql` file with `-- +goose Up/Down` sections
- ✅ Correct: `00001_create_users_table.sql`
- ❌ Wrong: `00001_create_users_table.up.sql` + `00001_create_users_table.down.sql`

---

### Issue 3: "Column already exists"

**Symptom:**

```
ERROR: column "status" of relation "users" already exists (SQLSTATE 42701)
```

**Solution:**
Always use `IF NOT EXISTS`:

```sql
-- Wrong
ALTER TABLE users ADD COLUMN status VARCHAR(20);

-- Correct
ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20);
```

---

### Issue 4: "Cannot drop table, foreign key constraint exists"

**Symptom:**

```
ERROR: cannot drop table users because other objects depend on it
```

**Solution:**
Drop dependent objects first:

```sql
-- +goose Down
DROP TABLE user_profiles;  -- Drop dependent table first
DROP TABLE users;          -- Then drop parent table
```

---

### Issue 5: "Database connection refused"

**Symptom:**

```
ERROR: could not connect to server: Connection refused
```

**Solution:**

1. Verify PostgreSQL is running:
   ```bash
   docker ps | grep postgres
   ```
2. Check database credentials in `config.development.json`
3. Verify database exists:
   ```bash
   psql -h localhost -p 54100 -U admin -d postgres -c \
       "SELECT datname FROM pg_database WHERE datname='go-sandbox';"
   ```

---

### Issue 6: "UUID type not supported"

**Symptom:**

```
ERROR: type "uuid" does not exist
```

**Solution:**
Enable the UUID extension in PostgreSQL:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

Or create a migration file `00000_enable_uuid_extension.sql`:

```sql
-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +goose Down
DROP EXTENSION IF EXISTS "uuid-ossp";
```

---

### Issue 7: "Syntax error in migration"

**Symptom:**

```
ERROR: syntax error at or near...
```

**Solution:**

1. Review the SQL syntax in your migration file
2. Test the SQL manually in psql:
   ```bash
   psql -h localhost -p 54100 -U admin -d go-sandbox < \
       db/migrations/goose-migrate/00001_create_users_table.sql
   ```
3. Check for:
   - Missing semicolons
   - Unclosed quotes or parentheses
   - Invalid SQL keywords
   - Database name quoting (use `"` for names with hyphens)

---

## Quick Reference

### Database Connection

```bash
# Connect to the service database
psql -h localhost -p 54100 -U admin -d go-sandbox

# List tables
\dt

# List indexes
\di

# Show table schema
\d users

# Show version tracking
SELECT * FROM goose_db_version;
```

### Migration Commands

```bash
# Run service (applies pending migrations)
cd examples/microservices/userservice
go run ./cmd/app/main.go

# Create a new migration file (template)
touch db/migrations/goose-migrate/00NNN_description.sql

# Check migration status
psql -h localhost -p 54100 -U admin -d go-sandbox -c \
    "SELECT * FROM goose_db_version ORDER BY version_id;"
```

### Common SQL Patterns

```sql
-- Create table with soft delete
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index for soft delete queries
CREATE INDEX idx_users_active ON users(id) WHERE deleted_at IS NULL;

-- Create index on foreign key
CREATE INDEX idx_orders_user_id ON orders(user_id);

-- Create composite index
CREATE INDEX idx_users_email_status ON users(email, status)
    WHERE deleted_at IS NULL;

-- Add unique constraint
ALTER TABLE users ADD CONSTRAINT unique_email_active
    UNIQUE (email) WHERE deleted_at IS NULL;
```

---

## Summary

The User Service migration system:

1. **Automatically runs on startup** - No manual intervention needed
2. **Uses Goose for tracking** - Version control in `goose_db_version` table
3. **Supports rollbacks** - Each migration has Up and Down sections
4. **Handles seeding** - Idempotent seed data insertion via GORM
5. **Type-safe** - Direct SQL for precise control over schema

Follow the patterns and best practices above to maintain a clean, reversible schema evolution!
