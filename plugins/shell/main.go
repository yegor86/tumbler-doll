package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/yegor86/tumbler-doll/plugins/shell/shared"
)

type ShellPlugin struct {
	shell  shared.Shell
	client *plugin.Client
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SHELL_PLUGIN",
	MagicCookieValue: "shell",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"shell": &shared.ShellPlugin{},
}

func (p *ShellPlugin) Start() error {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("plugins/shell/shell"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
		Logger: logger,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("shell")
	if err != nil {
		return err
	}

	p.shell = raw.(shared.Shell)
	p.client = client
	return nil
}

func (p *ShellPlugin) Stop() error {
	if p.client == nil {
		return fmt.Errorf("bash plugin is not initialized")
	}
	p.client.Kill()
	return nil
}

func (p *ShellPlugin) ListMethods() map[string]string {
	return map[string]string{
		"echo": "echo",
		"sh":   "sh",
	}
}

func (scmClient *ShellPlugin) Echo(args map[string]interface{}) string {
	err := scmClient.shell.Echo(args)
	return err.Error()
}

func (scmClient *ShellPlugin) Sh(args map[string]interface{}) string {
	err := scmClient.shell.Sh(args)
	return err.Error()
}
