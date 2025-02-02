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
	"google.golang.org/grpc"

	"github.com/yegor86/tumbler-doll/internal/workflow"
	pb "github.com/yegor86/tumbler-doll/plugins/shell/proto"
	"github.com/yegor86/tumbler-doll/plugins/shell/shared"
)

type ShellPluginImpl struct {
	logger hclog.Logger
}

func (g *ShellPluginImpl) Echo(req *pb.LogRequest, res grpc.ServerStreamingServer[pb.LogResponse]) error {

	g.logger.Info("[Shell] echo %s...", req.Command)
	text := req.Command

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
	return readAll(reader, res)
}

func (g *ShellPluginImpl) Sh(req *pb.LogRequest, res grpc.ServerStreamingServer[pb.LogResponse]) error {

	next := func (cmd, containerId string) (*bufio.Reader, error) {
		g.logger.Info("[Shell] sh '%s'...", cmd)
		terms := strings.Fields(cmd)

		execCommand := exec.Command(terms[0], terms[1:]...)
		stdout, err := execCommand.StdoutPipe()
		if err != nil {
			g.logger.Error("[Shell] Plugin error %v", err)
			return nil, err
		}
		if err = execCommand.Start(); err != nil {
			return nil, err
		}
		defer execCommand.Wait()

		reader := bufio.NewReader(stdout)
		return reader, nil
	}
	containerized := workflow.Containerize(next)
	reader, err := containerized(req.Command, req.ContainerId);
	if err != nil {
		return err
	}
	return readAll(reader, res)
}

func readAll(reader *bufio.Reader, res grpc.ServerStreamingServer[pb.LogResponse]) error {
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return err
		}
		if err != nil {
			return err
		}

		// Send back a chunk of logs
		res.Send(&pb.LogResponse{Chunk: string(line)})
		time.Sleep(100 * time.Millisecond) // Simulate streaming delay
	}
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
		"shell": &shared.ServerShellPlugin{Impl: shellImpl},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
