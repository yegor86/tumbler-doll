package plugin

import (
	"log/slog"

	"github.com/hashicorp/go-plugin"

	"github.com/yegor86/tumbler-doll/plugins/scm/shared"
)

type ScmImpl struct {
	logger slog.Logger
}

func (g *ScmImpl) Checkout(url string, branch string, credentialsIdL string) string {
	g.logger.Debug("message from GreeterHello.Greet")

	return "done"
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "GIT_PLUGIN",
	MagicCookieValue: "gitSCM",
}

func main() {
	scmImpl := &ScmImpl{
		logger: *slog.Default(),
	}

	var pluginMap = map[string]plugin.Plugin{
		"checkout": &shared.ScmPlugin{Impl: scmImpl},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
