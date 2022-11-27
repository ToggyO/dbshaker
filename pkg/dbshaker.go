package dbshaker

import (
	"errors"
	"fmt"
	"github.com/ToggyO/dbshaker/internal"
	"os"
	"path/filepath"
)

var logger = internal.NewStdLogger()

func Run(db *DB, command, directory string, args ...string) error {
	absMigrationDirectoryPath, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	switch command {
	case internal.CmdCreate:
		if len(args) == 0 {
			return fmt.Errorf("create must be of form: dbshaker create create --dir|-d <path> --name|-n <name> --type|-t [go|sql]")
		}

		migrationType := GoTemplate
		if len(args) == 2 {
			migrationType = MigrationTemplateType(args[1])
		}

		// TODO: check if statement
		const perm uint32 = 0o755
		_, err := os.Stat(absMigrationDirectoryPath)
		if errors.Is(err, os.ErrNotExist) {
			if err = os.Mkdir(absMigrationDirectoryPath, os.FileMode(perm)); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		return CreateMigrationTemplate(args[0], absMigrationDirectoryPath, migrationType)
	}
	return nil
}
