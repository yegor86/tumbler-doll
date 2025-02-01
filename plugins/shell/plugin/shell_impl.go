package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/yegor86/tumbler-doll/internal/workflow"
	"github.com/yegor86/tumbler-doll/plugins/shell/shared"
)

type ShellPluginImpl struct {
	logger hclog.Logger
}

func (g *ShellPluginImpl) Echo(args map[string]interface{}, reply *shared.StreamLogsReply) error {

	g.logger.Info("[Shell] echo %s...", args["text"])
	text := args["text"].(string)

	cmd := exec.Command("echo", text)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		g.logger.Error("[Shell] Plugin error %v", err)
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	defer cmd.Wait()

	reader := bufio.NewReader(stdout)
	return readAll(reader, reply)
}

func (g *ShellPluginImpl) Sh(params map[string]interface{}, reply *shared.StreamLogsReply) error {

	next := func (params map[string]interface{}) (*bufio.Reader, error) {
		g.logger.Info("[Shell] sh '%s'...", params["text"])
		text := params["text"].(string)
		terms := strings.Fields(text)

		cmd := exec.Command(terms[0], terms[1:]...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			g.logger.Error("[Shell] Plugin error %v", err)
			return nil, err
		}
		if err = cmd.Start(); err != nil {
			return nil, err
		}
		defer cmd.Wait()

		reader := bufio.NewReader(stdout)
		return reader, nil
	}
	containerized := workflow.Containerize(next)
	reader, err := containerized(params);
	if err != nil {
		return err
	}
	return readAll(reader, reply)
}

func readAll(reader *bufio.Reader, reply *shared.StreamLogsReply) error {
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return err
		}
		if err != nil {
			return err
		}

		// Send back a chunk of logs
		reply.Chunk = string(line)
		time.Sleep(100 * time.Millisecond) // Simulate streaming delay
	}
	return nil
}

// RemoveControlChars removes non-printable ASCII characters from byte array and return human readble string.
func removeControlChars(input []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, input)
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
