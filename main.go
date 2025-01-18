package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
			slog.Warn("Failed to register plugin %s: %v", name, err)
		}
	}

	jobDb := jobs.GetInstance()
	jobs, err := jobDb.LoadJobs()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Jobs: %v", jobs)

	exitOnSyscall(pluginManager)

	cmd.Execute()
}

func exitOnSyscall(pluginManager *plugins.PluginManager) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		slog.Warn("Shutting down...")

		pluginManager.UnregisterAll()

		os.Exit(0)
	}()
}
