package dbshaker

import (
	"github.com/ToggyO/dbshaker/internal"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestCreateMigrationTemplate(t *testing.T) {
	dir := t.TempDir()

	var err error
	migrationNames := []string{"create_users_table", "create_events_table", "create_products_table"}

	for _, name := range migrationNames {
		// tODO: добавит для sql template
		err = CreateMigrationTemplate(name, dir, GoTemplate)
		time.Sleep(1 * time.Second)
		require.NoError(t, err)
	}

	entry, err := os.ReadDir(dir)
	require.Len(t, entry, len(migrationNames))

	sort.Slice(entry, func(i, j int) bool {
		return entry[i].Name() < entry[j].Name()
	})

	for i, e := range entry {
		fileName := e.Name()

		_, err := internal.IsValidFileName(fileName)
		require.NoError(t, err)

		index := strings.Index(fileName, internal.FileNameSeparator) + 1
		migrationName := fileName[index:]
		migrationName = strings.TrimSuffix(filepath.Base(migrationName), filepath.Ext(migrationName))

		require.Equal(t, migrationNames[i], migrationName)
	}
}
