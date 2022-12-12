package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/ToggyO/dbshaker/shared"
)

type postgresDialect struct {
	TransactionManager

	tableName string
	isLocked  atomic.Bool
}

func NewPostgresDialect(db *sql.DB, tableName string) internal.ISqlDialect {
	return &postgresDialect{
		TransactionManager: TransactionManager{db: db},
		tableName:          tableName,
	}
}

func (p *postgresDialect) CreateVersionTable(ctx context.Context, queryRunner shared.IQueryRunner) error {
	if queryRunner == nil {
		queryRunner = p.GetQueryRunner(ctx)
	}

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS  %s (
			version BIGINT NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT NOW(),
    		description VARCHAR(300)
	);`, p.tableName)
	_, err := queryRunner.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s USING btree ("version");`,
		internal.VersionDBIndexName, p.tableName)
	_, err = queryRunner.ExecContext(ctx, query)

	return err
}

func (p *postgresDialect) InsertVersion(
	ctx context.Context,
	queryRunner shared.IQueryRunner,
	version int64,
	description string,
) error {
	if queryRunner == nil {
		queryRunner = p.GetQueryRunner(ctx)
	}
	query := fmt.Sprintf(`INSERT INTO %s (version, description) VALUES ($1,$2);`, p.tableName)
	_, err := queryRunner.ExecContext(ctx, query, version, description)
	return err
}

func (p *postgresDialect) RemoveVersion(ctx context.Context, queryRunner shared.IQueryRunner, version int64) error {
	if queryRunner == nil {
		queryRunner = p.GetQueryRunner(ctx)
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE version = $1;`, p.tableName)
	_, err := queryRunner.ExecContext(ctx, query, version)
	return err
}

func (p *postgresDialect) GetMigrationsList(
	ctx context.Context,
	queryRunner shared.IQueryRunner,
	filter *internal.MigrationListFilter,
) (internal.MigrationRecords, error) {
	if queryRunner == nil {
		queryRunner = p.GetQueryRunner(ctx)
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(`SELECT version, applied_at, description FROM %s OFFSET $1`, p.tableName))

	var params []any
	if filter == nil {
		filter = &internal.MigrationListFilter{}
	}

	if filter.Offset >= 0 {
		params = append(params, filter.Offset)
	}

	if filter.Limit > 0 {
		sb.WriteString("LIMIT $2")
		params = append(params, filter.Limit)
	}

	sb.WriteString(";")
	k := sb.String()
	rows, err := queryRunner.QueryContext(ctx, k, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}

	defer rows.Close()
	migrations := make(internal.MigrationRecords, 0)

	for rows.Next() {
		var model internal.MigrationRecord

		if err := rows.Scan(&model.Version, &model.AppliedAt, &model.Description); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		migrations = append(migrations, model)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get next row: %w", err)
	}

	return migrations, nil
}

func (p *postgresDialect) GetDBVersion(ctx context.Context, queryRunner shared.IQueryRunner) (int64, error) {
	query := fmt.Sprintf(`SELECT version FROM %s ORDER BY version DESC;`, p.tableName)
	if queryRunner == nil {
		queryRunner = p.GetQueryRunner(ctx)
	}

	var version int64

	if err := queryRunner.QueryRowContext(ctx, query).Scan(&version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return version, nil
		}
		return version, err
	}

	return version, nil
}

func (p *postgresDialect) Lock(ctx context.Context) error {
	return CasRestoreOnError(&p.isLocked, false, true, internal.ErrLockAcquired, func() error {
		// TODO: аргументы для генерации слабоваты.
		// TODO: P.S. рефлексию убрать
		lockID := GenerateLockID(p.tableName, reflect.TypeOf(p).Name())

		query := `SELECT pg_advisory_lock($1)`
		if _, err := p.GetQueryRunner(ctx).ExecContext(ctx, query, lockID); err != nil {
			return internal.ErrTryLockFailed(err)
		}

		return nil
	})
}

func (p *postgresDialect) Unlock(ctx context.Context) error {
	return CasRestoreOnError(&p.isLocked, true, false, internal.ErrLockNotAcquired, func() error {
		// TODO: аргументы для генерации слабоваты.
		// TODO: P.S. рефлексию убрать
		lockID := GenerateLockID(p.tableName, reflect.TypeOf(p).Name())

		query := `SELECT pg_advisory_unlock($1)`
		if _, err := p.GetQueryRunner(ctx).ExecContext(ctx, query, lockID); err != nil {
			return internal.ErrTryUnlockFailed(err)
		}

		return nil
	})
}
