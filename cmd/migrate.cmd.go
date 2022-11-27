package main

import (
	dbshaker "github.com/ToggyO/dbshaker/pkg"
	"github.com/spf13/cobra"
	"log"

	"github.com/ToggyO/dbshaker/internal"
)

var migrateCmd = &cobra.Command{
	Use:   internal.CmdMigrate,
	Short: "run migrations",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		driver := args[0]
		connString := args[1]
		dir, err := cmd.PersistentFlags().GetString(directoryCmdArg)
		if err != nil {
			log.Fatalln(err.Error())
		}

		db, err := dbshaker.OpenDBWithDriver(driver, connString)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if err = dbshaker.Up(db, dir); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {

	// TODO:
	//migrateCmd.Flags().
}
