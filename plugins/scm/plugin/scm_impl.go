package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/yegor86/tumbler-doll/internal/cryptography"
	"github.com/yegor86/tumbler-doll/internal/jenkins/xml"
	"github.com/yegor86/tumbler-doll/plugins/scm/shared"
)

type ScmPluginImpl struct {
	logger hclog.Logger
	git    Git
}

var (
	authMethods = map[string]func(credentials *xml.Credential) transport.AuthMethod{
		"ssh": func(credentials *xml.Credential) transport.AuthMethod {
			if credentials == nil {
				return nil
			}
			args := credentials.Tags

			var passphrase = ""
			if _, ok := args["passphrase"]; ok {
				passphrase = args["passphrase"]
			}
			var privateKey = ""
			if _, ok := args["privateKey"]; ok {
				privateKey = args["privateKey"]
			}
			var username = "git"
			if _, ok := args["username"]; ok {
				username = args["username"]
			}

			publicKey, err := ssh.NewPublicKeys(username, []byte(privateKey), passphrase)
			if err != nil {
				log.Printf("error generating public key: %v\n", err)
				return nil
			}
			return publicKey
		},
	}
)

func (g *ScmPluginImpl) Checkout(args map[string]interface{}) (string, error) {
	if _, ok := args["url"]; !ok {
		return "", fmt.Errorf("url is missing")
	}
	if _, ok := args["branch"]; !ok {
		return "", fmt.Errorf("git branch is missing")
	}

	url := args["url"].(string)
	branch := args["branch"].(string)
	var credentials *xml.Credential = nil
	if _, ok := args["credentialsId"]; ok {
		credentialsId := args["credentialsId"].(string)
		crypto := cryptography.GetInstance()
		credentials = crypto.GetCredentialsById(credentialsId)
		if credentials == nil {
			return "", fmt.Errorf("credentials not found by id %s", credentialsId)
		}

		g.logger.Info("PluginImpl credentialsId:", credentialsId)
		g.logger.Info("PluginImpl credentials list size:", len(crypto.Credentials))
	}

	var authMethod transport.AuthMethod = nil
	if _, ok := authMethods["ssh"]; ok {
		generateAuth := authMethods["ssh"]
		authMethod = generateAuth(credentials)
	}

	g.logger.Info("PluginImpl Checkout %s...", url)
	g.logger.Info("PluginImpl auth method %s:%s...", authMethod.Name(), authMethod.String())

	if err := g.git.CloneOrPull(url, branch, authMethod); err != nil {
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
	crypto := cryptography.GetInstance()
	crypto.LoadOrSeedCrypto()

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Debug,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	scmImpl := &ScmPluginImpl{
		logger: logger,
		git: &GitRepo{
			Changelog: true,
			Poll: true,
			ProgressWriter: logger.StandardWriter(&hclog.StandardLoggerOptions{
				InferLevels: true,
			}),
		},
	}

	var pluginMap = map[string]plugin.Plugin{
		"scm": &shared.ScmPlugin{Impl: scmImpl},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})

	log.Println("Test output")
}
