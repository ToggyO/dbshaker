package dbshaker

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ToggyO/dbshaker/internal"
)

// DB represent a database connection driver.
type DB struct {
	db      *sql.DB
	dialect internal.ISqlDialect
}

// OpenDBWithDriver creates a connection to a database, and creates
// compatible with the supplied driver by calling SQL dialect.
func OpenDBWithDriver(dialect, connectionString string) (*DB, error) {
	fmt.Printf("Connecting to `%s` database...", dialect)

	var connection *sql.DB
	var err error

	switch dialect {
	// tODO: check
	// case "postgres", "pgx", "sqlite3", "sqlite", "mysql", "sqlserver":
	case internal.PostgresDialect, internal.PgxDialect:
		connection, err = sql.Open(dialect, connectionString)
	default:
		return nil, fmt.Errorf("unsupported driver '%s'", dialect)
	}

	if err != nil {
		return nil, err
	}

	if err = connection.Ping(); err != nil {
		return nil, fmt.Errorf("ERROR: failed connect to database: %w", err)
	}

	sqlDialect, err := createDialect(connection, dialect)
	if err != nil {
		return nil, err
	}

	newDB := &DB{
		db:      connection,
		dialect: sqlDialect,
	}

	fmt.Println("Connected to database!")

	return newDB, nil
}

// EnsureDBVersion retrieves the current version for this DB (major version, patch).
// Create and initialize the DB version table if it doesn't exist.
func EnsureDBVersion(db *DB) (int64, byte, error) {
	return EnsureDBVersionContext(context.Background(), db)
}

// EnsureDBVersionContext retrieves the current version for this DB (major version, patch) with context.
// Create and initialize the DB version table if it doesn't exist.
func EnsureDBVersionContext(ctx context.Context, db *DB) (int64, byte, error) {
	version, err := db.dialect.GetDBVersion(ctx)
	if err != nil {
		return version.Version, version.Patch, db.dialect.CreateVersionTable(ctx)
	}

	return version.Version, version.Patch, nil
}
