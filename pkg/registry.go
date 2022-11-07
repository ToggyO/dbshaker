package dbshaker

import "github.com/ToggyO/dbshaker/internal"

type folderGoMigrationRegistry map[int64]*internal.Migration

// registry stores registered go migrations by key - path to migration folder.
var registry = make(map[string]folderGoMigrationRegistry)
