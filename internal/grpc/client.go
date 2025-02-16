package grpc

import (
	"context"

	pb "github.com/yegor86/tumbler-doll/internal/grpc/proto"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	conn   *grpc.ClientConn
	client pb.LogStreamingServiceClient
	stream grpc.ClientStreamingClient[pb.LogRequest, pb.LogResponse]
}

// hostPort: localhost:50051
func (c *GrpcClient) Connect(hostPort string) error {
	conn, err := grpc.NewClient(hostPort, grpc.EmptyDialOption{})
	if err != nil {
		return err
	}

	stream, err := c.client.Stream(context.Background(), grpc.EmptyCallOption{})
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = pb.NewLogStreamingServiceClient(conn)
	c.stream = stream

	return nil
}

func (c *GrpcClient) CloseStream() (*pb.LogResponse, error) {
	return c.stream.CloseAndRecv()
}

func (c *GrpcClient) Close() error {
	return c.conn.Close()
}

func (c *GrpcClient) Send(msg string) error {
	err := c.stream.Send(&pb.LogRequest{
		Message: msg,
	})

	if err != nil {
		return err
	}

	

	return err
}
