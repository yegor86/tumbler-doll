package scm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/yegor86/tumbler-doll/plugins/scm/shared"
)

type ScmPlugin struct {
	scm    shared.Scm
	client *plugin.Client
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

func (p *ScmPlugin) Start() error {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("plugins/scm/scm"),
		Logger:          logger,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("checkout")
	if err != nil {
		return err
	}

	p.scm = raw.(shared.Scm)
	p.client = client
	return nil
}

func (p *ScmPlugin) Stop() error {
	if p.client == nil {
		return fmt.Errorf("scm plugin is not initialized")
	}
	p.client.Kill()
	return nil
}

func (scmClient *ScmPlugin) Checkout(args map[string]interface{}) string {
	return scmClient.scm.Checkout(args)
}