package dbshaker

// TODO: add comment
// runSQLMigration allows to run a set of SQL statements
// TODO: add db versioning
//func runSQLMigration(ctx context.Context, db *DB, statements []string, useTx bool, version int64, direction bool) error {
//	if useTx {
//		return db.dialect.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
//			for _, statement := range statements {
//				if _, err := tx.ExecContext(ctx, internal.ClearStatement(statement)); err != nil {
//					return err
//				}
//			}
//			return nil
//		})
//	}
//
//	for _, statement := range statements {
//		if _, err := db.connection.ExecContext(ctx, internal.ClearStatement(statement)); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
