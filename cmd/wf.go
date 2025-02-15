package cmd

import (
	"log"

	temporal "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

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
			wfClient, ok := cmd.Context().Value("wfClient").(temporal.Client)
			if !ok {
				log.Fatalf("Failed to obtain temporal client")
			}

			w := worker.New(wfClient, "JobQueue", worker.Options{})

			w.RegisterWorkflow(workflow.GroovyDSLWorkflow)
			w.RegisterActivity(&workflow.StageActivities{})

			err := w.Run(worker.InterruptCh())
			if err != nil {
				log.Fatalln("Unable to start worker", err)
			}
		},
	}
)
