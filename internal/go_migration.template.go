package internal

import "text/template"

var GoMigrationTemplate = template.Must(template.New("dbshaker.go-migration").Parse(`package migrations

import (
    "github.com/ToggyO/dbshaker/pkg"
    "github.com/ToggyO/dbshaker/shared"
)

func init() {
	dbshaker.RegisterGOMigration(up{{.MName}}, down{{.MName}}, true)
}

func up{{.MName}}(runner shared.IQueryRunner) error {
	// This code is executed when the migration is applied.
	return nil
}

func down{{.MName}}(runner shared.IQueryRunner) error {
	// This code is executed when the migration is rolled back.
	return nil
}
`))
