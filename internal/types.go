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

func (mr MigrationRecords) ToMigrationsList() []*Migration {
	migrations := make([]*Migration, 0, len(mr))

	for _, migrationRecord := range mr {
		migrations = append(migrations, &Migration{
			Version: migrationRecord.Version,
		})
	}

	return migrations
}

type DBVersion struct {
	Version int64
	Patch   byte
}
