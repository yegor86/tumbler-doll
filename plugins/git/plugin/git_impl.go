package main

import (
	"log/slog"
	"github.com/hashicorp/go-plugin"

	"github.com/yegor86/tumbler-doll/plugins/git/shared"
)

type GitPlugin struct {
	logger slog.Logger
}

func (g *GitPlugin) checkout(url string, branch string, credentialsIdL string) string {
	g.logger.Debug("message from GreeterHello.Greet")
	
	
	
	return "done"
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "GIT_PLUGIN",
	MagicCookieValue: "gitSCM",
}

func main() {
	greeter := &GitPlugin{
		logger: *slog.Default(),
	}
	
	var pluginMap = map[string]plugin.Plugin{
		"greeter": &shared.GreeterPlugin{Impl: greeter},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}