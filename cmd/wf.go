package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	temporal "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"golang.org/x/net/context"

	cli "github.com/spf13/cobra"

	"github.com/yegor86/tumbler-doll/internal/workflow"
	"github.com/yegor86/tumbler-doll/plugins"
	"github.com/yegor86/tumbler-doll/plugins/scm"
	"github.com/yegor86/tumbler-doll/plugins/shell"
)

func init() {
	rootCmd.AddCommand(workflowCmd)
}

var (
	workflowCmd = &cli.Command{
		Use:   "worker",
		Short: "Run Worker",
		Long:  "Run Worker",
		Run: func(cmd *cli.Command, args []string) {
			wfClient, ok := cmd.Context().Value("wfClient").(temporal.Client)
			if !ok {
				log.Fatalf("Failed to obtain temporal client")
			}
			pluginManager := plugins.GetInstance()
			defer pluginManager.UnregisterAll()

			plugins := map[string]plugins.Plugin{
				"scm":   &scm.ScmPlugin{},
				"shell": &shell.ShellPlugin{},
			}
			
			ctx := context.WithValue(context.Background(), "temporalHostport", os.Getenv("TEMPORAL_HOSTPORT"))
			
			for name, plugin := range plugins {
				err := pluginManager.Register(ctx, name, plugin)
				if err != nil {
					log.Printf("Failed to register plugin %s: %v", name, err)
				}
			}
			exitOnSyscall(pluginManager)
			

			w := worker.New(wfClient, "JobQueue", worker.Options{})

			w.RegisterWorkflow(workflow.GroovyDSLWorkflow)
			w.RegisterActivity(&workflow.StageActivities{})

			// Start the worker
			err := w.Run(worker.InterruptCh())
			if err != nil {
				log.Fatalf("Unable to start worker: %v", err)
			}
		},
	}
)

func exitOnSyscall(pluginManager *plugins.PluginManager) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		log.Printf("Shutting down...")

		pluginManager.UnregisterAll()

		os.Exit(0)
	}()
}