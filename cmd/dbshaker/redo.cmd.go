package main

import (
	"github.com/ToggyO/dbshaker/internal"
	"github.com/spf13/cobra"
)

var redoCmd = &cobra.Command{
	Use:   internal.CmdRedo,
	Short: "repeat the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO:
	},
}
