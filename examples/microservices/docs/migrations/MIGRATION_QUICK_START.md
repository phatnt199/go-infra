# User Service - Migration Quick Start

A quick reference guide for common migration tasks.

## 🚀 Getting Started

### 1. Start the Service (Runs Migrations Automatically)

```bash
cd examples/microservices/userservice
go run ./cmd/app/main.go
```

Output indicates successful migration:

```
2025/11/05 14:35:56 OK   00001_create_users_table.sql (4.6ms)
2025/11/05 14:35:56 goose: successfully migrated database to version: 1
2025/11/05 14:35:56.946+0700 | INFO | User Service is listening on...
```

### 2. Connect to Database

```bash
psql -h localhost -p 54100 -U admin -d go-sandbox
```

Check migrations applied:

```sql
SELECT * FROM goose_db_version;
```

---

## 📝 Create a New Migration

### Step 1: Create Migration File

Create `db/migrations/goose-migrate/00NNN_description.sql`:

```bash
touch db/migrations/goose-migrate/00002_add_profile_column.sql
```

### Step 2: Write the Migration

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN IF NOT EXISTS profile TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS profile;
-- +goose StatementEnd
```

### Step 3: Restart Service

```bash
go run ./cmd/app/main.go
```

---

## 🔨 Common Tasks

### Add a Column

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS phone VARCHAR(20);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS phone;
-- +goose StatementEnd
```

### Drop a Column

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS phone;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20);
-- +goose StatementEnd
```

### Create an Index

```sql
-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_users_email
    ON users(email) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;
-- +goose StatementEnd
```

### Create a New Table

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_profiles;
-- +goose StatementEnd
```

### Add a Constraint

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD CONSTRAINT unique_email_active
    UNIQUE (email) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT IF EXISTS unique_email_active;
-- +goose StatementEnd
```

---

## 🌱 Add Seed Data

Edit `internal/shared/configurations/users/users_configurator_seed.go`:

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
        // ADD YOUR DATA HERE
        {
            ID:    uuid.NewV4(),
            Name:  "New User",
            Email: "newuser@example.com",
        },
    }

    err := db.CreateInBatches(users, len(users)).Error
    if err != nil {
        return errors.Wrap(err, "Error seeding users data")
    }

    return nil
}
```

### Reseed Data

```bash
# Connect to database
psql -h localhost -p 54100 -U admin -d go-sandbox

# Delete existing data
DELETE FROM users;

# Restart service
# Migration runs automatically and reseeds
```

---

## ❌ Troubleshooting

### Migration File Not Found

```
panic: no migration files found
```

✅ **Solution:**

```bash
# Check directory exists
ls -la db/migrations/goose-migrate/

# Verify config path
cat config.development.json | grep migrationsDir
```

### Duplicate Version Error

```
panic: goose: duplicate version 1 detected
```

✅ **Solution:** Use single `.sql` file with `-- +goose Up/Down` sections:

- ✅ `00001_create_table.sql`
- ❌ `00001_create_table.up.sql` + `.down.sql`

### Column Already Exists

```
ERROR: column "status" already exists
```

✅ **Solution:** Always use `IF NOT EXISTS`:

```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20);
```

### UUID Extension Not Found

```
ERROR: type "uuid" does not exist
```

✅ **Solution:** Create migration to enable UUID extension:

File: `db/migrations/goose-migrate/00000_enable_uuid.sql`

```sql
-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +goose Down
DROP EXTENSION IF EXISTS "uuid-ossp";
```

---

## 📊 Check Migration Status

```bash
# Connect to database
psql -h localhost -p 54100 -U admin -d go-sandbox

# View all applied migrations
SELECT * FROM goose_db_version ORDER BY version_id;

# List all tables
\dt

# View table structure
\d users
```

---

## 📖 Full Documentation

See `MIGRATION_GUIDE.md` for comprehensive documentation on:

- Architecture and how migrations work
- Detailed step-by-step examples
- Best practices
- Advanced troubleshooting

---

## 🎯 Workflow Summary

1. **Create migration file** → `00NNN_description.sql`
2. **Write SQL** → Add `-- +goose Up` and `-- +goose Down` sections
3. **Start service** → `go run ./cmd/app/main.go`
4. **Verify** → Check `goose_db_version` table
5. **Update models** → Update GORM datamodels if needed
6. **Test** → Restart service and verify data

Done! 🎉
