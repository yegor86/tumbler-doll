package main

import (
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

	pluginManager := plugins.GetInstance()
	defer pluginManager.UnregisterAll()

	plugins := map[string]plugins.Plugin{
		"scm":   &scm.ScmPlugin{},
		"shell": &shell.ShellPlugin{},
	}

	for name, plugin := range plugins {
		err := pluginManager.Register(name, plugin)
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

	client, err := temporal.Dial(temporal.Options{})
	if err != nil {
		log.Fatalf("Unable to create Workflow client", err)
	}
	defer client.Close()

	cmd.Execute(client)
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
