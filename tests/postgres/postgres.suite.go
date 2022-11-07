package postgres

import (
	"github.com/ToggyO/dbshaker/pkg"
	"github.com/ToggyO/dbshaker/tests/suites"
	_ "github.com/lib/pq" //lint:ignore revive
	"github.com/stretchr/testify/require"
)

type PgTestSuite struct {
	suites.ServiceFixtureSuite
}

func (s *PgTestSuite) SetupSuite() {
	s.Init("postgres", CreatePgConnectionString(suites.NewDBConf("postgres/.env")))
}

func (s *PgTestSuite) TestMigrationDownTo() {
	err := dbshaker.DownTo(s.DB, s.MigrationRoot, 15102022005)
	require.NoError(s.Suite.T(), err)

	migrations, err := dbshaker.ListMigrations(s.DB)
	require.NoError(s.Suite.T(), err)
	require.Len(s.T(), migrations, 2)
}
