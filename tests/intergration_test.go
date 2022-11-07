//go:build integration
// +build integration

package tests

import (
	"testing"

	_ "github.com/ToggyO/dbshaker/tests/migrations"
	"github.com/ToggyO/dbshaker/tests/postgres"
	"github.com/stretchr/testify/suite"
)

func TestIntegration(t *testing.T) {
	suite.Run(t, new(postgres.PgTestSuite))
}
