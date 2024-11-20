package scm

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/yegor86/tumbler-doll/plugins/scm/shared"
)

type ScmClint struct {

}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"checkout": &shared.ScmPlugin{},
}

func (scmClient *ScmClint) checkout(url string, branch string, credentialsId string) string {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./plugin/scm"),
		Logger:          logger,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("checkout")
	if err != nil {
		log.Fatal(err)
	}

	scm := raw.(shared.Scm)
	return scm.Checkout(url, branch, credentialsId)
}