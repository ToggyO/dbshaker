package dbshaker

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/stretchr/testify/require"
)

func TestCreateMigrationTemplate(t *testing.T) {
	dir := t.TempDir()

	var err error
	migrations := []struct {
		name         string
		templateType MigrationTemplateType
	}{
		{
			name:         "create_users_table",
			templateType: GoTemplate,
		},
		{
			name:         "create_events_table",
			templateType: SQLTemplate,
		},
		{
			name:         "create_products_table",
			templateType: GoTemplate,
		},
	}

	for _, m := range migrations {
		err = CreateMigrationTemplate(m.name, dir, m.templateType)
		time.Sleep(1 * time.Second)
		require.NoError(t, err)
	}

	entry, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entry, len(migrations))

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

		require.Equal(t, migrations[i].name, migrationName)

		expectedTemplateType := MigrationTemplateType(filepath.Ext(fileName)[1:])
		require.Equal(t, migrations[i].templateType, expectedTemplateType)
	}
}
