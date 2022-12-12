package main

import (
	"log"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/ToggyO/dbshaker/pkg"
	"github.com/spf13/cobra"
)

var migrateUpCmd = &cobra.Command{
	Use:   internal.CmdUp,
	Short: "run migrate up",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		db, strVersion, err := prepareMigrateCmdParams(cmd, args)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if err = dbshaker.Run(db, internal.CmdUp, migrationRoot, strVersion); err != nil {
			log.Fatalln(err.Error())
		}
	},
}
