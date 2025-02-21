package grpc

import (
	"context"

	pb "github.com/yegor86/tumbler-doll/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	conn   *grpc.ClientConn
	stream grpc.ClientStreamingClient[pb.LogRequest, pb.LogResponse]
}

// hostPort: localhost:50051
func NewClient(hostPort string) (*GrpcClient, error) {
	conn, err := grpc.NewClient(hostPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewLogStreamingServiceClient(conn)
	stream, err := client.Stream(context.Background(), grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}

	return &GrpcClient {
		conn: conn,
		stream: stream,
	}, nil
}

func (c *GrpcClient) CloseStream() error {
	err := c.stream.CloseSend()
	if err != nil {
		return err
	}
	return c.conn.Close()
}

func (c *GrpcClient) Send(workflowId string, msg string) error {
	return c.stream.Send(&pb.LogRequest{
		WorkflowId: workflowId,
		Message: msg,
	})
}
