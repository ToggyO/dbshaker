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
	markerMigrateUpStart   = "+dbshaker UpStart"
	markerMigrateUpEnd     = "+dbshaker UpEnd"
	markerMigrateDownStart = "+dbshaker DownStart"
	markerMigrateDownEnd   = "+dbshaker DownEnd"
	markerStatementBegin   = "+dbshaker StatementBegin"
	markerStatementEnd     = "+dbshaker StatementEnd"
	markerNoTransaction    = "+dbshaker NO_TRANSACTION"
)
