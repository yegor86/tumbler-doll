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

func (g *ScmPluginImpl) Checkout(args shared.CheckoutArgs) string {
	g.logger.Info("PluginImpl Checkout %s...", args.Url)

	return "Cloned repo " + args.Url + " and branch " + args.Branch
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
		"checkout": &shared.ScmPlugin{Impl: scmImpl},
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
