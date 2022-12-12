package dbshaker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ToggyO/dbshaker/internal"
)

var logger = internal.NewStdLogger()

func Run(db *DB, command, directory string, args ...string) error {
	absMigrationDirectoryPath, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	switch command {
	case internal.CmdCreate:
		return runCreateCmd(absMigrationDirectoryPath, args)

	case internal.CmdUp:
		return runUpCmd(db, absMigrationDirectoryPath, args)

	case internal.CmdDown:
		return runDownCmd(db, absMigrationDirectoryPath, args)

	case internal.CmdRedo:
		// TODO:

	default:
		if err := Status(db, directory); err != nil {
			return err
		}
	}

	return nil
}

func runCreateCmd(dirPath string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf(
			"create must be of form: dbshaker create create --dir|-d <path> --name|-n <name> --type|-t [go|sql]")
	}

	migrationType := GoTemplate
	if len(args) == 2 {
		migrationType = MigrationTemplateType(args[1])
	}

	const perm uint32 = 0o755
	_, err := os.Stat(dirPath)
	if errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(dirPath, os.FileMode(perm)); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return CreateMigrationTemplate(args[0], dirPath, migrationType)
}

func runUpCmd(db *DB, dirPath string, args []string) error {
	if len(args) > 0 && len(args[0]) > 0 {
		strVersion := args[0]
		toVersion, err := strconv.ParseInt(strVersion, 10, 64)
		if err != nil {
			return err
		}

		if err = UpTo(db, dirPath, toVersion); err != nil {
			return err
		}

		return nil
	}
	return Up(db, dirPath)
}

func runDownCmd(db *DB, dirPath string, args []string) error {
	if len(args) > 0 && len(args[0]) > 0 {
		strVersion := args[0]
		toVersion, err := strconv.ParseInt(strVersion, 10, 64)
		if err != nil {
			return err
		}

		if err = DownTo(db, dirPath, toVersion); err != nil {
			return err
		}

		return nil
	}
	return Down(db, dirPath)
}
