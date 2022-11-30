package main

import (
	"github.com/spf13/cobra"

	"github.com/ToggyO/dbshaker/internal"
)

var redoCmd = &cobra.Command{
	Use:   internal.CmdRedo,
	Short: "repeat the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO:
	},
}
