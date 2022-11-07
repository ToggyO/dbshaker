package internal

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidFileName(t *testing.T) {
	t.Run("get unrecognized migration type error", func(t *testing.T) {
		_, err := IsValidFileName("122_migration_1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrRecognizedMigrationType)
	})

	t.Run("get no filename separator error", func(t *testing.T) {
		_, err := IsValidFileName("3499498230482309.sql")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNoFilenameSeparator)
	})

	t.Run("get invalid migration id error", func(t *testing.T) {
		_, err := IsValidFileName("-12324324_solo.go")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidMigrationID)
	})

	t.Run("get ParseInt error", func(t *testing.T) {
		_, err := IsValidFileName("eldzhey_69696969_gay.go")
		require.Error(t, err)
		require.ErrorIs(t, err, strconv.ErrSyntax)
	})

	t.Run("valid file name", func(t *testing.T) {
		num, err := IsValidFileName("3049434234_hello_moto.sql")
		require.NoError(t, err)

		require.Greater(t, num, int64(0))
		require.Equal(t, int64(3049434234), num)
	})
}
