package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/yegor86/tumbler-doll/plugins/scm/shared"
)

type ScmPluginImpl struct {
	logger hclog.Logger
}

func (g *ScmPluginImpl) Checkout(args map[string]interface{}) (string, error) {
	if _, ok := args["url"]; !ok {
		return "", fmt.Errorf("url is missing")
	}
	if _, ok := args["branch"]; !ok {
		return "", fmt.Errorf("branch is missing")
	}
	
	url := args["url"].(string)
	branch := args["branch"].(string)
	// credentialsId, _ := args["credentialsId"]
	g.logger.Info("PluginImpl Checkout %s...", url)

	cloneDir, err := shared.DeriveCloneDir(url)
	if err != nil {
		return "", err
	}
	
	repo := &shared.GitRepo {
		Url: url,
		Branch: branch,
		CloneDir: "/tmp/" + cloneDir,
		Changelog: true,
		Credentials: "",
		Poll: true,
	}
	if err := repo.CloneOrPull(); err != nil {
		return "", err
	}

	return fmt.Sprintf("Cloned repo %s and branch %s", url, branch), nil
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
