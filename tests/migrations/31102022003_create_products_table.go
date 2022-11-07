package migrations

import (
	"database/sql"

	"github.com/ToggyO/dbshaker/pkg"
)

func init() {
	dbshaker.AddMigration(Up31102022003, Down31102022003)
}

func Up31102022003(tx *sql.Tx) error {
	_, err := tx.Exec(
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

func Down31102022003(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE products;")
	if err != nil {
		return err
	}
	return nil
}
