package cmd

import (
	"fmt"

	cli "github.com/spf13/cobra"
)

var (
	// Executable and GitVersion are overridden by Makefile with executable name
	Executable = "__Version__"
	GitVersion = "__Version__"
)

// Version command
func init() {
	rootCmd.AddCommand(&cli.Command{
		Use:   "version",
		Short: "Show version",
		Long:  `Show current version`,
		Run: func(cmd *cli.Command, args []string) {
			fmt.Println(Executable + " - " + GitVersion)
		},
	})
}
