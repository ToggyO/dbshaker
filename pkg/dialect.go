package dbshaker

import (
	"database/sql"
	"fmt"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/ToggyO/dbshaker/internal/db"
)

func createDialect(connection *sql.DB, d string) (internal.ISqlDialect, error) {
	// TODO: добавить поддержку диалектов других СУБД
	switch d {
	case internal.PostgresDialect, internal.PgxDialect:
		return db.NewPostgresDialect(connection, internal.ServiceTableName), nil
	default:
		return nil, fmt.Errorf("%q: unknown dialect", d)
	}
}
