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

	docker "github.com/yegor86/tumbler-doll/plugins/docker/shared"
	"github.com/yegor86/tumbler-doll/plugins/shell/shared"
	pb "github.com/yegor86/tumbler-doll/plugins/shell/proto"
)

type ShellPluginImpl struct {
	logger hclog.Logger
	docker docker.DockerClient
}

func (g *ShellPluginImpl) Echo(req *pb.ShellRequest, res grpc.ServerStreamingServer[pb.ShellResponse]) error {
	return g.execShell(req, res)
}

func (g *ShellPluginImpl) Sh(req *pb.ShellRequest, res grpc.ServerStreamingServer[pb.ShellResponse]) error {
	return g.execShell(req, res)	
}

func (g *ShellPluginImpl) execShell(req *pb.ShellRequest, res grpc.ServerStreamingServer[pb.ShellResponse]) error {
	g.logger.Info("[Shell] sh '%s'...", req.Command)
	terms := strings.Fields(req.Command)
	
	cmd := exec.Command(terms[0], terms[1:]...)
	inputStreamConsumer, closeStreamConsumer := func() (*bufio.Scanner, error) {
		stdout, err := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout
		if err != nil {
			return nil, err
		}
		if err = cmd.Start(); err != nil {
			return nil, err
		}
		return bufio.NewScanner(stdout), nil
	}, func() error {
		return cmd.Wait()
	}
	
	if (req.ContainerId != "") {
		var attachResp *types.HijackedResponse = nil
		inputStreamConsumer, closeStreamConsumer = func() (*bufio.Scanner, error) {
			var err error
			attachResp, err = g.docker.ExecContainer(context.Background(), req.ContainerId, terms)

			if err != nil {
				return nil, fmt.Errorf("error attaching to container %s: %v", req.ContainerId, err)
			}
			return bufio.NewScanner(attachResp.Reader), nil
		}, func() error {
			attachResp.Close()
			return nil
		}
	}
	
	inStream, err := inputStreamConsumer()
	if err != nil {
		g.logger.Error("[Shell] Plugin.inputStreamConsumer error %v", err)
		return err
	}
	defer closeStreamConsumer()

	err = g.readAndSendBack(inStream, res)
	if err != nil {
		g.logger.Error("[Shell] Plugin.readAndSendBack error %v", err)
	}

	return err
}

func (g *ShellPluginImpl) readAndSendBack(scanner *bufio.Scanner, res grpc.ServerStreamingServer[pb.ShellResponse]) error {
	for scanner.Scan() {
		// Simulate streaming delay
		time.Sleep(100 * time.Millisecond)
		
		// Send back a chunk of logs
		data := removeControlChars(scanner.Bytes())
		res.Send(&pb.ShellResponse{Chunk: string(data)})
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

	dockerClient, err := docker.NewDockerClient(context.Background())
	if err != nil {
		logger.Warn("error initializing docker client %v", err)
	}
	defer dockerClient.Stop()

	shellImpl := &ShellPluginImpl{
		docker: dockerClient,
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
