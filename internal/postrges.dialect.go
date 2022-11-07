package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (p *postgresDialect) CreateVersionTable(ctx context.Context) error {
	query := fmt.Sprintf(`CREATE TABLE %s (
			version BIGINT NOT NULL UNIQUE,
			patch INTEGER DEFAULT 0,
			applied_at TIMESTAMP DEFAULT NOW()
	);`, p.tableName)

	_, err := p.GetQueryRunner(ctx).ExecContext(ctx, query)
	return err
}

func (p *postgresDialect) InsertVersion(ctx context.Context, version int64) error {
	query := fmt.Sprintf(`INSERT INTO %s (version) VALUES ($1);`, p.tableName)
	_, err := p.GetQueryRunner(ctx).ExecContext(ctx, query, version)
	return err
}

func (p *postgresDialect) IncrementVersionPatch(ctx context.Context, version int64) error {
	query := fmt.Sprintf(`UPDATE %s SET patch = patch + 1 WHERE version = $1`, p.tableName)
	_, err := p.GetQueryRunner(ctx).ExecContext(ctx, query, version)
	return err
}

func (p *postgresDialect) RemoveVersion(ctx context.Context, version int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE version = $1;`, p.tableName)
	_, err := p.GetQueryRunner(ctx).ExecContext(ctx, query, version)
	return err
}

func (p *postgresDialect) GetMigrationsList(
	ctx context.Context,
	filter *MigrationListFilter,
) (MigrationRecords, error) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(`SELECT version, patch, applied_at FROM %s OFFSET $1`, p.tableName))

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
	rows, err := p.GetQueryRunner(ctx).QueryContext(ctx, k, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}

	defer rows.Close()
	migrations := make(MigrationRecords, 0)

	for rows.Next() {
		var model MigrationRecord

		if err := rows.Scan(&model.Version, &model.Patch, &model.AppliedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		migrations = append(migrations, model)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get next row: %w", err)
	}

	return migrations, nil
}

func (p *postgresDialect) GetDBVersion(ctx context.Context) (DBVersion, error) {
	query := fmt.Sprintf(`SELECT version, patch FROM %s ORDER BY version DESC;`, p.tableName)
	queryRunner := p.GetQueryRunner(ctx)

	var version DBVersion

	if err := queryRunner.QueryRowContext(ctx, query).Scan(&version.Version, &version.Patch); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return version, nil
		}
		return version, err
	}

	return version, nil
}
