package internal

import "time"

type TransactionKey string

type Dialect string

type MigrationListFilter struct {
	Offset int
	Limit  int
}

type MigrationRecord struct {
	Version   int64     `db:"version"`
	Patch     byte      `db:"patch"`
	AppliedAt time.Time `db:"applied_at"`
}

type MigrationRecords []MigrationRecord

type DBVersion struct {
	Version int64
	Patch   byte
}
