package main

import (
	"path"

	"github.com/ToggyO/dbshaker/internal"
	"github.com/spf13/cobra"
)

const (
	directoryCmdArg      = "dir"
	directoryCmdArgS     = "d"
	directoryCmdArgUsage = "Path to folder, contains migrations"
)

var (
	defaultMigrationRoot = path.Join(".", "migrations")
	migrationRoot        string
)

var rootCmd = cobra.Command{
	Use:              internal.ToolName,
	TraverseChildren: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().
		StringVarP(&migrationRoot, directoryCmdArg, directoryCmdArgS, defaultMigrationRoot, directoryCmdArgUsage)

	rootCmd.AddCommand(createCmd, migrateCmd, statusCmd)
}
