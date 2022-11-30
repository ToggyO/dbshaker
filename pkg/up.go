package dbshaker

import (
	"context"

	"github.com/ToggyO/dbshaker/internal"
)

// Up - migrates up to a max version.
func Up(db *DB, directory string) error {
	return UpTo(db, directory, internal.MaxVersion)
}

// UpContext migrates up to a max version with context.
func UpContext(ctx context.Context, db *DB, directory string) error {
	return UpToContext(ctx, db, directory, internal.MaxVersion)
}

// UpTo migrates up to a specific version.
func UpTo(db *DB, directory string, targetVersion int64) error {
	return UpToContext(context.Background(), db, directory, targetVersion)
}

// UpToContext migrates up to a specific version with context.
func UpToContext(ctx context.Context, db *DB, directory string, targetVersion int64) error {
	logger.Println("starting migration up process...")
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := lockDb(ctx, db); err != nil {
		return err
	}

	_, err := EnsureDBVersionContext(ctx, db)
	if err != nil {
		return err
	}

	foundMigrations, err := scanMigrations(directory, targetVersion, true)
	if err != nil {
		return err
	}

	knownMigrations, err := db.dialect.GetMigrationsList(ctx, nil, nil)
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

	if err := db.dialect.Unlock(ctx); err != nil {
		return err
	}

	logger.Println(internal.GetSuccessMigrationMessage(currentDBVersion))
	return nil
}
