# User Service - Migration Documentation Index

Complete migration documentation for the User Service microservice.

## 📚 Documentation Files

| Document                                                     | Purpose                          | Audience   | Time   |
| ------------------------------------------------------------ | -------------------------------- | ---------- | ------ |
| **[MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md)**   | Quick reference for common tasks | Everyone   | 5 min  |
| **[MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md)**               | Comprehensive migration guide    | Developers | 30 min |
| **[MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md)**         | Real-world copy-paste examples   | Developers | 20 min |
| **[MIGRATION_ARCHITECTURE.md](./MIGRATION_ARCHITECTURE.md)** | System architecture & diagrams   | Tech leads | 15 min |

---

## 🎯 Quick Navigation

### I want to...

#### **Get Started Quickly**

→ Read [MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md)

- 🚀 Start the service
- 📝 Create your first migration
- 🔍 Check migration status

#### **Understand How It Works**

→ Read [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) → Overview & Architecture sections

- How migrations run on startup
- How Goose version tracking works
- Configuration and file structure

#### **Add a New Column to Users Table**

→ Read [MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md) → "Add User Profile Fields" example

1. Copy the migration SQL template
2. Modify for your column
3. Update GORM model
4. Restart service

#### **Create a Completely New Table**

→ Read [MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md) → "Create Users Table" or "Create User Preferences Table" examples

1. Use the appropriate example template
2. Customize for your table
3. Create corresponding GORM model
4. Restart service

#### **Handle Schema Changes (rename, type change)**

→ Read [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) → Common Tasks section
→ Or [MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md) → Data Transformations section

- Rename column
- Change column type
- Add/remove constraints

#### **Seed Initial Data**

→ Read [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) → Task 6: Seed Initial Data
→ Or [MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md) → Seed Data section

1. Edit `users_configurator_seed.go`
2. Add records to the slice
3. Restart service

#### **Understand the System Architecture**

→ Read [MIGRATION_ARCHITECTURE.md](./MIGRATION_ARCHITECTURE.md)

- Component diagrams
- Data flow diagrams
- State transition diagrams
- Integration points

#### **Fix an Error**

→ Read [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) → Troubleshooting section

- Migration file not found
- Duplicate version error
- Column already exists
- Database connection issues

---

## 📖 Learning Paths

### Path 1: First-Time Developer (30 minutes)

1. **[MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md)** (5 min)

   - Understand basic concepts
   - Run your first migration

2. **[MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md)** (15 min)

   - Review "Add User Profile Fields" example
   - See how it's done in practice

3. **[MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md)** (10 min)
   - Read Overview section
   - Skim Common Tasks section

### Path 2: Experienced Developer (15 minutes)

1. **[MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md)** (3 min)

   - Quick reference check

2. **[MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md)** (12 min)
   - Browse examples relevant to your needs
   - Copy templates

### Path 3: System Architect (45 minutes)

1. **[MIGRATION_ARCHITECTURE.md](./MIGRATION_ARCHITECTURE.md)** (15 min)

   - System design and flows
   - Integration points

2. **[MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md)** (20 min)

   - Architecture section
   - How migrations work section
   - Best practices

3. **[MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md)** (10 min)
   - Advanced patterns section

---

## 🎓 Key Concepts

### What is a Migration?

A migration is a versioned SQL script that transforms your database schema. Each migration has:

- **UP section**: Code to apply the change
- **DOWN section**: Code to reverse the change

### How Goose Works

```
Service Starts
    ↓
Read migrations from: db/migrations/goose-migrate/
    ↓
Check goose_db_version table for current version
    ↓
Execute all migrations with version > current_version
    ↓
Update goose_db_version with new version
    ↓
Run seeding if needed
    ↓
Service Ready
```

### File Naming

```
00001_create_users_table.sql
├─ 00001: Version number (5 digits, zero-padded)
└─ create_users_table: Human-readable description (snake_case)
```

### Migration Structure

```sql
-- +goose Up
-- +goose StatementBegin
-- Your forward SQL here
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Your reverse SQL here
-- +goose StatementEnd
```

---

## 💾 Configuration

Located in `config.development.json` (and other config files):

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

---

## 📁 Project Structure

```
userservice/
├── MIGRATION_GUIDE.md .................. Comprehensive guide
├── MIGRATION_QUICK_START.md ........... Quick reference
├── MIGRATION_EXAMPLES.md .............. Copy-paste examples
├── MIGRATION_ARCHITECTURE.md ......... System design
├── MIGRATION_README.md ................ This file
│
├── config.development.json ............ Configuration (dev)
├── config.production.json ............ Configuration (prod)
│
├── db/migrations/goose-migrate/
│   ├── 00001_create_users_table.sql
│   └── ... (more migrations)
│
├── internal/shared/configurations/users/
│   ├── users_configurator_migration.go . Triggers migrations
│   ├── users_configurator_seed.go ..... Populates initial data
│   └── users_configurator.go ......... Main configuration
│
└── cmd/app/
    └── main.go ..................... Application entry point
```

