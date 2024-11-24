package scm

import (
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/yegor86/tumbler-doll/plugins/scm/shared"
)

type ScmClint struct {
	scm shared.Scm
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "GIT_PLUGIN",
	MagicCookieValue: "gitSCM",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"checkout": &shared.ScmPlugin{},
}

func NewScmClient() *shared.ScmRPCClient {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	os.Setenv(handshakeConfig.MagicCookieKey, handshakeConfig.MagicCookieValue)
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("plugins/scm/scm"),
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
		log.Panic(err)
	}

	return raw.(*shared.ScmRPCClient)
}

func (scmClient *ScmClint) Checkout(args shared.CheckoutArgs) string {
	return scmClient.scm.Checkout(args)
}
