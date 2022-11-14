package dbshaker

import (
	"fmt"
	"github.com/ToggyO/dbshaker/internal"
	"path/filepath"
	"runtime"
)

type folderGoMigrationRegistry map[int64]*Migration

// registry stores registered go migrations by key - path to migration folder.
var registry = make(map[string]folderGoMigrationRegistry)

// TODO: add comment
func RegisterGOMigration(up MigrationFunc, down MigrationFunc, useTx bool) {
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

	migration := &Migration{
		Name:    filename,
		Version: version,
		UpFn:    up,
		DownFn:  down,
		UseTx:   useTx,
	}

	if exists, ok := folderRegistry[version]; ok {
		logger.Fatal(fmt.Sprintf("failed to add migration %q: conflicts with exitsting %q", filename, exists.Name))
	}

	folderRegistry[version] = migration
	registry[key] = folderRegistry
}
