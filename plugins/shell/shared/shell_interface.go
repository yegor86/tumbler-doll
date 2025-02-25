package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	pb "github.com/yegor86/tumbler-doll/plugins/shell/proto"
)

type ClientShell interface {
	Echo(ctx context.Context, args map[string]interface{}) (grpc.ServerStreamingClient[pb.ShellResponse], error)
	Sh(ctx context.Context, args map[string]interface{}) (grpc.ServerStreamingClient[pb.ShellResponse], error)
}

type ServerShell interface {
	Echo(request *pb.ShellRequest, response grpc.ServerStreamingServer[pb.ShellResponse]) error
	Sh(request *pb.ShellRequest, response grpc.ServerStreamingServer[pb.ShellResponse]) error
}

type ShellRPCClient struct {
	client   pb.ShellStreamingServiceClient
	broker   *plugin.GRPCBroker
}

func (g *ShellRPCClient) Echo(ctx context.Context, args map[string]interface{}) (grpc.ServerStreamingClient[pb.ShellResponse], error) {

	cmd := "echo " + args["text"].(string)
	containerId, _ := args["containerId"].(string)
	
	return g.client.Echo(context.Background(), &pb.ShellRequest{
		Command:     cmd,
		ContainerId: containerId,
	})
}

func (g *ShellRPCClient) Sh(ctx context.Context, args map[string]interface{}) (grpc.ServerStreamingClient[pb.ShellResponse], error) {
	cmd := args["text"].(string)
	containerId, _ := args["containerId"].(string)

	return g.client.Sh(ctx, &pb.ShellRequest{
		Command:     cmd,
		ContainerId: containerId,
	})
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
	Impl ClientShell
}

func (p *ShellPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ShellRPCClient{
		client: pb.NewShellStreamingServiceClient(c),
		broker: broker,
	}, nil
}
