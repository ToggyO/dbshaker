package dbshaker

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/ToggyO/dbshaker/internal"
)

const defaultLockTimeout = 15 * time.Second

// DB represent a database connection driver.
type DB struct {
	connection *sql.DB
	dialect    internal.ISqlDialect
	mu         *sync.Mutex
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
		mu:         &sync.Mutex{},
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
	return db.dialect.GetDBVersion(ctx, queryRunner)
}

func lockDB(ctx context.Context, db *DB) error {
	// create done channel, used in the timeout goroutine
	done := make(chan bool, 1)
	defer func() {
		done <- true
	}()

	// use errChan to signal error back to this context
	errChan := make(chan error, 2)

	// TODO: настроить конфигурирование
	timeout := time.After(defaultLockTimeout)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-timeout:
				errChan <- internal.ErrLockTimeout
			}
		}
	}()

	// now try to acquire the lock
	go func() {
		err := db.dialect.Lock(ctx)
		if err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	// wait until we either receive ErrLockTimeout or error from Lock operation
	return <-errChan
}
