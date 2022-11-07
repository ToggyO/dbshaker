package dbshaker

import (
	"database/sql"
	"fmt"

	"github.com/ToggyO/dbshaker/internal"
)

func createDialect(connection *sql.DB, d string) (internal.ISqlDialect, error) {
	// TODO: добавить поддержку диалектов других СУБД
	switch d {
	case internal.PostgresDialect, internal.PgxDialect:
		return internal.NewPostgresDialect(connection, internal.ServiceTableName), nil
	default:
		return nil, fmt.Errorf("%q: unknown dialect", d)
	}
}
