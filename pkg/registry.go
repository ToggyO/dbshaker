package dbshaker

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ToggyO/dbshaker/internal"
)

// registry stores registered go migrations.
var registry = make(map[int64]*Migration)

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

	sourceDir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		panic(err)
	}

	migration := &Migration{
		Name:      filepath.Base(filename),
		Version:   version,
		UpFn:      up,
		DownFn:    down,
		UseTx:     useTx,
		Source:    filename,
		SourceDir: sourceDir,
	}

	if exists, ok := registry[version]; ok {
		logger.Fatal(fmt.Sprintf("failed to add migration %q: conflicts with exitsting %q", filename, exists.Name))
	}

	registry[version] = migration
}
