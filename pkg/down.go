package dbshaker

import (
	"context"
	"database/sql"

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

	//currentDBVersion, _, err := EnsureDBVersionContext(ctx, db)
	//if err != nil {
	//	return err
	//}
	//
	//if currentDBVersion < targetVersion {
	//	logger.Println("database is already up-to-date. current version: %d", currentDBVersion)
	//	return nil
	//}

	return db.dialect.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		currentDBVersion, _, err := EnsureDBVersionContext(ctx, db)
		if err != nil {
			return err
		}

		if currentDBVersion < targetVersion {
			logger.Println("database is already up to date. current version: %d", currentDBVersion)
			return nil
		}

		migrations, err := lookupMigrations(directory, maxVersion)
		if err != nil {
			return err
		}

		migrationsMap := make(map[int64]*Migration)
		for _, m := range migrations {
			migrationsMap[m.Version] = m
		}

		for {
			currentDBVersion, _, err = EnsureDBVersionContext(ctx, db)
			if err != nil {
				return err
			}

			if currentDBVersion == 0 {
				logger.Println(internal.GetSuccessMigrationMessage(currentDBVersion))
				return nil
			}

			currentMigration, ok := migrationsMap[currentDBVersion]
			if !ok {
				logger.Println(internal.GetSuccessMigrationMessage(currentDBVersion))
				return nil
			}

			if currentMigration.Version <= targetVersion {
				logger.Println(internal.GetSuccessMigrationMessage(currentDBVersion))
				return nil
			}

			if err = currentMigration.DownContext(ctx, db); err != nil {
				return err
			}
		}
	})
}
