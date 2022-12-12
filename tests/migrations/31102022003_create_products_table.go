package migrations

import (
	"github.com/ToggyO/dbshaker/pkg"
	"github.com/ToggyO/dbshaker/shared"
)

func init() {
	dbshaker.RegisterGOMigration(Up31102022003, Down31102022003, true)
}

func Up31102022003(runner shared.IQueryRunner) error {
	_, err := runner.Exec(
		`CREATE TABLE products(
		id SERIAL PRIMARY KEY,
		name VARCHAR NOT NULL,
		price DECIMAL NOT NULL
   	);`)
	if err != nil {
		return err
	}
	return nil
}

func Down31102022003(runner shared.IQueryRunner) error {
	_, err := runner.Exec("DROP TABLE products;")
	if err != nil {
		return err
	}
	return nil
}
