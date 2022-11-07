package migrations

import (
	"database/sql"

	"github.com/ToggyO/dbshaker/pkg"
)

func init() {
	dbshaker.AddMigration(Up11092022001, Down11092022001)
}

func Up11092022001(tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE users(
		id SERIAL PRIMARY KEY,
		name VARCHAR NOT NULL
   	);`)
	if err != nil {
		return err
	}
	return nil
}

func Down11092022001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE users;")
	if err != nil {
		return err
	}
	return nil
}
