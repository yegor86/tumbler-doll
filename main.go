package main

import (
	"os"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/yegor86/tumbler-doll/cmd"
	"github.com/yegor86/tumbler-doll/plugins"
	"github.com/yegor86/tumbler-doll/plugins/scm"
	"github.com/yegor86/tumbler-doll/plugins/shell"
)

func main() {

	pluginManager := plugins.GetInstance()
	defer pluginManager.UnregisterAll()

	err := pluginManager.Register("scm", &scm.ScmPlugin{})
	if err != nil {
		slog.Warn("Failed to register plugin %s: %v", "scm", err)
	}
	err = pluginManager.Register("shell", &shell.ShellPlugin{})
	if err != nil {
		slog.Warn("Failed to register plugin %s: %v", "shell", err)
	}

	signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
        <-signals
        slog.Warn("Shutting down...")

        pluginManager.UnregisterAll()

        os.Exit(0)
    }()

	cmd.Execute()
}
