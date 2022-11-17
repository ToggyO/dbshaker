package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ToggyO/dbshaker/shared"
	"strings"
)

type postgresDialect struct {
	TransactionManager
	tableName string
}

func NewPostgresDialect(db *sql.DB, tableName string) ISqlDialect {
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
		VersionDBIndexName, p.tableName)
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
	filter *MigrationListFilter,
) (MigrationRecords, error) {
	if queryRunner == nil {
		queryRunner = p.GetQueryRunner(ctx)
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(`SELECT version, applied_at, description FROM %s OFFSET $1`, p.tableName))

	var params []any
	if filter == nil {
		filter = &MigrationListFilter{}
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
	migrations := make(MigrationRecords, 0)

	for rows.Next() {
		var model MigrationRecord

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

// TODO: check (mb obsolete)
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
