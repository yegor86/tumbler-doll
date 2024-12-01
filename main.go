package main

import (
	"os"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/yegor86/tumbler-doll/cmd"
	"github.com/yegor86/tumbler-doll/plugins"
	"github.com/yegor86/tumbler-doll/plugins/scm"
)

func main() {

	pluginManager := plugins.GetInstance()
	defer pluginManager.UnregisterAll()

	pluginManager.Register("scm", &scm.ScmPlugin{})

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
