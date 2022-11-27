package main

import (
	"github.com/ToggyO/dbshaker/internal"
	"github.com/spf13/cobra"
)

const (
	directoryCmdArg      = "dir"
	directoryCmdArgS     = "d"
	directoryCmdArgUsage = "Path to folder, contains migrations"
)

var rootCmd = cobra.Command{
	Use:              internal.ToolName,
	TraverseChildren: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd, migrateCmd)
	createCmd.PersistentFlags().StringP(directoryCmdArg, directoryCmdArgS, "", directoryCmdArgUsage)
}
