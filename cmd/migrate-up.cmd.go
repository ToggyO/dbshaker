package main

import (
	"github.com/spf13/cobra"

	"github.com/ToggyO/dbshaker/internal"
)

const (
	toVersionCmdArgName = "to"
)

var migrateUpCmd = &cobra.Command{
	Use:   internal.CmdUp,
	Short: "run migrate up",
	Run: func(cmd *cobra.Command, args []string) {
		//db, err := dbshaker.OpenDBWithDriver()
		//err := dbshaker.Up()
	},
}

func init() {
	//migrateUpCmd.PersistentFlags()
	//migrateUpCmd.Flags().Int64(toVersionCmdArgName, , "")
}
