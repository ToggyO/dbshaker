package dbshaker

import (
	"errors"
	"fmt"
	"github.com/ToggyO/dbshaker/internal"
	"os"
	"path/filepath"
	"strconv"
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

		const perm uint32 = 0o755
		_, err := os.Stat(absMigrationDirectoryPath)
		if errors.Is(err, os.ErrNotExist) {
			if err = os.MkdirAll(absMigrationDirectoryPath, os.FileMode(perm)); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		return CreateMigrationTemplate(args[0], absMigrationDirectoryPath, migrationType)

	case internal.CmdUp:
		if len(args) > 0 {
			strVersion := args[0]
			if len(strVersion) > 0 {
				toVersion, err := strconv.ParseInt(strVersion, 10, 64)
				if err != nil {
					return err
				}

				if err = UpTo(db, directory, toVersion); err != nil {
					return err
				}

				return nil
			}
		}
		return Up(db, absMigrationDirectoryPath)

	case internal.CmdDown:
		if len(args) > 0 {
			strVersion := args[0]
			if len(strVersion) > 0 {
				toVersion, err := strconv.ParseInt(strVersion, 10, 64)
				if err != nil {
					return err
				}

				if err = DownTo(db, directory, toVersion); err != nil {
					return err
				}

				return nil
			}
		}
		return Down(db, absMigrationDirectoryPath)

	case internal.CmdRedo:
		// TODO:

	default:
		if err := Status(db, directory); err != nil {
			return err
		}
	}

	return nil
}
