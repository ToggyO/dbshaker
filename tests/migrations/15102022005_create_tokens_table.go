package migrations

import (
	"github.com/ToggyO/dbshaker/pkg"
	"github.com/ToggyO/dbshaker/shared"
)

func init() {
	dbshaker.RegisterGOMigration(Up15102022005, Down15102022005, true)
}

func Up15102022005(runner shared.IQueryRunner) error {
	_, err := runner.Exec(
		`CREATE TABLE tokens(
		id SERIAL PRIMARY KEY,
		body VARCHAR NOT NULL
   	);`)
	if err != nil {
		return err
	}
	return nil
}

func Down15102022005(runner shared.IQueryRunner) error {
	_, err := runner.Exec("DROP TABLE tokens;")
	if err != nil {
		return err
	}
	return nil
}
