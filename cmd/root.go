package cmd

import (
	"fmt"
	"log"
	"os"

	cli "github.com/spf13/cobra"
	temporal "go.temporal.io/sdk/client"
	"golang.org/x/net/context"
)

func init() {
	// Parse defaults, config file and environment.
	_, _, err := Load()
	if err != nil {
		log.Fatalf("could not parse YAML config: %v", err)
	}
}

var (
	// The Root Cli Handler
	rootCmd = &cli.Command{
		Version: GitVersion,
		Use:     Executable,
	}
)

// Execute starts the program
func Execute(client temporal.Client) {
	ctx := context.WithValue(context.Background(), "wfClient", client)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
