package grpc

import (
	"context"

	pb "github.com/yegor86/tumbler-doll/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	conn   *grpc.ClientConn
	Stream grpc.ClientStreamingClient[pb.LogRequest, pb.LogResponse]
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
		Stream: stream,
	}, nil
}

func (c *GrpcClient) CloseStream() error {
	err := c.Stream.CloseSend()
	if err != nil {
		return err
	}
	return c.conn.Close()
}

func (c *GrpcClient) Send(workflowId string, msg string) error {
	return c.Stream.Send(&pb.LogRequest{
		WorkflowId: workflowId,
		Message: msg,
	})
}
