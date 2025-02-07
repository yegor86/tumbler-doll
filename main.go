package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	temporal "go.temporal.io/sdk/client"

	"github.com/yegor86/tumbler-doll/cmd"
	"github.com/yegor86/tumbler-doll/internal/cryptography"
	"github.com/yegor86/tumbler-doll/internal/env"
	"github.com/yegor86/tumbler-doll/internal/jenkins/jobs"
	"github.com/yegor86/tumbler-doll/plugins"
	"github.com/yegor86/tumbler-doll/plugins/scm"
	"github.com/yegor86/tumbler-doll/plugins/shell"
)

func main() {
	env.LoadEnvVars()
	crypto := cryptography.GetInstance()
	crypto.LoadOrSeedCrypto()

	wfClient, err := temporal.Dial(temporal.Options{})
	if err != nil {
		log.Fatalf("Unable to create Workflow client: %v", err)
	}
	defer wfClient.Close()

	pluginManager := plugins.GetInstance()
	defer pluginManager.UnregisterAll()

	plugins := map[string]plugins.Plugin{
		"scm":   &scm.ScmPlugin{},
		"shell": &shell.ShellPlugin{},
	}

	ctx := context.WithValue(context.Background(), "wfClient", wfClient)
	for name, plugin := range plugins {
		err := pluginManager.Register(ctx, name, plugin)
		if err != nil {
			log.Printf("Failed to register plugin %s: %v", name, err)
		}
	}

	jobDb := jobs.GetInstance()
	jobs, err := jobDb.LoadJobs()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Jobs: %v", jobs)

	exitOnSyscall(pluginManager)

	cmd.Execute(wfClient)
}

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
