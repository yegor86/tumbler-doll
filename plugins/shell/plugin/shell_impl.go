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

	return g.Sh(&pb.LogRequest{
		Command:     "echo " + req.Command,
		ContainerId: req.ContainerId,
	}, res)
}

func (g *ShellPluginImpl) Sh(req *pb.LogRequest, res grpc.ServerStreamingServer[pb.LogResponse]) error {

	next := func(cmd, containerId string) (*bufio.Scanner, func() error, error) {
		g.logger.Info("[Shell] sh '%s'...", cmd)
		terms := strings.Fields(cmd)

		execCommand := exec.Command(terms[0], terms[1:]...)
		stdout, err := execCommand.StdoutPipe()
		execCommand.Stderr = execCommand.Stdout
		if err != nil {
			g.logger.Error("[Shell] Plugin error %v", err)
			return nil, nil, err
		}
		if err = execCommand.Start(); err != nil {
			return nil, nil, err
		}
		scanner := bufio.NewScanner(stdout)
		return scanner, func() error {
			return execCommand.Wait()
		}, nil
	}
	containerized := workflow.Containerize(next)
	scanner, waitFunc, err := containerized(req.Command, req.ContainerId)
	if err != nil {
		return err
	}
	err = readAll(scanner, res)
	if err != nil && err != io.EOF {
		return nil
	}

	err = waitFunc()
	return err
}

func readAll(scanner *bufio.Scanner, res grpc.ServerStreamingServer[pb.LogResponse]) error {
	for scanner.Scan() {
		// Send back a chunk of logs
		res.Send(&pb.LogResponse{Chunk: scanner.Text()})
		time.Sleep(100 * time.Millisecond) // Simulate streaming delay
	}

	// Check for errors in scanner
	if err := scanner.Err(); err != nil {
		return err
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
		"shell": &shared.ServerShellPlugin{Impl: shellImpl},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
