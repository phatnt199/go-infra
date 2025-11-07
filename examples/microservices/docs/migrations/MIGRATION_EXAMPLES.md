# User Service - Migration Examples

Practical, copy-paste ready migration examples.

---

## 📋 Table of Contents

1. [Initial Setup](#initial-setup)
2. [Schema Evolution](#schema-evolution)
3. [Performance Optimization](#performance-optimization)
4. [Data Transformations](#data-transformations)
5. [Advanced Patterns](#advanced-patterns)

---

## Initial Setup

### Example 1: Create Users Table (00001)

**File:** `db/migrations/goose-migrate/00001_create_users_table.sql`

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

## Schema Evolution

### Example 2: Add User Profile Fields (00002)

**File:** `db/migrations/goose-migrate/00002_add_profile_fields.sql`

**Scenario:** Users want to add profile information like bio, avatar, phone number.

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS bio TEXT,
    ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500),
    ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
    ADD COLUMN IF NOT EXISTS location VARCHAR(255),
    ADD COLUMN IF NOT EXISTS website_url VARCHAR(500);

-- Create indexes for profile searches
CREATE INDEX IF NOT EXISTS idx_users_location ON users(location)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_location;
ALTER TABLE users
    DROP COLUMN IF EXISTS website_url,
    DROP COLUMN IF EXISTS location,
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS avatar_url,
    DROP COLUMN IF EXISTS bio;
-- +goose StatementEnd
```

**Update GORM Model:**

```go
type UserDataModel struct {
    ID         uuid.UUID  `gorm:"column:id;type:uuid;primaryKey"`
    Name       string     `gorm:"column:name"`
    Email      string     `gorm:"column:email"`
    Password   string     `gorm:"column:password"`
    Bio        string     `gorm:"column:bio;type:text"`
    AvatarURL  string     `gorm:"column:avatar_url"`
    Phone      string     `gorm:"column:phone"`
    Location   string     `gorm:"column:location"`
    WebsiteURL string     `gorm:"column:website_url"`
    CreatedAt  time.Time  `gorm:"column:created_at"`
    UpdatedAt  time.Time  `gorm:"column:updated_at"`
    DeletedAt  *time.Time `gorm:"column:deleted_at"`
}
```

---

### Example 3: Add User Status (00003)

**File:** `db/migrations/goose-migrate/00003_add_user_status.sql`

**Scenario:** Track user account status (active, inactive, suspended, pending_verification).

```sql
-- +goose Up
-- +goose StatementBegin
-- Create ENUM type for status
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'pending_verification');

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS status user_status DEFAULT 'pending_verification',
    ADD COLUMN IF NOT EXISTS status_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Create index for status-based queries
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)
    WHERE deleted_at IS NULL;

-- Create composite index for common queries
CREATE INDEX IF NOT EXISTS idx_users_status_created_at ON users(status, created_at DESC)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_status_created_at;
DROP INDEX IF EXISTS idx_users_status;
ALTER TABLE users
    DROP COLUMN IF EXISTS status_updated_at,
    DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS user_status;
-- +goose StatementEnd
```

**Update GORM Model:**

```go
import "database/sql"

type UserStatus string

const (
    UserStatusActive              UserStatus = "active"
    UserStatusInactive            UserStatus = "inactive"
    UserStatusSuspended           UserStatus = "suspended"
    UserStatusPendingVerification UserStatus = "pending_verification"
)

type UserDataModel struct {
    // ... existing fields ...
    Status           sql.NullString `gorm:"column:status;type:user_status"`
    StatusUpdatedAt  *time.Time     `gorm:"column:status_updated_at"`
}
```

---

## Performance Optimization

### Example 4: Add Search Indexes (00004)

**File:** `db/migrations/goose-migrate/00004_add_search_indexes.sql`

**Scenario:** Improve query performance for common search operations.

```sql
-- +goose Up
-- +goose StatementBegin
-- Case-insensitive email search
CREATE INDEX IF NOT EXISTS idx_users_email_lower ON users(LOWER(email));

-- Full-text search on name and bio
CREATE INDEX IF NOT EXISTS idx_users_name_tsvector
    ON users USING GIN(to_tsvector('english', COALESCE(name, '') || ' ' || COALESCE(bio, '')));

-- Range query optimization for created_at
CREATE INDEX IF NOT EXISTS idx_users_created_at_desc ON users(created_at DESC)
    WHERE deleted_at IS NULL;

-- Soft delete awareness index
CREATE INDEX IF NOT EXISTS idx_users_active_email ON users(email)
    WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_active_email;
DROP INDEX IF EXISTS idx_users_created_at_desc;
DROP INDEX IF EXISTS idx_users_name_tsvector;
DROP INDEX IF EXISTS idx_users_email_lower;
-- +goose StatementEnd
```

---

### Example 5: Add Audit Trail (00005)

**File:** `db/migrations/goose-migrate/00005_create_audit_log.sql`

**Scenario:** Track all changes to user records for compliance.

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,  -- 'CREATE', 'UPDATE', 'DELETE'
    changed_fields JSONB,  -- JSON object of changes
    old_values JSONB,      -- Previous values
    new_values JSONB,      -- New values
    changed_by UUID,       -- Who made the change
    ip_address INET,       -- IP address of change
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_audit_logs_user_id ON user_audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_audit_logs_action ON user_audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_user_audit_logs_created_at ON user_audit_logs(created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_audit_logs;
-- +goose StatementEnd
```

---

## Data Transformations

### Example 6: Rename Column (00006)

**File:** `db/migrations/goose-migrate/00006_rename_column.sql`

**Scenario:** Rename `avatar_url` to `profile_image_url` for consistency.

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN avatar_url TO profile_image_url;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN profile_image_url TO avatar_url;
-- +goose StatementEnd
```

---

### Example 7: Change Column Type (00007)

**File:** `db/migrations/goose-migrate/00007_change_column_type.sql`

**Scenario:** Change `phone` from VARCHAR to a structured phone number format.

```sql
-- +goose Up
-- +goose StatementBegin
-- Create new phone column with correct type
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_new VARCHAR(20);

-- Copy data from old column (with transformation if needed)
UPDATE users SET phone_new = phone WHERE phone IS NOT NULL;

-- Drop old column and rename
ALTER TABLE users DROP COLUMN IF EXISTS phone;
ALTER TABLE users RENAME COLUMN phone_new TO phone;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert: create old column and copy back
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_old VARCHAR(20);
UPDATE users SET phone_old = phone WHERE phone IS NOT NULL;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
ALTER TABLE users RENAME COLUMN phone_old TO phone;
-- +goose StatementEnd
```

---

### Example 8: Add NOT NULL Constraint (00008)

**File:** `db/migrations/goose-migrate/00008_add_not_null_constraint.sql`

**Scenario:** Make `name` field required (NOT NULL) after cleaning up data.

```sql
-- +goose Up
-- +goose StatementBegin
-- First, set default values for NULL records
UPDATE users SET name = 'User' WHERE name IS NULL;

-- Then add the constraint
ALTER TABLE users ALTER COLUMN name SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the constraint (allowing NULLs again)
ALTER TABLE users ALTER COLUMN name DROP NOT NULL;
-- +goose StatementEnd
```

---

## Advanced Patterns

### Example 9: Create User Preferences Table with Inheritance (00009)

**File:** `db/migrations/goose-migrate/00009_create_user_preferences.sql`

**Scenario:** Store user preferences and settings separately.

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,

    -- Communication preferences
    email_notifications BOOLEAN DEFAULT true,
    sms_notifications BOOLEAN DEFAULT false,
    push_notifications BOOLEAN DEFAULT true,

    -- Privacy settings
    profile_visibility VARCHAR(20) DEFAULT 'public',  -- public, friends, private
    show_email BOOLEAN DEFAULT false,
    show_phone BOOLEAN DEFAULT false,

    -- Theme and language
    theme VARCHAR(20) DEFAULT 'light',  -- light, dark, auto
    language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',

    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_preferences;
-- +goose StatementEnd
```

**GORM Model:**

```go
type UserPreferenceDataModel struct {
    ID                   uuid.UUID  `gorm:"column:id;type:uuid;primaryKey"`
    UserID               uuid.UUID  `gorm:"column:user_id;type:uuid"`
    EmailNotifications   bool       `gorm:"column:email_notifications;default:true"`
    SMSNotifications     bool       `gorm:"column:sms_notifications;default:false"`
    PushNotifications    bool       `gorm:"column:push_notifications;default:true"`
    ProfileVisibility    string     `gorm:"column:profile_visibility;default:public"`
    ShowEmail            bool       `gorm:"column:show_email;default:false"`
    ShowPhone            bool       `gorm:"column:show_phone;default:false"`
    Theme                string     `gorm:"column:theme;default:light"`
    Language             string     `gorm:"column:language;default:en"`
    Timezone             string     `gorm:"column:timezone;default:UTC"`
    CreatedAt            time.Time  `gorm:"column:created_at"`
    UpdatedAt            time.Time  `gorm:"column:updated_at"`
}

func (u *UserPreferenceDataModel) TableName() string {
    return "user_preferences"
}
```

---

### Example 10: Add Data Partitioning (00010)

**File:** `db/migrations/goose-migrate/00010_partition_users_by_created_date.sql`

**Scenario:** Partition large `users` table by creation date for better performance.

```sql
-- +goose Up
-- +goose StatementBegin
-- This is a complex operation - create partitioned table and migrate data
-- Typically done in phases during off-peak hours

-- Phase 1: Create partitioned table
CREATE TABLE IF NOT EXISTS users_partitioned (
    id UUID,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Create partitions for each year
CREATE TABLE users_2023 PARTITION OF users_partitioned
    FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');

CREATE TABLE users_2024 PARTITION OF users_partitioned
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

CREATE TABLE users_2025 PARTITION OF users_partitioned
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

-- Phase 2: Copy existing data
INSERT INTO users_partitioned
    SELECT * FROM users ON CONFLICT DO NOTHING;

-- Phase 3: Rename tables (requires brief downtime)
-- ALTER TABLE users RENAME TO users_old;
-- ALTER TABLE users_partitioned RENAME TO users;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop partitioned table (only if not in use)
-- DROP TABLE IF EXISTS users_partitioned CASCADE;
-- +goose StatementEnd
```

---

### Example 11: Add Foreign Key Constraint (00011)

**File:** `db/migrations/goose-migrate/00011_add_user_organization.sql`

**Scenario:** Add organizations and link users to organizations.

```sql
-- +goose Up
-- +goose StatementBegin
-- Create organizations table
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Add organization_id to users
ALTER TABLE users ADD COLUMN IF NOT EXISTS organization_id UUID;

-- Add foreign key constraint
ALTER TABLE users ADD CONSTRAINT IF NOT EXISTS fk_users_organization_id
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL;

-- Create index for join queries
CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id)
    WHERE deleted_at IS NULL;

-- Create index for organization queries
CREATE INDEX IF NOT EXISTS idx_organizations_name ON organizations(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_organizations_name;
DROP INDEX IF EXISTS idx_users_organization_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_organization_id;
ALTER TABLE users DROP COLUMN IF EXISTS organization_id;
DROP TABLE IF EXISTS organizations;
-- +goose StatementEnd
```

---

### Example 12: Add Unique Constraint with Conditions (00012)

**File:** `db/migrations/goose-migrate/00012_add_unique_email_per_org.sql`

**Scenario:** Ensure emails are unique within each organization (but same email OK across orgs).

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD CONSTRAINT IF NOT EXISTS unique_email_per_org
    UNIQUE (organization_id, email) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS unique_email_per_org;
-- +goose StatementEnd
```

---

## Quick Copy-Paste Template

Use this template for new migrations:

```sql
-- +goose Up
-- +goose StatementBegin
-- Your migration SQL here
-- Always use IF NOT EXISTS / IF EXISTS for safety
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Reverse of your migration SQL
-- Must be able to undo all changes
-- +goose StatementEnd
```

---

## Testing Migrations

### Verify Migration Applied

```bash
# Connect to database
psql -h localhost -p 54100 -U admin -d go-sandbox

# Check version
SELECT * FROM goose_db_version ORDER BY version_id;

# List tables
\dt

# View table structure
\d users

# View indexes
\di
```

### Rollback a Migration (Manual)

```sql
-- Connect to database
psql -h localhost -p 54100 -U admin -d go-sandbox

-- Execute the Down SQL from the migration file manually
-- Then update version
UPDATE goose_db_version SET version_id = X;  -- Set to previous version
```

---

## Performance Tips

1. **Add indexes on foreign keys**

   ```sql
   CREATE INDEX idx_fk_name ON table_name(foreign_key_column);
   ```

2. **Use partial indexes for soft deletes**

   ```sql
   CREATE INDEX idx_active_records ON table_name(id)
       WHERE deleted_at IS NULL;
   ```

3. **Create composite indexes for common WHERE clauses**

   ```sql
   CREATE INDEX idx_user_status_org ON users(status, organization_id)
       WHERE deleted_at IS NULL;
   ```

4. **Avoid large data type columns in indexes**

   ```sql
   -- Good: Index on ID
   CREATE INDEX idx_users_id ON users(id);

   -- Bad: Index on TEXT column
   CREATE INDEX idx_users_bio ON users(bio);  -- Don't do this!
   ```

5. **Use JSONB for flexible schemas**
   ```sql
   ALTER TABLE users ADD COLUMN metadata JSONB DEFAULT '{}'::JSONB;
   CREATE INDEX idx_metadata ON users USING GIN(metadata);
   ```

---

Done! Use these examples as templates for your own migrations. 🚀
