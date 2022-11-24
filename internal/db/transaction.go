package db

import (
	"context"
	"database/sql"
	"github.com/ToggyO/dbshaker/internal"
	"time"

	"github.com/ToggyO/dbshaker/shared"
)

const (
	transactionKey internal.TransactionKey = "t_x_transaction"
)

type TransactionManager struct {
	db *sql.DB
}

func (tm *TransactionManager) Transaction(ctx context.Context, action internal.TransactionAction) error {
	return tm.TransactionConfigurable(ctx, nil, action)
}

// TODO: доработать логику retry.
func (tm *TransactionManager) TransactionConfigurable(
	ctx context.Context,
	options *internal.TxBuilderOptions,
	action internal.TransactionAction,
) error {
	if options == nil {
		options = &internal.TxBuilderOptions{RetryCount: 1, TimeoutBetweenRetries: time.Millisecond}
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
		// TODO: добавить проверку на занятость транзакцией
		if err == nil {
			break
		}
		tm.sleep(ctx, options.TimeoutBetweenRetries)
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

func (tm *TransactionManager) sleep(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}
