package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/ToggyO/dbshaker/pkg"
)

const (
	createCmdArgName      = "name"
	createCmdArgNameS     = "n"
	createCmdArgNameUsage = "Name for the migration"

	createCmdArgType      = "type"
	createCmdArgTypeS     = "t"
	createCmdArgTypeUsage = "Migration type (go or sql)"
)

var createCmd = &cobra.Command{
	Use:   internal.CmdCreate,
	Short: "creates migration template",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString(createCmdArgName)
		if err != nil {
			log.Fatalln(err.Error())
		}

		migrationType, err := cmd.Flags().GetString(createCmdArgType)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if err = dbshaker.Run(nil, internal.CmdCreate, migrationRoot, name, migrationType); err != nil {
			log.Fatalln(err.Error())
		}
	},
}

func init() {
	createCmd.PersistentFlags().StringP(createCmdArgName, createCmdArgNameS, "", createCmdArgNameUsage)
	createCmd.PersistentFlags().StringP(createCmdArgType, createCmdArgTypeS, string(dbshaker.GoTemplate), createCmdArgTypeUsage)
}
