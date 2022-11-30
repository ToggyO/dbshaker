package main

import (
	"log"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/ToggyO/dbshaker/pkg"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   internal.CmdStatus,
	Short: "prints migrations status for provided directory",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		db, _, err := prepareMigrateCmdParams(cmd, args)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if err := dbshaker.Status(db, migrationRoot); err != nil {
			log.Fatalln(err.Error())
		}
	},
}
