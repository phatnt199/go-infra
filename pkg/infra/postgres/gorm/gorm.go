package postgresgorm

import (
	"database/sql"
	"fmt"

	"emperror.dev/errors"
	defaultlogger "github.com/phatnt199/go-infra/pkg/logger/default_logger"
	gromlog "github.com/phatnt199/go-infra/pkg/logger/external/gormlog"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultPostgresDB = "postgres"
	DNSFormat         = "host=%s port=%d user=%s dbname=%s password=%s sslmode=disable"
)

type BuildPostgresDSNArgs struct {
	host     string
	port     int
	user     string
	dbName   string
	password string
}

func NewGorm(cfg *GormOptions) (*gorm.DB, error) {
	if cfg == nil {
		return nil, errors.New("Gorm configuration cannot be nil")
	}
	if cfg.DBName == "" {
		return nil, errors.New("Gorm database name cannot be empty")
	}

	switch cfg.Type {
	case InMemory:
		return nil, errors.New("Gorm In-Memory database type not yet supported")
	case SQLite:
		return nil, errors.New("Gorm SQLite database type not yet supported")
	case Postgres:
		if err := createPostgresDB(cfg); err != nil {
			return nil, errors.Wrap(err, "Failed to create Postgres database")
		}
	default:
		return nil, errors.New("Unsupported Gorm database type")
	}

	gorm, err := OpenPostgresConnection(BuildPostgresDSNArgs{
		host:     cfg.Host,
		port:     cfg.Port,
		dbName:   cfg.DBName,
		user:     cfg.User,
		password: cfg.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open Postgres connection")
	}

	return gorm, nil
}

func buildPostgresDSN(opts BuildPostgresDSNArgs) string {
	return fmt.Sprintf(
		DNSFormat,
		opts.host,
		opts.port,
		opts.user,
		opts.dbName,
		opts.password,
	)
}

func OpenPostgresConnection(opts BuildPostgresDSNArgs) (*gorm.DB, error) {
	dsn := buildPostgresDSN(opts)

	gormDB, err := gorm.Open(
		gormPostgres.Open(dsn),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultlogger.GetLogger()),
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to Postgres default db")
	}

	return gormDB, err
}

func createPostgresDB(cfg *GormOptions) error {
	// Connect to default Postgres database to create the target database if it doesn't exist
	postgresGormDB, err := OpenPostgresConnection(BuildPostgresDSNArgs{
		host:     cfg.Host,
		port:     cfg.Port,
		dbName:   defaultPostgresDB,
		user:     cfg.User,
		password: cfg.Password,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to connect to Postgres default db")
	}

	// Get sql.DB from gorm DB
	db, err := postgresGormDB.DB()
	if err != nil {
		return errors.Wrap(err, "Failed to get sql.DB from gorm DB")
	}
	defer db.Close()

	// Check if the target database exists
	exists, err := checkDatabaseExists(db, cfg.DBName)
	if err != nil {
		return errors.Wrap(err, "Failed to check if database exists")
	}

	// Create the database if it does not exist
	if !exists {
		// Create the target database
		if err := createDatabase(db, cfg.DBName); err != nil {
			return errors.Wrap(err, "Failed to create database")
		}
	}

	return nil
}

func checkDatabaseExists(db *sql.DB, dbName string) (bool, error) {
	query := "SELECT 1 FROM pg_catalog.pg_database WHERE datname = $1"

	rows, err := db.Query(query, dbName)
	if err != nil {
		return false, errors.Wrap(err, "Failed to query pg_catalog.pg_database")
	}
	defer rows.Close()

	if rows.Next() {
		var exists int
		if err := rows.Scan(&exists); err != nil {
			return false, errors.Wrap(err, "Failed to scan pg_catalog.pg_database")
		}
		return exists == 1, nil
	}

	if err := rows.Err(); err != nil {
		return false, errors.Wrap(err, "Failed to iterate over pg_catalog.pg_database")
	}

	return false, nil
}

func createDatabase(db *sql.DB, dbName string) error {
	query := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err := db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "Failed to create database %s", dbName)
	}

	return nil
}
