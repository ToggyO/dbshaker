package dbshaker

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ToggyO/dbshaker/internal"
)

func AddMigration(up internal.MigrationFunc, down internal.MigrationFunc) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("ERROR: error during parsing file name")
	}

	version, err := internal.IsValidFileName(filename)
	if err != nil {
		panic(err)
	}

	key := filepath.Dir(filename)
	if err != nil {
		panic(err)
	}

	folderRegistry, ok := registry[key]
	if !ok {
		folderRegistry = make(folderGoMigrationRegistry)
	}

	migration := &internal.Migration{
		Name:    filename,
		Version: version,
		UpFn:    up,
		DownFn:  down,
	}

	if exists, ok := folderRegistry[version]; ok {
		logger.Fatal(fmt.Sprintf("failed to add migration %q: conflicts with exitsting %q", filename, exists.Name))
	}

	folderRegistry[version] = migration
	registry[key] = folderRegistry
}
