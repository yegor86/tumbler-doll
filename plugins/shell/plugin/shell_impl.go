package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/yegor86/tumbler-doll/internal/workflow"
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

func (g *ShellPluginImpl) Sh(params map[string]interface{}) string {

	g.logger.Info("[Shell] sh '%s'...", params["text"])
	text := params["text"].(string)
	terms := strings.Fields(text)
	if containerId, ok := params["containerId"]; ok {
		output, err := workflow.ExecContainer(context.Background(), containerId.(string), terms)
		
		if err != nil {
			g.logger.Error("Error attaching to the container %s: %v. Continue running command on host machine", containerId, err)
			return fmt.Errorf("error attaching to container %s: %v", containerId, err).Error()
		}
		return string(output)
	}

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
