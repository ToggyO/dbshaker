package internal

const (
	ToolName         = "dbshaker"
	ServiceTableName = "dbshaker_version"

	GoExt             = ".go"
	SQLExt            = ".sql"
	GoFilesPattern    = "*.go"
	SQLFilesPattern   = "*.sql"
	FileNameSeparator = "_"
	SQLCommentPrefix  = "--"
	SQLSemicolon      = ";"

	PostgresDialect = "postgres"
	PgxDialect      = "pgx"

	// SQL migration statement markers
	markerMigrateUpStart   = "+dbshaker Up start"
	markerMigrateUpEnd     = "+dbshaker Up end"
	markerMigrateDownStart = "+dbshaker Down start"
	markerMigrateDownEnd   = "+dbshaker Down end"
	markerStatementBegin   = "+dbshaker StatementBegin"
	markerStatementEnd     = "+dbshaker StatementEnd"
	markerNoTransaction    = "+dbshaker NO_TRANSACTION"
)
