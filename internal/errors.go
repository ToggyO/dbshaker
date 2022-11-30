package internal

import (
	"errors"
	"fmt"
)

var (
	ErrRecognizedMigrationType = errors.New("[dbshaker]: not a recognized migration file type")
	ErrNoFilenameSeparator     = errors.New("[dbshaker]: no filename separator '_' found")
	ErrInvalidMigrationID      = errors.New("[dbshaker]: migration IDs must be greater than zero")
	ErrUnregisteredGoMigration = errors.New("[dbshaker]: go migration functions must be registered via `RegisterGOMigration`")
	ErrLockAcquired            = errors.New("[dbshaker]: can't acquire lock")
	ErrLockNotAcquired         = errors.New("[dbshaker]: can't unlock, as not currently locked")
	ErrLockTimeout             = errors.New("[dbshaker]: timeout: can't acquire database lock")

	ErrTryLockFailed = func(err error) error {
		return fmt.Errorf("[dbshaker]: try lock failed: %w", err)
	}

	ErrTryUnlockFailed = func(err error) error {
		return fmt.Errorf("[dbshaker]: try unlock failed: %w", err)
	}

	ErrCouldNotParseMigration = func(source string, err error) error {
		return fmt.Errorf("[dbshaker]: could not parse go migration file %q: %w", source, err)
	}

	ErrDuplicateVersion = func(version int64, source1, source2 string) error {
		return fmt.Errorf("[dbshaker]: duplicate version %v detected:\n%v\n%v", version, source1, source2)
	}

	ErrNoMigrationsInDirectory = func(dir string) error {
		return fmt.Errorf("[dbshaker]: no migrations found in provided directory: %s", dir)
	}

	ErrFailedToRunMigration = func(source string, migrationType string, migrationFunc interface{}, err error) error {
		return fmt.Errorf(
			"ERROR %v: failed to run %s migration function %T: %w", migrationType, source, migrationFunc, err)
	}

	ErrFailedToCreateMigration = func(err error) error {
		return fmt.Errorf("[dbshaker]: failed to create migration file: %w", err)
	}

	ErrMissingSQLParsingAnnotation = func(annotation string) error {
		return fmt.Errorf("failed to parse migration: missing `-- %s` annotation", annotation)
	}

	ErrUnfinishedSQLQuery = func(state int, direction bool, remaining string) error {
		return fmt.Errorf(
			"failed to parse migration: state %q, direction: %v: unexpected unfinished SQL query: %q:"+
				" missing semicolon", state, direction, remaining)
	}
)
