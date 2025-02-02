package shared

import (
	"context"
	"fmt"
	
	"google.golang.org/grpc"
	"github.com/hashicorp/go-plugin"
	pb "github.com/yegor86/tumbler-doll/plugins/shell/proto"
)

type Shell interface {
	Echo(args map[string]interface{}) error
	Sh(args map[string]interface{}) error
}

type ServerShell interface {
	Echo(request *pb.LogRequest, response grpc.ServerStreamingServer[pb.LogResponse]) error
	Sh(request *pb.LogRequest, response grpc.ServerStreamingServer[pb.LogResponse]) error
}

// Here is an implementation that talks over RPC
type ShellRPCClient struct{
	client pb.LogStreamingServiceClient
	broker *plugin.GRPCBroker
}

func (g *ShellRPCClient) Echo(args map[string]interface{}) error {
	// err := g.client.Call("Plugin.Echo", args, reply)
	
	cmd := "echo " + args["text"].(string)
	containerId := ""
	if _, ok := args["containerId"]; ok {
		containerId = args["containerId"].(string)
	}
	stream, err := g.client.Echo(context.Background(), &pb.LogRequest{
		Command: cmd,
		ContainerId: containerId,
	})
	if err != nil {
		return err
	}

	var streamErr error
	for {
		resp, err := stream.Recv()
		if err != nil {
			streamErr = err
			break // Stream finished
		}
		fmt.Println(resp.Chunk)
	}
	
	return streamErr
}

func (g *ShellRPCClient) Sh(args map[string]interface{}) error {
	// err := g.client.Call("Plugin.Sh", args, reply)
	cmd := args["text"].(string)
	containerId := ""
	// if _, ok := args["containerId"]; ok {
	// 	containerId = args["containerId"].(string)
	// }
	stream, err := g.client.Sh(context.Background(), &pb.LogRequest{
		Command: cmd,
		ContainerId: containerId,
	})
	if err != nil {
		return err
	}

	var streamErr error
	for {
		resp, err := stream.Recv()
		if err != nil {
			streamErr = err
			break
		}
		fmt.Println(resp.Chunk)
	}
	
	return streamErr
}

type ShellRPCServer struct {
	pb.UnsafeLogStreamingServiceServer
	Impl   ServerShell
	broker *plugin.GRPCBroker
}

func (s *ShellRPCServer) Echo(request *pb.LogRequest, response grpc.ServerStreamingServer[pb.LogResponse]) error {
	return s.Impl.Echo(request, response)
}

func (s *ShellRPCServer) Sh(request *pb.LogRequest, response grpc.ServerStreamingServer[pb.LogResponse]) error {
	return s.Impl.Echo(request, response)
}

type ShellPlugin struct {
	plugin.GRPCPlugin
	plugin.NetRPCUnsupportedPlugin
	Impl Shell
}

type ServerShellPlugin struct {
	plugin.GRPCPlugin
	plugin.NetRPCUnsupportedPlugin
	Impl ServerShell
}

func (p *ServerShellPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterLogStreamingServiceServer(s, &ShellRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *ShellPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ShellRPCClient{
		client: pb.NewLogStreamingServiceClient(c), 
		broker: broker,
	}, nil
}
