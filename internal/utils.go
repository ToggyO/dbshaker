package internal

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func IsValidFileName(value string) (int64, error) {
	base := filepath.Base(value)
	if ext := filepath.Ext(base); ext != GoExt && ext != SqlExt {
		return 0, ErrRecognizedMigrationType
	}

	index := strings.Index(base, FileNameSeparator)
	if index < 0 {
		return 0, ErrNoFilenameSeparator
	}

	num, err := strconv.ParseInt(base[:index], 10, 64)
	if err == nil && num <= 0 {
		return 0, ErrInvalidMigrationId
	}

	return num, err
}

func GetSuccessMigrationMessage(currentDbVersion int64) string {
	return fmt.Sprintf("no migrations to run. current version: %d\n", currentDbVersion)
}
