package dbshaker

import (
	"context"
	"path/filepath"
	"sort"

	"github.com/ToggyO/dbshaker/internal"
)

const (
	maxUint    = ^uint64(0)
	maxVersion = int64(maxUint >> 1) // max(int64)

)

// Migrations runtime slice of Migration struct pointers.
type Migrations []*Migration

func (ms Migrations) Len() int {
	return len(ms)
}

func (ms Migrations) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms Migrations) Less(i, j int) bool {
	if ms[i].Version == ms[j].Version {
		logger.Fatal(internal.ErrDuplicateVersion(ms[i].Version, ms[i].Source, ms[j].Source))
	}
	return ms[i].Version < ms[j].Version
}

// ListMigrations lists all applied migrations in database.
func ListMigrations(db *DB) (Migrations, error) {
	return ListMigrationsContext(context.Background(), db)
}

// ListMigrationsContext lists all applied migrations in database with context.
func ListMigrationsContext(ctx context.Context, db *DB) (Migrations, error) {
	records, err := db.dialect.GetMigrationsList(ctx, db.dialect.GetQueryRunner(ctx), nil)
	if err != nil {
		return nil, err
	}
	return toMigrationsList(records), nil
}

// scanMigrations returns a slice of valid migrations in the migrations folder and migration registry,
// acceptable for current direction and target version.
// Returned slice of migrations is sorted by version in ascending direction.
// TODO: `embed` support in future by embed.FS.
func scanMigrations(directory string, targetVersion int64, direction bool) (Migrations, error) {
	// TODO: convert directory to absolute path
	sqlMigrationFiles, err := filepath.Glob(filepath.Join(directory, internal.SQLFilesPattern))
	if err != nil {
		return nil, err
	}

	migrations := make(Migrations, 0, len(sqlMigrationFiles)+len(registry))

	for _, file := range sqlMigrationFiles {
		v, err := internal.IsValidFileName(file)
		if err != nil {
			return nil, internal.ErrCouldNotParseMigration(file, err)
		}

		if !checkVersion(v, targetVersion, direction) {
			continue
		}

		migrations = append(migrations, &Migration{
			Name:      filepath.Base(file),
			Version:   v,
			Source:    file,
			SourceDir: filepath.Dir(file),
		})
	}

	migrationRootDir, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}

	// Migrations in `.go` files, registered via RegisterGOMigration
	for _, migration := range registry {
		if migration.SourceDir != migrationRootDir {
			continue
		}

		if !checkVersion(migration.Version, targetVersion, direction) {
			continue
		}
		migrations = append(migrations, migration)
	}

	// Unregistered `.go` migrations
	gGoMigrationsFiles, err := filepath.Glob(filepath.Join(directory, internal.GoFilesPattern))
	if err != nil {
		return nil, err
	}

	for _, file := range gGoMigrationsFiles {
		v, err := internal.IsValidFileName(file)
		if err != nil {
			continue // Пропускаем файлы, которые не имею версионного префикса
		}

		if _, ok := registry[v]; !ok {
			return nil, internal.ErrUnregisteredGoMigration
		}
	}

	sort.Sort(migrations)

	return migrations, nil
}

func lookupNotAppliedMigrations(known, found Migrations) Migrations {
	return filterMigrationsByDirection(known, found, true)
}

func lookupAppliedMigrations(known, found Migrations) Migrations {
	return filterMigrationsByDirection(known, found, false)
}

func filterMigrationsByDirection(known, found Migrations, direction bool) Migrations {
	existing := make(map[int64]bool)
	for _, k := range known {
		existing[k.Version] = true
	}

	var migrations Migrations
	for _, f := range found {
		_, ok := existing[f.Version]
		if direction && !ok {
			migrations = append(migrations, f)
		} else if !direction && ok {
			migrations = append(migrations, f)
		}
	}

	// Reverse migration for down direction to apply them in reversed order and avoid conflicts
	if !direction {
		sort.Slice(migrations, func(i, j int) bool {
			return migrations[i].Version > migrations[j].Version
		})
	}

	return migrations
}

func toMigrationsList(mr internal.MigrationRecords) []*Migration {
	migrations := make([]*Migration, 0, len(mr))

	for _, migrationRecord := range mr {
		migrations = append(migrations, &Migration{
			Version: migrationRecord.Version,
		})
	}

	return migrations
}

func checkVersion(version, targetVersion int64, direction bool) bool {
	if direction {
		return version < targetVersion
	}
	return version > targetVersion
}
