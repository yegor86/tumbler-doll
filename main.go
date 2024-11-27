package main

import (
	"fmt"
	"os"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/yegor86/tumbler-doll/cmd"
	"github.com/yegor86/tumbler-doll/plugins"
	"github.com/yegor86/tumbler-doll/plugins/scm"
	"github.com/yegor86/tumbler-doll/plugins/scm/shared"
)

func main() {

	pluginManager := plugins.NewPluginManager()
	defer pluginManager.UnregisterAll()

	pluginManager.Register("scm", &scm.ScmPlugin{})

	args := shared.CheckoutArgs{
		Url:           "http://testurl.com",
		Branch:        "master",
		CredentialsId: "",
	}
	res, err := pluginManager.Execute("scm", "Checkout", args)
	if err != nil {
		fmt.Printf("Error executing scm.Checkout %v", err)
	}
	fmt.Println(res)

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
