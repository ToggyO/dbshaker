package internal

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
)

// MigrationFunc migration action in database.
type MigrationFunc func(tx *sql.Tx) error

// Migration represents a database migration, manages by go runtime.
type Migration struct {
	Name    string // migration file name.
	Version int64  // version of migration.
	Patch   int16  // patch version of migration (increments when new migrations were applied,
	// but the greatest migration version in not changed)

	UpFn   MigrationFunc // Up migrations function.
	DownFn MigrationFunc // Down migrations function.

	Source string // path to migration file.
}

// Up executes an up migration.
func (m *Migration) Up(tx *sql.Tx, dialect ISqlDialect) error {
	return m.UpContext(context.Background(), tx, dialect)
}

// UpContext executes an up migration with context.
func (m *Migration) UpContext(ctx context.Context, tx *sql.Tx, dialect ISqlDialect) error {
	return m.run(ctx, tx, dialect, true)
}

// Down executes an up migration.
func (m *Migration) Down(tx *sql.Tx, dialect ISqlDialect) error {
	return m.DownContext(context.Background(), tx, dialect)
}

// DownContext executes an up migration with context.
func (m *Migration) DownContext(ctx context.Context, tx *sql.Tx, dialect ISqlDialect) error {
	return m.run(ctx, tx, dialect, false)
}

func (m *Migration) run(ctx context.Context, tx *sql.Tx, dialect ISqlDialect, direction bool) error {
	ext := filepath.Ext(m.Name)
	switch ext {
	case SQLExt:
	case GoExt:
		var err error

		fn := m.UpFn
		if !direction {
			fn = m.DownFn
		}

		if fn != nil {
			if err = fn(tx); err != nil {
				_ = tx.Rollback()
				return ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
			}
		}

		if direction {
			if err = dialect.InsertVersion(ctx, m.Version); err != nil {
				// TODO: check multiple rollback
				_ = tx.Rollback()
				return ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
			}
		} else {
			if err = dialect.RemoveVersion(ctx, m.Version); err != nil {
				_ = tx.Rollback()
				return ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
			}
		}

		if fn != nil {
			log.Println("OK   ", filepath.Base(m.Name))
		} else {
			log.Println("EMPTY", filepath.Base(m.Name))
		}

		return nil
	}

	return nil
}
