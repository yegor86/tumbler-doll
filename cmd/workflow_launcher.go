package cmd

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	_ "github.com/jackc/pgx/v4/stdlib"
	cli "github.com/spf13/cobra"

	"github.com/yegor86/tumbler-doll/internal/workflow"
)

func init() {
	rootCmd.AddCommand(workflowCmd)
}

var (
	workflowCmd = &cli.Command{
		Use:   "wf",
		Short: "Start Workflow",
		Long:  `Start Workflow`,
		Run: func(cmd *cli.Command, args []string) {
			// The client and worker are heavyweight objects that should be created once per process.
			c, err := client.Dial(client.Options{})
			if err != nil {
				log.Fatalln("Unable to create Workflow client", err)
			}
			defer c.Close()

			w := worker.New(c, "dsl", worker.Options{})

			w.RegisterWorkflow(workflow.SimpleDSLWorkflow)
			w.RegisterActivity(&workflow.SampleActivities{})

			err = w.Run(worker.InterruptCh())
			if err != nil {
				log.Fatalln("Unable to start worker", err)
			}
		},
	}
)
