package shared

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	pb "github.com/yegor86/tumbler-doll/plugins/shell/proto"
)

type Shell interface {
	Echo(args map[string]interface{}) error
	Sh(args map[string]interface{}) error
}

type ServerShell interface {
	Echo(request *pb.ShellRequest, response grpc.ServerStreamingServer[pb.ShellResponse]) error
	Sh(request *pb.ShellRequest, response grpc.ServerStreamingServer[pb.ShellResponse]) error
}

// Here is an implementation that talks over RPC
type ShellRPCClient struct {
	client   pb.ShellStreamingServiceClient
	broker   *plugin.GRPCBroker
}

func (g *ShellRPCClient) Echo(args map[string]interface{}) error {
	// err := g.client.Call("Plugin.Echo", args, reply)

	cmd := "echo " + args["text"].(string)
	containerId := ""
	if _, ok := args["containerId"]; ok {
		containerId = args["containerId"].(string)
	}
	// workflowExecutionId := ""
	// if _, ok := args["workflowExecutionId"]; ok {
	// 	workflowExecutionId = args["workflowExecutionId"].(string)
	// }
	// logSignalName := ""
	// if _, ok := args["logSignalName"]; ok {
	// 	workflowExecutionId = args["logSignalName"].(string)
	// }

	stream, err := g.client.Echo(context.Background(), &pb.ShellRequest{
		Command:     cmd,
		ContainerId: containerId,
	})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		// if g.wfClient != nil && workflowExecutionId != "" && logSignalName != "" {
		// 	err = g.wfClient.SignalWorkflow(ctx, workflowExecutionId, "", logSignalName, resp.Chunk)
		// }
		// if err != nil {
		// 	return err
		// }
		fmt.Println(resp.Chunk)
	}
}

func (g *ShellRPCClient) Sh(args map[string]interface{}) error {
	ctx := context.Background()
	cmd := args["text"].(string)
	containerId := ""
	if _, ok := args["containerId"]; ok {
		containerId = args["containerId"].(string)
	}
	// workflowExecutionId := ""
	// if _, ok := args["workflowExecutionId"]; ok {
	// 	workflowExecutionId = args["workflowExecutionId"].(string)
	// }
	// logSignalName := ""
	// if _, ok := args["logSignalName"]; ok {
	// 	workflowExecutionId = args["logSignalName"].(string)
	// }

	stream, err := g.client.Sh(ctx, &pb.ShellRequest{
		Command:     cmd,
		ContainerId: containerId,
	})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		// if g.wfClient != nil && workflowExecutionId != "" && logSignalName != "" {
		// 	err = g.wfClient.SignalWorkflow(ctx, workflowExecutionId, "", logSignalName, resp.Chunk)
		// }
		// if err != nil {
		// 	return err
		// }
		fmt.Println(resp.Chunk)
	}
}

type ShellRPCServer struct {
	pb.UnsafeShellStreamingServiceServer
	Impl   ServerShell
	broker *plugin.GRPCBroker
}

func (s *ShellRPCServer) Echo(request *pb.ShellRequest, response grpc.ServerStreamingServer[pb.ShellResponse]) error {
	return s.Impl.Echo(request, response)
}

func (s *ShellRPCServer) Sh(request *pb.ShellRequest, response grpc.ServerStreamingServer[pb.ShellResponse]) error {
	return s.Impl.Echo(request, response)
}

type ServerShellPlugin struct {
	plugin.GRPCPlugin
	plugin.NetRPCUnsupportedPlugin
	Impl ServerShell
}

func (p *ServerShellPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterShellStreamingServiceServer(s, &ShellRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

type ShellPlugin struct {
	plugin.GRPCPlugin
	plugin.NetRPCUnsupportedPlugin
	Impl Shell
}

func (p *ShellPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ShellRPCClient{
		client: pb.NewShellStreamingServiceClient(c),
		broker: broker,
	}, nil
}
