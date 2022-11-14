package internal

import (
	"context"
	"github.com/ToggyO/dbshaker/shared"
)

type ISqlDialect interface {
	shared.ITransactionBuilder

	CreateVersionTable(ctx context.Context) error
	InsertVersion(ctx context.Context, version int64) error
	IncrementVersionPatch(ctx context.Context, version int64) error
	RemoveVersion(ctx context.Context, version int64) error
	GetMigrationsList(ctx context.Context, filter *MigrationListFilter) (MigrationRecords, error)
	GetDBVersion(ctx context.Context) (DBVersion, error)
}
