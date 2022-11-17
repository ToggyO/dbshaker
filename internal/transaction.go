package internal

import (
	"context"
	"database/sql"

	"github.com/ToggyO/dbshaker/shared"
)

const transactionKey TransactionKey = "t_x_transaction"

type TransactionManager struct {
	db *sql.DB
}

func (tm *TransactionManager) Transaction(
	ctx context.Context,
	options *TxBuilderOptions,
	action TransactionAction,
) error {
	if options == nil {
		options = &TxBuilderOptions{RetryCount: 1}
	}

	var tx *sql.Tx
	var err error
	for i := 0; i < options.RetryCount; i++ {
		tx, err = tm.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			return err
		}

		ctx = context.WithValue(ctx, transactionKey, tx)

		err = action(ctx, tx)
		if err == nil {
			break
		}
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			return
		}

		if err != nil {
			xerr := tx.Rollback()
			if xerr != nil {
				err = xerr
			}
		} else {
			err = tx.Commit()
		}
	}()

	return err
}

func (tm *TransactionManager) GetQueryRunner(ctx context.Context) shared.IQueryRunner {
	if txRunner, ok := ctx.Value(transactionKey).(*sql.Tx); ok {
		return txRunner
	}
	return tm.db
}
