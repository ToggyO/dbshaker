package main

import (
	"log"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/ToggyO/dbshaker/pkg"
	"github.com/spf13/cobra"
)

var migrateDownCmd = &cobra.Command{
	Use:   internal.CmdDown,
	Short: "run migrate down",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		db, strVersion, err := prepareMigrateCmdParams(cmd, args)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if err = dbshaker.Run(db, internal.CmdDown, migrationRoot, strVersion); err != nil {
			log.Fatalln(err.Error())
		}
	},
}
