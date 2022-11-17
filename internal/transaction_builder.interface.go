package internal

import (
	"context"
	"database/sql"
	"time"

	"github.com/ToggyO/dbshaker/shared"
)

// TransactionAction function that will be executed while the transaction is running.
type TransactionAction = func(ctx context.Context, tx *sql.Tx) error

// ITransactionBuilder represent an SQL transaction process runner.
type ITransactionBuilder interface {
	Transaction(ctx context.Context, options *TxBuilderOptions, action TransactionAction) error
	GetQueryRunner(ctx context.Context) shared.IQueryRunner
}

type TxBuilderOptions struct {
	RetryCount            int
	TimeoutBetweenRetries time.Duration
}
