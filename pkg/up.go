package dbshaker

import (
	"context"
	"database/sql"

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

	currentDBVersion, _, err := EnsureDBVersionContext(ctx, db)
	if err != nil {
		return err
	}

	if currentDBVersion > targetVersion {
		logger.Println("database is already up to date. current version: %d", currentDBVersion)
		return nil
	}

	return db.dialect.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		foundMigrations, err := lookupMigrations(directory, targetVersion)
		if err != nil {
			return err
		}

		dbMigrations, err := db.dialect.GetMigrationsList(ctx, nil)
		if err != nil {
			return err
		}

		notAppliedMigrations := lookupNotAppliedMigrations(dbMigrations.ToMigrationsList(), foundMigrations)

		for _, migration := range notAppliedMigrations {
			if err = migration.UpContext(ctx, tx, db.dialect); err != nil {
				return err
			}
		}

		notAppliedMigrationsLen := len(notAppliedMigrations)
		if notAppliedMigrationsLen > 0 {
			if notAppliedMigrations[notAppliedMigrationsLen-1].Version < currentDBVersion {
				err = db.dialect.IncrementVersionPatch(ctx, currentDBVersion)
				if err != nil {
					return err
				}
			}
		}

		currentDBVersion, _, err = EnsureDBVersionContext(ctx, db)
		if err != nil {
			return err
		}

		logger.Println(internal.GetSuccessMigrationMessage(currentDBVersion))
		return nil
	})
}