---

## 🚀 Getting Started

### 1. Run the Service

```bash
cd examples/microservices/userservice
go run ./cmd/app/main.go
```

Expected output:

```
2025/11/05 14:35:56 OK   00001_create_users_table.sql (4.6ms)
2025/11/05 14:35:56 goose: successfully migrated database to version: 1
```

### 2. Create a Migration

Create file: `db/migrations/goose-migrate/00002_description.sql`

```sql
-- +goose Up
-- +goose StatementBegin
-- Your SQL here
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Your reverse SQL here
-- +goose StatementEnd
```

### 3. Restart Service

The migration runs automatically on startup.

---

## ✨ Best Practices

1. ✅ **Always use IF EXISTS/IF NOT EXISTS** - Makes migrations idempotent
2. ✅ **Make migrations reversible** - Every UP needs a matching DOWN
3. ✅ **Use meaningful names** - `00001_create_users_table.sql` not `00001_fix.sql`
4. ✅ **Test migrations locally** - Before committing to version control
5. ✅ **One concern per migration** - Don't mix unrelated changes
6. ✅ **Add comments** - Explain the "why" not just the "what"

---

## 🔍 Common Tasks

### Add a Column

→ See [MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md) → Add a Column

### Create a Table

→ See [MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md) → Initial Setup

### Add an Index

→ See [MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md) → Create an Index

### Seed Data

→ See [MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md) → Add Seed Data

### Fix an Error

→ See [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) → Troubleshooting

---

## 🎯 FAQ

**Q: Will migrations run automatically?**
A: Yes! Migrations run on service startup if `skipMigration` is `false` in config.

**Q: Can I rollback a migration?**
A: Yes, by manually executing the DOWN SQL and updating the `goose_db_version` table.

**Q: What if I make a mistake in a migration?**
A: Create a new migration to fix it (never modify already-applied migrations).

**Q: Do I need to update GORM models?**
A: Yes, after schema changes, update the corresponding datamodel struct.

**Q: Is seeding automatic?**
A: Yes, but only runs if the table is empty (idempotent).

**Q: Can migrations have transaction issues?**
A: Each migration file is executed as a single transaction, so all-or-nothing.

**Q: Where are migration files stored?**
A: `db/migrations/goose-migrate/` (configurable in config file)

**Q: What tracks applied migrations?**
A: Goose creates a `goose_db_version` table in your database.

---

## 📞 Support

**Found a bug?**

- Check [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) → Troubleshooting

**Have a question?**

- Search [MIGRATION_EXAMPLES.md](./MIGRATION_EXAMPLES.md) for similar example

**Want to understand the system?**

- Read [MIGRATION_ARCHITECTURE.md](./MIGRATION_ARCHITECTURE.md)

**Need a quick reference?**

- Use [MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md)

---

## 📚 External Resources

- **Goose Documentation**: https://github.com/pressly/goose
- **PostgreSQL Documentation**: https://www.postgresql.org/docs/
- **GORM Documentation**: https://gorm.io/docs/

---

## 🎉 Summary

The User Service uses **Goose** for database migrations:

| Aspect            | Details                                    |
| ----------------- | ------------------------------------------ |
| **Tool**          | Goose v3 (SQL-based migrations)            |
| **Versioning**    | Tracked in `goose_db_version` table        |
| **Execution**     | Automatic on service startup               |
| **Rollback**      | Manual (execute DOWN SQL + update version) |
| **Storage**       | `db/migrations/goose-migrate/`             |
| **Configuration** | `config.*.json`                            |
| **Seeding**       | GORM-based, idempotent                     |
| **Status**        | ✅ Fully automated & production-ready      |

---

## 📖 Document Versions

| Document                  | Version | Last Updated |
| ------------------------- | ------- | ------------ |
| MIGRATION_README.md       | 1.0     | 2025-11-05   |
| MIGRATION_QUICK_START.md  | 1.0     | 2025-11-05   |
| MIGRATION_GUIDE.md        | 1.0     | 2025-11-05   |
| MIGRATION_EXAMPLES.md     | 1.0     | 2025-11-05   |
| MIGRATION_ARCHITECTURE.md | 1.0     | 2025-11-05   |

---

**Ready to get started? Start with [MIGRATION_QUICK_START.md](./MIGRATION_QUICK_START.md)!** 🚀
