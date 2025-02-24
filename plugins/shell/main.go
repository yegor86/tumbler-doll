package shell

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/yegor86/tumbler-doll/internal/grpc"
	"github.com/yegor86/tumbler-doll/plugins"
	"github.com/yegor86/tumbler-doll/plugins/shell/shared"

	logstream "github.com/yegor86/tumbler-doll/internal/grpc/proto"
	pb "github.com/yegor86/tumbler-doll/plugins/shell/proto"
)

type ShellPlugin struct {
	shell  shared.ClientShell
	pluginClient *plugin.Client
	streamClient *grpc.GrpcClient
	ctx context.Context
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SHELL_PLUGIN",
	MagicCookieValue: "shell",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"shell": &shared.ShellPlugin{},
}

func (p *ShellPlugin) Start(ctx context.Context) error {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	temporalHostPort, ok := ctx.Value("temporalHostport").(string)
	if !ok {
		return fmt.Errorf("failed to extract TEMPORAL_ADDRESS from context: %v", ctx.Value("temporalHostport"))
	}
	streamClient, err := grpc.NewClient(temporalHostPort)
	if err != nil {
		return err
	}
	p.streamClient = streamClient

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("plugins/shell/shell"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
		Logger: logger,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("shell")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	p.shell = raw.(shared.ClientShell)
	p.pluginClient = client
	p.ctx = ctx
	return nil
}

func (p *ShellPlugin) Stop() error {
	if p.pluginClient == nil {
		return errors.New("shell plugin is not initialized")
	}
	p.pluginClient.Kill()
	p.streamClient.CloseStream()
	
	return nil
}

func (p *ShellPlugin) ListMethods() map[string]string {
	return map[string]string{
		"echo": "echo",
		"sh":   "sh",
	}
}

func (scmClient *ShellPlugin) Echo(args map[string]interface{}) error {
	workflowExecutionId := ""
	if _, ok := args["workflowExecutionId"]; ok {
		workflowExecutionId = args["workflowExecutionId"].(string)
	}

	serverStream, err := scmClient.shell.Echo(scmClient.ctx, args)
	if err != nil {
		return err
	}
	return plugins.Redirect(serverStream, scmClient.streamClient.Stream, func(resp *pb.ShellResponse) *logstream.LogRequest {
		return &logstream.LogRequest{
			Message: resp.Chunk,
			WorkflowId: workflowExecutionId,
		}
	})
}

func (scmClient *ShellPlugin) Sh(args map[string]interface{}) error {
	workflowExecutionId := ""
	if _, ok := args["workflowExecutionId"]; ok {
		workflowExecutionId = args["workflowExecutionId"].(string)
	}

	serverStream, err := scmClient.shell.Sh(scmClient.ctx, args)
	if err != nil {
		return err
	}
	return plugins.Redirect(serverStream, scmClient.streamClient.Stream, func(resp *pb.ShellResponse) *logstream.LogRequest {
		return &logstream.LogRequest{
			Message: resp.Chunk,
			WorkflowId: workflowExecutionId,
		}
	})
}
