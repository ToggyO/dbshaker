package dbshaker

import (
	"github.com/ToggyO/dbshaker/internal"
	dbshaker "github.com/ToggyO/dbshaker/pkg"
	"github.com/spf13/cobra"
)

const (
	toVersionCmdArgName  = "to"
	toVersionCmdArgUsage = "Migrate database to specific version"
)

var migrateCmd = &cobra.Command{
	Use:              internal.CmdMigrate,
	Short:            "run migrations",
	TraverseChildren: true,
}

func init() {
	migrateCmd.PersistentFlags().String(toVersionCmdArgName, "", toVersionCmdArgUsage)
	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd, redoCmd)
}

func prepareMigrateCmdParams(cmd *cobra.Command, args []string) (*dbshaker.DB, string, error) {
	driver := args[0]
	connString := args[1]
	db, err := dbshaker.OpenDBWithDriver(driver, connString)
	if err != nil {
		return nil, "", err
	}

	strVersion, _ := cmd.Flags().GetString(toVersionCmdArgName)
	return db, strVersion, nil

}
