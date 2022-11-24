package dbshaker

import (
	"context"
	"github.com/ToggyO/dbshaker/internal"
)

// Down rolls back all existing migrations.
func Down(db *DB, directory string) error {
	return DownContext(context.Background(), db, directory)
}

// DownContext rolls back all existing migrations with context.
func DownContext(ctx context.Context, db *DB, directory string) error {
	return DownToContext(ctx, db, directory, 0)
}

// DownTo rolls back migrations to a specific version.
func DownTo(db *DB, directory string, targetVersion int64) error {
	return DownToContext(context.Background(), db, directory, targetVersion)
}

// DownToContext rolls back migrations to a specific version with context.
func DownToContext(ctx context.Context, db *DB, directory string, targetVersion int64) error {
	logger.Printf("starting migration down process...")
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := lockDb(ctx, db); err != nil {
		return err
	}

	_, err := EnsureDBVersionContext(ctx, db)
	if err != nil {
		return err
	}

	migrations, err := scanMigrations(directory, targetVersion, false)
	if err != nil {
		return err
	}

	knownMigrations, err := db.dialect.GetMigrationsList(ctx, nil, nil)
	if err != nil {
		return err
	}

	appliedMigrations := lookupAppliedMigrations(toMigrationsList(knownMigrations), migrations)

	for _, applied := range appliedMigrations {
		if err = applied.DownContext(ctx, db); err != nil {
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
