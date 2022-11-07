package migrations

import (
	"database/sql"

	"github.com/ToggyO/dbshaker/pkg"
)

func init() {
	dbshaker.AddMigration(Up15102022005, Down15102022005)
}

func Up15102022005(tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE tokens(
		id SERIAL PRIMARY KEY,
		body VARCHAR NOT NULL
   	);`)
	if err != nil {
		return err
	}
	return nil
}

func Down15102022005(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE tokens;")
	if err != nil {
		return err
	}
	return nil
}
