package dbshaker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupNotAppliedMigrations(t *testing.T) {
	known := Migrations{
		{Version: 1},
		{Version: 3},
		{Version: 4},
		{Version: 6},
		{Version: 7},
	}

	found := Migrations{
		{Version: 1},
		{Version: 2}, // not applied and below current db version
		{Version: 3},
		{Version: 4},
		{Version: 5}, // not applied and below current db version
		{Version: 6},
		{Version: 7}, // curren db version
		{Version: 8}, // not applied and above current db version
	}

	migrations := lookupNotAppliedMigrations(known, found)
	require.Len(t, migrations, 3)

	require.Equal(t, migrations[0].Version, int64(2))
	require.Equal(t, migrations[1].Version, int64(5))
	require.Equal(t, migrations[2].Version, int64(8))
}
