package dbshaker

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/ToggyO/dbshaker/shared"
)

// MigrationFunc migration action in database.
type MigrationFunc func(queryRunner shared.IQueryRunner) error

// Migration represents a database migration, manages by go runtime.
type Migration struct {
	Name    string // migration file name.
	Version int64  // version of migration.
	Patch   int16  // patch version of migration (increments when new migrations were applied,
	// but the greatest migration version in not changed)

	UpFn   MigrationFunc // Up migrations function.
	DownFn MigrationFunc // Down migrations function.

	Source    string // path to migration file.
	SourceDir string // path ti migration directory
	UseTx     bool   // indicate whether to run migration in transaction or not.
}

// Up executes an up migration.
func (m *Migration) Up(db *DB) error {
	return m.UpContext(context.Background(), db)
}

// UpContext executes an up migration with context.
func (m *Migration) UpContext(ctx context.Context, db *DB) error {
	return m.run(ctx, db, true)
}

// Down executes an up migration.
func (m *Migration) Down(db *DB) error {
	return m.DownContext(context.Background(), db)
}

// DownContext executes an up migration with context.
func (m *Migration) DownContext(ctx context.Context, db *DB) error {
	return m.run(ctx, db, false)
}

func (m *Migration) run(ctx context.Context, db *DB, direction bool) error {
	ext := filepath.Ext(m.Source)
	switch ext {
	case internal.SQLExt:
		file, err := os.Open(m.Source)
		if err != nil {
			// TODO: check error
			return fmt.Errorf("ERROR %v: failed to open SQL migration file: %w", filepath.Base(m.Source), err)
		}
		defer file.Close()

		statements, useTx, err := internal.ParseSQLMigration(bufio.NewReader(file), direction)
		if err != nil {
			// TODO: вынести
			return fmt.Errorf("ERROR %v: failed to parse SQL migration file: %w", filepath.Base(m.Source), err)
		}

		m.UseTx = useTx
		if err := m.runSQLMigration(ctx, db, statements, direction); err != nil {
			return fmt.Errorf("ERROR %v: failed to run SQL migration: %w", filepath.Base(m.Source), err)
		}

		// TODO: duplicatte
		if len(statements) > 0 {
			log.Printf("OK   %s \n", filepath.Base(m.Source))
		} else {
			log.Printf("EMPTY %s \n", filepath.Base(m.Source))
		}

	case internal.GoExt:
		if !m.UseTx {
			// TODO:
			//return db.dialect.Transaction(ctx, func(context context.Context, tx *sql.Tx) error {
			//	return m.runGoMigration(context, tx, db.dialect, direction)
			//})
			return m.runGoMigration(ctx, db.connection, db.dialect, direction)
		}

		return m.runGoMigration(ctx, nil, db.dialect, direction)
	}

	return nil
}

func (m *Migration) runGoMigration(ctx context.Context, queryRunner shared.IQueryRunner, dialect internal.ISqlDialect, direction bool) error {
	if queryRunner == nil {
		queryRunner = dialect.GetQueryRunner(ctx)
	}

	var err error

	fn := m.UpFn
	if !direction {
		fn = m.DownFn
	}

	if fn != nil {
		if err = fn(queryRunner); err != nil {
			//_ = tx.Rollback()
			return internal.ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
		}
	}

	if direction {
		if err = dialect.InsertVersion(ctx, queryRunner, m.Version); err != nil {
			return internal.ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
		}
	} else {
		if err = dialect.RemoveVersion(ctx, queryRunner, m.Version); err != nil {
			return internal.ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
		}
	}

	// TODO: duplicatte
	if fn != nil {
		log.Println("OK   ", filepath.Base(m.Name))
	} else {
		log.Println("EMPTY", filepath.Base(m.Name))
	}

	return err
}

func (m *Migration) runSQLMigration(ctx context.Context, db *DB, statements []string, direction bool) error {
	// TODO: add versioning
	if !m.UseTx {
		for _, statement := range statements {
			if _, err := db.connection.ExecContext(ctx, internal.ClearStatement(statement)); err != nil {
				return err
			}
		}

		// TODO: check
		// TODO: add db versioning
	}

	queryRunner := db.dialect.GetQueryRunner(ctx)
	for _, statement := range statements {
		if _, err := queryRunner.ExecContext(ctx, internal.ClearStatement(statement)); err != nil {
			return err
		}
	}

	// TODO: check
	//if direction {
	//	if err := db.dialect.InsertVersion(ctx, queryRunner, m.Version); err != nil {
	//		return internal.ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
	//	}
	//} else {
	//	if err := db.dialect.RemoveVersion(ctx, queryRunner, m.Version); err != nil {
	//		return internal.ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
	//	}
	//}

	// TODO: add db versioning
	//m.Version

	return nil
}

//func (m *Migration) run(ctx context.Context, tx *sql.Tx, dialect internal.ISqlDialect, direction bool) error {
//	ext := filepath.Ext(m.Name)
//	var err error
//
//	switch ext {
//	case internal.SQLExt:
//		file, err := os.Open(m.Source)
//		if err != nil {
//			// TODO: check error
//			return fmt.Errorf("ERROR %v: failed to open SQL migration file: %w", filepath.Base(m.Source), err)
//		}
//		defer file.Close()
//
//		statements, useTx, err := internal.ParseSQLMigration(bufio.NewReader(file), direction)
//		err = runSQLMigration(ctx, statements, useTx, m.Version, direction)
//
//	case internal.GoExt:
//		fn := m.UpFn
//		if !direction {
//			fn = m.DownFn
//		}
//
//		if fn != nil {
//			if err = fn(tx); err != nil {
//				_ = tx.Rollback()
//				return internal.ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
//			}
//		}
//
//		if direction {
//			if err = dialect.InsertVersion(ctx, m.Version); err != nil {
//				// TODO: check multiple rollback
//				_ = tx.Rollback()
//				return internal.ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
//			}
//		} else {
//			if err = dialect.RemoveVersion(ctx, m.Version); err != nil {
//				_ = tx.Rollback()
//				return internal.ErrFailedToRunMigration(filepath.Base(m.Name), fn, err)
//			}
//		}
//
//		if fn != nil {
//			log.Println("OK   ", filepath.Base(m.Name))
//		} else {
//			log.Println("EMPTY", filepath.Base(m.Name))
//		}
//
//		return nil
//	}
//
//	return nil
//}
