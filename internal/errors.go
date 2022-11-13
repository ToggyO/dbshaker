package internal

import (
	"errors"
	"fmt"
)

var (
	ErrRecognizedMigrationType = errors.New("[dbshaker]: not a recognized migration file type")
	ErrNoFilenameSeparator     = errors.New("[dbshaker]: no filename separator '_' found")
	ErrInvalidMigrationID      = errors.New("[dbshaker]: migration IDs must be greater than zero")
	ErrUnregisteredGoMigration = errors.New("[dbshaker]: go migration functions must be registered via `AddMigration`")

	ErrCouldNotParseMigration = func(source string, err error) error {
		return fmt.Errorf("[dbshaker]: could not parse go migration file %q: %w", source, err)
	}

	ErrDuplicateVersion = func(version int64, source1, source2 string) error {
		return fmt.Errorf("[dbshaker]: duplicate version %v detected:\n%v\n%v", version, source1, source2)
	}

	ErrDBAlreadyIsUpToDate = func(version int64) error {
		return fmt.Errorf("[dbshaker]: database is already up to date. current version: %d", version)
	}

	ErrFailedToRunMigration = func(source string, migrationFunc MigrationFunc, err error) error {
		return fmt.Errorf("ERROR %v: failed to run Go migration function %T: %w", source, migrationFunc, err)
	}

	ErrFailedToCreateMigration = func(err error) error {
		return fmt.Errorf("[dbshaker]: failed to create migration file: %w", err)
	}

	errMissingSQLParsingAnnotation = func(annotation string) error {
		return fmt.Errorf("failed to parse migration: missing `-- %s` annotation", annotation)
	}

	errUnfinishedSQLQuery = func(state int, direction bool, remaining string) error {
		return fmt.Errorf(
			"failed to parse migration: state %q, direction: %v: unexpected unfinished SQL query: %q:"+
				" missing semicolon", state, direction, remaining)
	}
)
