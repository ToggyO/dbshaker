package internal

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func IsValidFileName(value string) (int64, error) {
	base := filepath.Base(value)
	if ext := filepath.Ext(base); ext != GoExt && ext != SQLExt {
		return 0, ErrRecognizedMigrationType
	}

	index := strings.Index(base, FileNameSeparator)
	if index < 0 {
		return 0, ErrNoFilenameSeparator
	}

	num, err := strconv.ParseInt(base[:index], 10, 64)
	if err == nil && num <= 0 {
		return 0, ErrInvalidMigrationID
	}

	return num, err
}

func GetSuccessMigrationMessage(currentDBVersion int64) string {
	return fmt.Sprintf("no migrations to run. current version: %d\n", currentDBVersion)
}

var (
	matchSQLComments = regexp.MustCompile(`(?m)^--.*$[\r\n]*`)
	matchEmptyEOL    = regexp.MustCompile(`(?m)^$[\r\n]*`) // TODO: Duplicate
)

func ClearStatement(s string) string {
	s = matchSQLComments.ReplaceAllString(s, ``)
	return matchEmptyEOL.ReplaceAllString(s, ``)
}
