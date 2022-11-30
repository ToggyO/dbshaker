package dbshaker

import (
	"bufio"
	"context"
	"github.com/ToggyO/dbshaker/internal/sql"
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
	UseTx     bool   // indicates whether to run migration in transaction or not.
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
			return internal.ErrFailedToRunMigration(m.Name, "SQL", m.Source, err)
		}
		defer file.Close()

		statements, useTx, err := sql.ParseSQLMigration(bufio.NewReader(file), direction)
		if err != nil {
			return internal.ErrFailedToRunMigration(m.Name, "SQL", m.Source, err)
		}

		m.UseTx = useTx
		if err := m.runSQLMigration(ctx, db, statements, direction); err != nil {
			return internal.ErrFailedToRunMigration(m.Name, "SQL", m.Source, err)
		}
		break
	case internal.GoExt:
		if !m.UseTx {
			return m.runGoMigration(ctx, db.connection, db.dialect, direction)
		}

		if err := m.runGoMigration(ctx, nil, db.dialect, direction); err != nil {
			return err
		}
		break
	}

	return nil
}

func (m *Migration) runGoMigration(
	ctx context.Context,
	queryRunner shared.IQueryRunner,
	dialect internal.ISqlDialect,
	direction bool,
) error {
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
			return internal.ErrFailedToRunMigration(m.Name, "Go", fn, err)
		}
	}

	if err := m.modifyVersion(ctx, dialect, queryRunner, direction, fn, "Go"); err != nil {
		return err
	}

	m.reportSuccess(fn != nil, m.Name)

	return err
}

func (m *Migration) runSQLMigration(ctx context.Context, db *DB, statements []string, direction bool) error {
	var runner shared.IQueryRunner
	if !m.UseTx {
		runner = db.connection
	} else {
		runner = db.dialect.GetQueryRunner(ctx)
	}

	for _, statement := range statements {
		if _, err := runner.ExecContext(ctx, internal.ClearStatement(statement)); err != nil {
			return err
		}
	}

	if err := m.modifyVersion(ctx, db.dialect, runner, direction, m.Source, "SQL"); err != nil {
		return err
	}

	m.reportSuccess(len(statements) > 0, filepath.Base(m.Source))

	return nil
}

func (m *Migration) modifyVersion(
	ctx context.Context,
	dialect internal.ISqlDialect,
	queryRunner shared.IQueryRunner,
	direction bool,
	source interface{},
	migrationType string,
) error {
	if direction {
		if err := dialect.InsertVersion(ctx, queryRunner, m.Version, m.Name); err != nil {
			return internal.ErrFailedToRunMigration(m.Name, migrationType, source, err)
		}
	} else {
		if err := dialect.RemoveVersion(ctx, queryRunner, m.Version); err != nil {
			return internal.ErrFailedToRunMigration(m.Name, migrationType, source, err)
		}
	}
	return nil
}

func (m *Migration) reportSuccess(condition bool, source interface{}) {
	if condition {
		log.Printf("OK   %s \n", source)
	} else {
		log.Printf("EMPTY %s \n", source)
	}
}
