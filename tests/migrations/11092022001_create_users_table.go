package migrations

import (
	"github.com/ToggyO/dbshaker/pkg"
	"github.com/ToggyO/dbshaker/shared"
)

func init() {
	dbshaker.RegisterGOMigration(Up11092022001, Down11092022001, true)
}

func Up11092022001(runner shared.IQueryRunner) error {
	_, err := runner.Exec(
		`CREATE TABLE users(
		id SERIAL PRIMARY KEY,
		name VARCHAR NOT NULL
   	);`)
	if err != nil {
		return err
	}
	return nil
}

func Down11092022001(runner shared.IQueryRunner) error {
	_, err := runner.Exec("DROP TABLE users;")
	if err != nil {
		return err
	}
	return nil
}
