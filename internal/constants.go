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

	VersionDBIndexName = "DBSH_Version_unique_clustered"

	// SQL migration statement markers
	MarkerMigrateUpStart   = "+dbshaker UpStart"
	MarkerMigrateUpEnd     = "+dbshaker UpEnd"
	MarkerMigrateDownStart = "+dbshaker DownStart"
	MarkerMigrateDownEnd   = "+dbshaker DownEnd"
	MarkerStatementBegin   = "+dbshaker StatementBegin"
	MarkerStatementEnd     = "+dbshaker StatementEnd"
	MarkerNoTransaction    = "+dbshaker NO_TRANSACTION"

	DbLockIDSalt uint = 1234567890

	// CLI command names
	CmdCreate  = "create"
	CmdMigrate = "migrate"
	CmdUp      = "up"
	CmdUpTo    = "up-"
)
