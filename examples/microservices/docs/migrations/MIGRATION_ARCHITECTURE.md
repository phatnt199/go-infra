# User Service - Migration Architecture & Diagrams

Visual representations and architecture documentation for the migration system.

---

## 🏗️ Architecture Overview

### Migration System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    User Service Application                      │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ├─ main.go
                 │  │
                 └─ App.Run()
                    │
                    └─ CreateFxApp() (Dependency Injection)
                       │
                       ├─ UsersServiceConfigurator.ConfigureUsersService()
                       │  │
                       │  ├─ migrateUsers()
                       │  │  │
                       │  │  └─ migrateGoose()
                       │  │     │
                       │  │     └─ PostgresMigrationRunner.Up()
                       │  │        │
                       │  │        └─ Goose Library
                       │  │           │
                       │  │           ├─ Read config from config.development.json
                       │  │           │
                       │  │           ├─ Connect to PostgreSQL
                       │  │           │
                       │  │           ├─ Query goose_db_version table
                       │  │           │
                       │  │           ├─ Read migration files from:
                       │  │           │  db/migrations/goose-migrate/*.sql
                       │  │           │
                       │  │           ├─ Execute pending migrations
                       │  │           │
                       │  │           └─ Update goose_db_version
                       │  │
                       │  ├─ seedUsers()
                       │  │  │
                       │  │  └─ seedDataManually()
                       │  │     │
                       │  │     ├─ Check if data exists
                       │  │     │
                       │  │     └─ Insert seed data via GORM
                       │  │
                       │  └─ configSwagger()
                       │
                       └─ HTTP Server Listen on :8080
```

---

## 📊 Database Migration Flow

### Sequential Migration Execution

```
Step 1: Application Start
  ↓
Step 2: Load Configuration
  ├─ Database connection details
  ├─ Migration directory: db/migrations/goose-migrate
  └─ Skip migration flag: false
  ↓
Step 3: Connect to PostgreSQL
  ├─ Host: localhost
  ├─ Port: 54100
  ├─ Database: go-sandbox
  └─ User: admin
  ↓
Step 4: Check goose_db_version Table
  │
  ├─ If empty: Current version = 0
  │
  ├─ If has records: Current version = MAX(version_id)
  │
  └─ Example: version_id = 1
  ↓
Step 5: Scan Migration Files
  │
  ├─ Find all *.sql files in db/migrations/goose-migrate/
  │
  ├─ Parse version numbers:
  │  ├─ 00001_create_users_table.sql → Version 1
  │  ├─ 00002_add_profile_fields.sql → Version 2
  │  └─ 00003_add_user_status.sql → Version 3
  │
  └─ Sort by version number (ascending)
  ↓
Step 6: Identify Pending Migrations
  │
  ├─ Current version: 1
  │
  ├─ Available versions: 1, 2, 3
  │
  ├─ Pending migrations: 2, 3
  │
  └─ Process: Execute files with version > 1
  ↓
Step 7: Execute Pending Migrations (In Order)
  │
  ├─ Migration #2: 00002_add_profile_fields.sql
  │  ├─ Parse -- +goose Up section
  │  ├─ Execute SQL
  │  ├─ If success: INSERT INTO goose_db_version (version_id=2, is_applied=true)
  │  └─ Update current version to 2
  │
  ├─ Migration #3: 00003_add_user_status.sql
  │  ├─ Parse -- +goose Up section
  │  ├─ Execute SQL
  │  ├─ If success: INSERT INTO goose_db_version (version_id=3, is_applied=true)
  │  └─ Update current version to 3
  │
  └─ Complete: All pending migrations applied
  ↓
Step 8: Run Seeding (After Migrations)
  │
  ├─ seedUsers()
  │
  ├─ Check: COUNT(*) FROM users
  │
  ├─ If count == 0:
  │  └─ INSERT seed data
  │
  └─ If count > 0:
     └─ Skip (idempotent)
  ↓
Step 9: Application Ready
  │
  └─ Start HTTP server on :8080
     ✓ Migrations complete
     ✓ Database schema up-to-date
     ✓ Seed data loaded
     ✓ Ready to handle requests
```

---

## 📁 File Structure

### Directory Layout

```
userservice/
│
├── 📄 MIGRATION_GUIDE.md (Comprehensive guide)
├── 📄 MIGRATION_QUICK_START.md (Quick reference)
├── 📄 MIGRATION_EXAMPLES.md (Copy-paste examples)
├── 📄 MIGRATION_ARCHITECTURE.md (This file)
│
├── 📂 db/
│   └── 📂 migrations/
│       └── 📂 goose-migrate/
│           ├── 00001_create_users_table.sql
│           ├── 00002_add_profile_fields.sql
│           ├── 00003_add_user_status.sql
│           └── ... (more migrations)
│
├── 📄 config.development.json
├── 📄 config.production.json
├── 📄 config.staging.json
│
├── 📂 internal/
│   └── 📂 shared/
│       └── 📂 configurations/
│           └── 📂 users/
│               ├── users_configurator.go
│               ├── users_configurator_migration.go
│               └── users_configurator_seed.go
│
└── 📂 cmd/
    └── 📂 app/
        └── main.go
```

---

## 🔄 State Transitions

### Migration State Diagram

```
                    ┌─────────────────────┐
                    │  Application Start  │
                    └──────────┬──────────┘
                               │
                      ┌────────▼────────┐
                      │  Load Config    │
                      └────────┬────────┘
                               │
                      ┌────────▼────────┐
                      │ Connect to DB   │
                      └────────┬────────┘
                               │
                      ┌────────▼──────────────────┐
                      │ Query goose_db_version   │
                      │ Get current_version      │
                      └────────┬──────────────────┘
                               │
          ┌────────────────────┴────────────────────┐
          │                                         │
    ┌─────▼──────┐                        ┌────────▼────────┐
    │ Table Empty │                        │ Has Records     │
    │  version=0  │                        │ version=N       │
    └─────┬──────┘                        └────────┬────────┘
          │                                         │
          └──────────────┬──────────────────────────┘
                         │
                ┌────────▼────────────┐
                │ Scan Migration      │
                │ Files in Directory  │
                └────────┬────────────┘
                         │
        ┌────────────────┴────────────────┐
        │                                 │
    ┌───▼──────┐                    ┌────▼───────┐
    │ Pending? │                    │ No Pending │
    │ YES ✓    │                    │ ✓ DONE     │
    └───┬──────┘                    └────────────┘
        │
    ┌───▼──────────────────┐
    │ Execute Next Migration
    │ (-- +goose Up)       │
    └───┬──────────────────┘
        │
    ┌───▼─────────────────────┐
    │ Success?                │
    └───┬─────────────────────┘
        │
    ┌───┴───────────┬──────────────┐
    │               │              │
YES │            NO │              │
    │               │              │
┌───▼──────────┐ ┌──▼────────────┐ │
│ Record in    │ │ Panic/Error   │ │
│ DB Version   │ └───────────────┘ │
│ More pending?│                    │
└───┬──────────┘                    │
    │                              │
    ├─── YES ──────────────────────┘
    │  (Loop back to exec next)
    │
    └─── NO
        │
        ┌────────────────────┐
        │ All Migrations OK  │
        └────────┬───────────┘
                 │
        ┌────────▼─────────┐
        │ Run Seeding      │
        │ seedUsers()      │
        └────────┬─────────┘
                 │
        ┌────────▼──────────────┐
        │ Start HTTP Server     │
        │ Listen on :8080       │
        └───────────────────────┘
```

---

## 🔍 Goose Version Tracking

### goose_db_version Table

```
┌────────────┬────────────┬─────────────────────────────┐
│ version_id │ is_applied │            tstamp            │
├────────────┼────────────┼─────────────────────────────┤
│     1      │     t      │ 2025-11-05 14:35:56.456713 │
│     2      │     t      │ 2025-11-05 14:36:12.789456 │
│     3      │     t      │ 2025-11-05 14:36:28.123789 │
└────────────┴────────────┴─────────────────────────────┘

t = true (migration applied successfully)
f = false (migration rolled back)
```

---

## 📝 Migration File Parsing

### SQL File Structure Parsing

```
File: 00001_create_users_table.sql

┌────────────────────────────────────────────────┐
│ -- +goose Up                                   │ ← Marker
│ -- +goose StatementBegin                       │ ← Begin
│ CREATE TABLE users ( ... );                    │
│ CREATE INDEX idx_users_email ON users(...);    │ ← Forward SQL
│ -- +goose StatementEnd                         │ ← End
│                                                │
│ -- +goose Down                                 │ ← Marker
│ -- +goose StatementBegin                       │ ← Begin
│ DROP TABLE users;                              │ ← Reverse SQL
│ -- +goose StatementEnd                         │ ← End
└────────────────────────────────────────────────┘

Goose Parser:
  1. Reads entire file
  2. Splits by "-- +goose Up" and "-- +goose Down"
  3. Extracts SQL between StatementBegin/End
  4. For UP migration: Execute forward SQL
  5. For DOWN migration: Execute reverse SQL (when rolling back)
```

---

## 🎯 Seeding Process

### Seed Data Flow

```
seedUsers() Called
       │
       ├─ Check if data exists
       │
       ├─────────────────────────────────────┐
       │                                     │
   YES │ COUNT > 0                       NO │ COUNT = 0
       │                                     │
       ▼                                     ▼
   ┌─────────────┐                    ┌──────────────┐
   │  Return OK  │                    │ Create Data  │
   │  (Idempotent)                    │ Objects      │
   └─────────────┘                    └──────┬───────┘
                                             │
                                    ┌────────▼──────────┐
                                    │ Generate UUIDs    │
                                    │ for Each Record   │
                                    └────────┬──────────┘
                                             │
                                    ┌────────▼──────────────┐
                                    │ db.CreateInBatches() │
                                    │ (Insert via GORM)    │
                                    └────────┬──────────────┘
                                             │
                                    ┌────────▼──────────┐
                                    │   Success?        │
                                    └────────┬──────────┘
                                             │
                                 ┌───────────┴───────────┐
                                 │                       │
                            YES  │                   NO  │
                                 │                       │
                            ┌────▼────┐          ┌──────▼───┐
                            │ Return  │          │ Return   │
                            │ Error   │          │ Error    │
                            └─────────┘          └──────────┘
```

---

## 🔐 Database Connection Flow

### Connection Configuration

```
config.development.json
│
├─ migrationOptions.host: "localhost"
├─ migrationOptions.port: 54100
├─ migrationOptions.user: "admin"
├─ migrationOptions.password: "123456"
├─ migrationOptions.dbName: "go-sandbox"
└─ migrationOptions.sslMode: false
   │
   └─ Build Connection String
      │
      ├─ Format: postgres://user:password@host:port/database
      │
      └─ Result: postgres://admin:123456@localhost:54100/go-sandbox
         │
         └─ PostgreSQL Driver (GORM)
            │
            ├─ Parse DSN
            ├─ Validate credentials
            ├─ Establish TCP connection
            ├─ Authenticate
            └─ Ready for queries
```

---

## 🧩 Integration Points

### Component Integration

```
┌─────────────────────────────────────────────────────────┐
│                    main.go                              │
│                                                         │
│  ├─ app.Run()                                          │
│  └─ UsersServiceConfigurator                           │
│     │                                                   │
│     ├─ ConfigureInfrastructure()  ──────────┐          │
│     │  ├─ GORM Config                       │          │
│     │  ├─ HTTP Adapter                      │          │
│     │  └─ Goose Migration Setup             │          │
│     │                                       │          │
│     ├─ ConfigureUsersService()  ◄──────────┤────────┐ │
│     │  ├─ migrateUsers()         (Calls)   │        │ │
│     │  │  └─ PostgresMigrationRunner.Up()  │  SQL  │ │
│     │  │                                    │ Files │ │
│     │  ├─ seedUsers()                      │        │ │
│     │  └─ ResolveFunc() (DI container)     │        │ │
│     │                                       │        │ │
│     └─ MapUsersEndpoints()                 │        │ │
│        └─ Setup HTTP routes                │        │ │
│                                            │        │ │
│     ┌────────────────────────────────────┐ │        │ │
│     │ PostgreSQL Database                │ │        │ │
│     │                                    │ │        │ │
│     │ ├─ goose_db_version              │ │        │ │
│     │ ├─ users                         │◄┼────────┤ │
│     │ ├─ user_profiles (if migration)  │ │        │ │
│     │ └─ user_preferences (if migration)│◄┼────────┘ │
│     │                                    │ │          │
│     └────────────────────────────────────┘ │          │
│                                            │          │
│         db/migrations/goose-migrate/ ◄─────┘          │
│         ├─ 00001_create_users_table.sql              │
│         ├─ 00002_add_profile_fields.sql              │
│         └─ 00003_add_user_status.sql                 │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## 🚀 Deployment Flow

### Production Deployment

```
┌──────────────────────────────┐
│ New Service Version Released │
│ (Contains new migration file)│
└──────────┬───────────────────┘
           │
    ┌──────▼──────────┐
    │ Start Service   │
    │ (New Version)   │
    └──────┬──────────┘
           │
    ┌──────▼─────────────────────────┐
    │ Read Migration Directory        │
    │ - Old migrations: 00001, 00002  │
    │ - New migration: 00003          │
    └──────┬──────────────────────────┘
           │
    ┌──────▼─────────────────────────┐
    │ Query goose_db_version          │
    │ Current version: 2              │
    └──────┬──────────────────────────┘
           │
    ┌──────▼─────────────────────────┐
    │ Identify Pending: 3             │
    └──────┬──────────────────────────┘
           │
    ┌──────▼─────────────────────────┐
    │ Execute 00003 migration         │
    │ (Zero downtime if supported)    │
    └──────┬──────────────────────────┘
           │
    ┌──────▼─────────────────────────┐
    │ Update goose_db_version to 3    │
    └──────┬──────────────────────────┘
           │
    ┌──────▼─────────────────────────┐
    │ Service Ready                   │
    │ Serving requests                │
    └─────────────────────────────────┘

✓ Zero downtime (if migrations are backward compatible)
✓ Automatic deployment
✓ Version tracked in database
✓ Rollback possible via code + manual DB update
```

---

## 📊 Common Scenarios

### Scenario 1: Adding a Column

```
Commit: "feat: add user.phone column"
├─ Create: db/migrations/goose-migrate/00004_add_phone_to_users.sql
├─ Update: UserDataModel struct
├─ Deploy to production
└─ On first run:
   ├─ Read 00004 migration
   ├─ Execute: ALTER TABLE users ADD COLUMN phone VARCHAR(20)
   ├─ Record version 4 in goose_db_version
   └─ Existing code works seamlessly (new column nullable by default)
```

### Scenario 2: Creating Related Table

```
Commit: "feat: add user preferences table"
├─ Create: db/migrations/goose-migrate/00005_create_user_preferences.sql
├─ Create: UserPreferenceDataModel
├─ Create: User-Preference Repository
├─ Update: API endpoints for preferences
├─ Deploy
└─ On first run:
   ├─ Execute: CREATE TABLE user_preferences...
   ├─ Execute: CREATE INDEX...
   ├─ Record version 5 in goose_db_version
   └─ New features available
```

### Scenario 3: Data Type Change

```
Commit: "refactor: change phone type from varchar to structured format"
├─ Create: db/migrations/goose-migrate/00006_refactor_phone_column.sql
│  ├─ UP: Migrate data, change type
│  └─ DOWN: Revert changes
├─ Update: UserDataModel phone field type
├─ Deploy
└─ On first run:
   ├─ Execute migration in transaction
   ├─ Record version 6
   └─ Backward compatible with old code for brief period
```

---

## 🎓 Learning Path

1. **Start here**: [MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md)

   - Basic concepts
   - Running the service
   - Creating your first migration

2. **Then read**: [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md)

   - Deep dive into architecture
   - Detailed steps for common tasks
   - Best practices
   - Troubleshooting

3. **Reference**: [MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md)

   - Copy-paste ready examples
   - Real-world scenarios
   - Advanced patterns

4. **Understand**: [MIGRATION_ARCHITECTURE.md](./MIGRATION_ARCHITECTURE.md) (this file)
   - System architecture
   - Data flows
   - Integration points

---

## 🎯 Key Takeaways

1. ✅ Goose handles version tracking automatically
2. ✅ Migrations run sequentially in version order
3. ✅ All changes are reversible (Down sections)
4. ✅ Seeding is idempotent (can run multiple times safely)
5. ✅ Configuration in `config.development.json`
6. ✅ SQL files in `db/migrations/goose-migrate/`
7. ✅ Automatic on service startup
8. ✅ Database state always tracked in `goose_db_version`

---

Need help? Refer to the specific guide:

- 🚀 Quick start: [MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md)
- 📖 Comprehensive: [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md)
- 💡 Examples: [MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md)
