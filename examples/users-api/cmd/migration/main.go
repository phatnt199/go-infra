package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/phatnt199/go-infra/pkg/logger"
	defaultlogger "github.com/phatnt199/go-infra/pkg/logger/default_logger"
	"github.com/phatnt199/go-infra/pkg/migration"
	"github.com/phatnt199/go-infra/pkg/migration/gomigrate"

	_ "github.com/lib/pq"
)

func main() {
	// Load .env file (optional - will use system env if not found)
	// _ = godotenv.Load("../../.env")

	// Parse command line flags
	direction := flag.String("direction", "up", "Migration direction: up, down, or status")
	version := flag.Uint("version", 0, "Target migration version (0 = all)")
	skip := flag.Bool("skip", false, "Skip migration")
	flag.Parse()

	// Initialize logger
	log := defaultlogger.GetLogger()

	log.Info("Database Migration Tool")

	// Load configuration from environment
	migrationConfig := &migration.MigrationOptions{
		Host:          getEnv("MIGRATION_OPTIONS_HOST", "localhost"),
		Port:          getEnvAsInt("MIGRATION_OPTIONS_PORT", 5432),
		User:          getEnv("MIGRATION_OPTIONS_USER", "postgres"),
		Password:      getEnv("MIGRATION_OPTIONS_PASSWORD", "postgres"),
		DBName:        getEnv("MIGRATION_OPTIONS_DBNAME", "usersdb"),
		MigrationsDir: getEnv("MIGRATION_OPTIONS_MIGRATIONS_DIR", "migrations"),
		SkipMigration: *skip,
	}

	if migrationConfig.SkipMigration {
		log.Info("Skipping database migrations")
		return
	}

	log.Infow("Migration configuration loaded", logger.Fields{
		"host":   migrationConfig.Host,
		"port":   migrationConfig.Port,
		"dbname": migrationConfig.DBName,
		"dir":    migrationConfig.MigrationsDir,
	})

	// Connect to database
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		migrationConfig.User,
		migrationConfig.Password,
		migrationConfig.Host,
		migrationConfig.Port,
		migrationConfig.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Info("Database connection established")

	// Create migrator
	migrator, err := gomigrate.NewGoMigratorPostgres(migrationConfig, db, log)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	// Execute migration
	ctx := context.Background()
	switch *direction {
	case "up":
		log.Infow("Running migrations UP", logger.Fields{
			"version": *version,
		})
		if err := migrator.Up(ctx, *version); err != nil {
			log.Fatalf("Migration UP failed: %v", err)
		}
		log.Info("Migrations UP completed successfully")

	case "down":
		log.Infow("Running migrations DOWN", logger.Fields{
			"version": *version,
		})
		if err := migrator.Down(ctx, *version); err != nil {
			log.Fatalf("Migration DOWN failed: %v", err)
		}
		log.Info("Migrations DOWN completed successfully")

	case "status":
		fmt.Println("Getting migration version...")
		version, dirty, err := migrator.Version(ctx)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Migration version: %d, dirty: %t\n", version, dirty)

	default:
		log.Fatalf("Invalid direction: %s (use 'up', 'down', or 'status')", *direction)
	}

	log.Info("Migration tool completed successfully")
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
