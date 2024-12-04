package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/yegor86/tumbler-doll/plugins/shell/shared"
)

type ShellPluginImpl struct {
	logger hclog.Logger
}

func (g *ShellPluginImpl) Echo(args map[string]interface{}) string {

	g.logger.Info("[Shell] echo %s...", args["text"])
	text := args["text"].(string)

	cmd := exec.Command("echo", text)
	result, err := cmd.Output()
	if err != nil {
		g.logger.Error("[Shell] Plugin error %v", err)
		return err.Error()
	}

	return string(result)
}

func (g *ShellPluginImpl) Sh(args map[string]interface{}) string {

	g.logger.Info("[Shell] sh '%s'...", args["text"])
	text := args["text"].(string)

	terms := strings.Fields(text)
	cmd := exec.Command(terms[0], terms[1:]...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		g.logger.Error("[Shell] Plugin error %v", err)
		return err.Error()
	}

	return string(result)
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SHELL_PLUGIN",
	MagicCookieValue: "shell",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Debug,
		Output:     os.Stdout,
		JSONFormat: true,
	})

	shellImpl := &ShellPluginImpl{
		logger: logger,
	}

	var pluginMap = map[string]plugin.Plugin{
		"shell": &shared.ShellPlugin{Impl: shellImpl},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
