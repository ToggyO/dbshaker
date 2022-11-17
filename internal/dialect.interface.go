package internal

import (
	"context"
	"github.com/ToggyO/dbshaker/shared"
)

type ISqlDialect interface {
	ITransactionBuilder

	CreateVersionTable(ctx context.Context, queryRunner shared.IQueryRunner) error
	InsertVersion(ctx context.Context, queryRunner shared.IQueryRunner, version int64, description string) error
	RemoveVersion(ctx context.Context, queryRunner shared.IQueryRunner, version int64) error
	GetMigrationsList(ctx context.Context, queryRunner shared.IQueryRunner, filter *MigrationListFilter) (MigrationRecords, error)
	GetDBVersion(ctx context.Context, queryRunner shared.IQueryRunner) (int64, error)
}
