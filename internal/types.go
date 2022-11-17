package internal

import "time"

type TransactionKey string

type Dialect string

type MigrationListFilter struct {
	Offset int
	Limit  int
}

type MigrationRecord struct {
	Version     int64     `db:"version"`
	AppliedAt   time.Time `db:"applied_at"`
	Description string    `db:"description"`
}

type MigrationRecords []MigrationRecord
