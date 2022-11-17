package dbshaker

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ToggyO/dbshaker/internal"
)

// DB represent a database connection driver.
type DB struct {
	connection *sql.DB
	dialect    internal.ISqlDialect
}

// OpenDBWithDriver creates a connection to a database, and creates
// compatible with the supplied driver by calling SQL dialect.
func OpenDBWithDriver(dialect, connectionString string) (*DB, error) {
	logger.Printf("Connecting to `%s` database...", dialect)

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
		connection: connection,
		dialect:    sqlDialect,
	}

	logger.Println("Connected to database!")

	return newDB, nil
}

// EnsureDBVersion retrieves the current version for this DB (major version, patch).
// Create and initialize the DB version table if it doesn't exist.
func EnsureDBVersion(db *DB) (int64, error) {
	return EnsureDBVersionContext(context.Background(), db)
}

// EnsureDBVersionContext retrieves the current version for this DB (major version, patch) with context.
// Create and initialize the DB version table if it doesn't exist.
func EnsureDBVersionContext(ctx context.Context, db *DB) (int64, error) {
	queryRunner := db.dialect.GetQueryRunner(ctx)
	err := db.dialect.CreateVersionTable(ctx, queryRunner)
	if err != nil {
		return 0, err
	}

	version, err := db.dialect.GetDBVersion(ctx, queryRunner)
	return version, nil
}

func ensureVersionTableExists(ctx context.Context, db *DB) error {
	return db.dialect.CreateVersionTable(ctx, db.dialect.GetQueryRunner(ctx))
}
