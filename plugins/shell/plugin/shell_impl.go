package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"github.com/docker/docker/api/types"

	"github.com/yegor86/tumbler-doll/internal/workflow"
	pb "github.com/yegor86/tumbler-doll/plugins/shell/proto"
	"github.com/yegor86/tumbler-doll/plugins/shell/shared"
)

type ShellPluginImpl struct {
	logger hclog.Logger
}

func (g *ShellPluginImpl) Echo(req *pb.LogRequest, res grpc.ServerStreamingServer[pb.LogResponse]) error {
	return g.execShell(req, res)
}

func (g *ShellPluginImpl) Sh(req *pb.LogRequest, res grpc.ServerStreamingServer[pb.LogResponse]) error {
	return g.execShell(req, res)	
}

func (g *ShellPluginImpl) execShell(req *pb.LogRequest, res grpc.ServerStreamingServer[pb.LogResponse]) error {
	g.logger.Info("[Shell] sh '%s'...", req.Command)
	terms := strings.Fields(req.Command)
	cmd := exec.Command(terms[0], terms[1:]...)
	
	
	inputStreamConsumer := func() (*bufio.Scanner, error) {
		stdout, err := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout
		if err != nil {
			g.logger.Error("[Shell] Plugin error %v", err)
			return nil, err
		}
		if err = cmd.Start(); err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(stdout)
		return scanner, nil
	}
	
	if (req.ContainerId != "") {
		var attachResp *types.HijackedResponse = nil
		inputStreamConsumer = func() (*bufio.Scanner, error) {
			var err error
			attachResp, err = workflow.ExecContainer(context.Background(), req.ContainerId, terms)

			if err != nil {
				return nil, fmt.Errorf("error attaching to container %s: %v", req.ContainerId, err)
			}
			return bufio.NewScanner(attachResp.Reader), nil
		}
		if attachResp != nil {
			defer attachResp.Close()
		}
	}

	inStream, err := inputStreamConsumer()
	if err != nil {
		return err
	}

	err = g.readAndSendBack(inStream, res)
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (g *ShellPluginImpl) readAndSendBack(scanner *bufio.Scanner, res grpc.ServerStreamingServer[pb.LogResponse]) error {
	for scanner.Scan() {
		// Simulate streaming delay
		time.Sleep(100 * time.Millisecond)
		
		// Send back a chunk of logs
		data := removeControlChars(scanner.Bytes())
		res.Send(&pb.LogResponse{Chunk: string(data)})
	}

	return scanner.Err()
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
		Output:     os.Stderr,
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
