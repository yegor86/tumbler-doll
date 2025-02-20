package shared

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	stream "github.com/yegor86/tumbler-doll/internal/grpc"
	pb "github.com/yegor86/tumbler-doll/plugins/shell/proto"
)

type ClientShell interface {
	Echo(ctx context.Context, args map[string]interface{}, streamClient *stream.GrpcClient) error
	Sh(ctx context.Context, args map[string]interface{}, streamClient *stream.GrpcClient) error
}

type ServerShell interface {
	Echo(request *pb.ShellRequest, response grpc.ServerStreamingServer[pb.ShellResponse]) error
	Sh(request *pb.ShellRequest, response grpc.ServerStreamingServer[pb.ShellResponse]) error
}

type ShellRPCClient struct {
	client   pb.ShellStreamingServiceClient
	broker   *plugin.GRPCBroker
}

func (g *ShellRPCClient) Echo(ctx context.Context, args map[string]interface{}, streamClient *stream.GrpcClient) error {
	// err := g.client.Call("Plugin.Echo", args, reply)

	cmd := "echo " + args["text"].(string)
	containerId := ""
	if _, ok := args["containerId"]; ok {
		containerId = args["containerId"].(string)
	}
	workflowExecutionId := ""
	if _, ok := args["workflowExecutionId"]; ok {
		workflowExecutionId = args["workflowExecutionId"].(string)
	}

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
		if workflowExecutionId != "" {
			err = streamClient.Send(workflowExecutionId, resp.Chunk)
		}
		if err != nil {
			return err
		}
		fmt.Println(resp.Chunk)
		time.Sleep(100 * time.Millisecond)
	}
}

func (g *ShellRPCClient) Sh(ctx context.Context, args map[string]interface{}, streamClient *stream.GrpcClient) error {
	cmd := args["text"].(string)
	containerId := ""
	if _, ok := args["containerId"]; ok {
		containerId = args["containerId"].(string)
	}
	workflowExecutionId := ""
	if _, ok := args["workflowExecutionId"]; ok {
		workflowExecutionId = args["workflowExecutionId"].(string)
	}

	serverStream, err := g.client.Sh(ctx, &pb.ShellRequest{
		Command:     cmd,
		ContainerId: containerId,
	})
	if err != nil {
		return err
	}

	// Create a channel for event transfer
	eventChannel := make(chan *pb.ShellResponse, 100)
	var eg errgroup.Group

	// Goroutine to receive events from server streaming
	eg.Go(func() error {
		defer close(eventChannel) // Close when server stream ends

		for {
			event, err := serverStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error receiving event: %v", err)
				return err
			}
			eventChannel <- event
		}
		return nil
	})

	 // Goroutine to send events to client streaming
	 eg.Go(func() error {
		for event := range eventChannel {
			if err := streamClient.Send(workflowExecutionId, event.Chunk); err != nil {
				log.Printf("Error sending event: %v", err)
				return err
			}
		}

		_, err = streamClient.CloseStream()
		return err
	})

	return eg.Wait()
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
