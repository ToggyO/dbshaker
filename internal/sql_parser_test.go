package internal

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestParseSQLMigration(t *testing.T) {
	testCases := []struct {
		sql                 string
		upStatementsCount   int
		downStatementsCount int
		withTransaction     bool
	}{
		{
			sql:                 multilineSQL,
			upStatementsCount:   1,
			downStatementsCount: 1,
			withTransaction:     true,
		},
	}

	for _, testCase := range testCases {
		statements, _, err := ParseSQLMigration(strings.NewReader(testCase.sql), true)
		require.NoError(t, err)
		require.Len(t, statements, testCase.upStatementsCount)

		statements, _, err = ParseSQLMigration(strings.NewReader(testCase.sql), false)
		require.NoError(t, err)
		require.Len(t, statements, testCase.upStatementsCount)
	}
}

const multilineSQL = `-- +dbshaker Up
	CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		age INTEGER NOT NULL,
		name VARCHAR NOT NULL
	);

-- +dbshaker Down
	DROP TABLE users;
`
