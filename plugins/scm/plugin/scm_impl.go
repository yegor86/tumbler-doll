package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/yegor86/tumbler-doll/plugins/scm/shared"
)

type ScmPluginImpl struct {
	logger hclog.Logger
}

func (g *ScmPluginImpl) Checkout(args map[string]interface{}) string {
	url := args["url"]
	branch := args["branch"]
	// credentialsId, _ := args["credentialsId"]
	g.logger.Info("PluginImpl Checkout %s...", url)

	return "Cloned repo " + url.(string) + " and branch " + branch.(string)
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "GIT_PLUGIN",
	MagicCookieValue: "gitSCM",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Debug,
		Output:     os.Stdout,
		JSONFormat: true,
	})

	scmImpl := &ScmPluginImpl{
		logger: logger,
	}

	var pluginMap = map[string]plugin.Plugin{
		"scm": &shared.ScmPlugin{Impl: scmImpl},
	}

	// logger.Warn("[scp] message from plugin")
	// logger.Warn("[scp] yet another message")
	os.Setenv(handshakeConfig.MagicCookieKey, handshakeConfig.MagicCookieValue)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
		// Logger:          logger,
	})
}
