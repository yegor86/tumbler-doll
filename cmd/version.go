package cmd

import (
	"fmt"

	cli "github.com/spf13/cobra"

	"github.com/snowzach/golib/version"
)

var (
	// Executable is overridden by Makefile with executable name
	Executable = "NoExecutable"
	// GitVersion is overridden by Makefile with git information
	GitVersion = "NoGitVersion"
)

// Version command
func init() {
	rootCmd.AddCommand(&cli.Command{
		Use:   "version",
		Short: "Show version",
		Long:  `Show version`,
		Run: func(cmd *cli.Command, args []string) {
			fmt.Println(version.Executable + " - " + version.GitVersion)
		},
	})
}
