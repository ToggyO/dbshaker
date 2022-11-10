package internal

import "text/template"

var GoMigrationTemplate = template.Must(template.New("dbshaker.go-migration").Parse(`package migrations

import (
    "database/sql"

    "github.com/ToggyO/dbshaker/pkg"
)

func init() {
	dbshaker.AddMigration(up{{.MName}}, down{{.MName}})
}

func up{{.MName}}(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func down{{.MName}}(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
`))
