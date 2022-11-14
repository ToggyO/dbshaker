package shared

import (
	"context"
	"database/sql"
)

// TransactionAction function that will be executed while the transaction is running.
type TransactionAction = func(ctx context.Context, tx *sql.Tx) error

// ITransactionBuilder represent an SQL transaction process runner.
type ITransactionBuilder interface {
	Transaction(ctx context.Context, action TransactionAction) error
	GetQueryRunner(ctx context.Context) IQueryRunner
}
