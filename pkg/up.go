package dbshaker

import (
	"context"
	"database/sql"
	"time"

	"github.com/ToggyO/dbshaker/internal"
)

// Up - migrates up to a max version.
func Up(db *DB, directory string) error {
	return UpTo(db, directory, maxVersion)
}

// UpContext migrates up to a max version with context.
func UpContext(ctx context.Context, db *DB, directory string) error {
	return UpToContext(ctx, db, directory, maxVersion)
}

// UpTo migrates up to a specific version.
func UpTo(db *DB, directory string, targetVersion int64) error {
	return UpToContext(context.Background(), db, directory, targetVersion)
}

// UpToContext migrates up to a specific version with context.
func UpToContext(ctx context.Context, db *DB, directory string, targetVersion int64) error {
	logger.Println("starting migration up process...")

	err := ensureVersionTableExists(ctx, db)
	if err != nil {
		return err
	}

	return db.dialect.Transaction(ctx,
		&internal.TxBuilderOptions{RetryCount: 10, TimeoutBetweenRetries: time.Millisecond * 5},
		func(ctx context.Context, tx *sql.Tx) error {
			_, err := EnsureDBVersionContext(ctx, db)
			if err != nil {
				return err
			}

			foundMigrations, err := scanMigrations(directory, targetVersion, true)
			if err != nil {
				return err
			}

			knownMigrations, err := db.dialect.GetMigrationsList(ctx, tx, nil)
			if err != nil {
				return err
			}

			notAppliedMigrations := lookupNotAppliedMigrations(toMigrationsList(knownMigrations), foundMigrations)

			for _, migration := range notAppliedMigrations {
				if err = migration.UpContext(ctx, db); err != nil {
					return err
				}
			}

			currentDBVersion, err := EnsureDBVersionContext(ctx, db)
			if err != nil {
				return err
			}

			logger.Println(internal.GetSuccessMigrationMessage(currentDBVersion))
			return nil
		})
}
