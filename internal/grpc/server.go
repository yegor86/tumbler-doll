package grpc

import (
	"io"
	"log"
	"net"
	"time"

	pb "github.com/yegor86/tumbler-doll/internal/grpc/proto"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	server *grpc.Server
	pb.UnimplementedLogStreamingServiceServer
	onReceived func(workflowId string, msg string)
}

func NewServer() *GrpcServer {
	grpcServer := grpc.NewServer()

	server := &GrpcServer{
		server: grpcServer,
	}
	pb.RegisterLogStreamingServiceServer(grpcServer, server)

	return server
}

func (s *GrpcServer) Stream(stream grpc.ClientStreamingServer[pb.LogRequest, pb.LogResponse]) error {

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			resp := &pb.LogResponse{
				Status: "Done", // Client finished streaming
			}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			return err
		}
		if s.onReceived != nil {
			s.onReceived(msg.WorkflowId, msg.Message)
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (s *GrpcServer) ListenAndServe(onReceived func(workflowId string, msg string)) error {
	s.onReceived = onReceived
	return s.ListenAndServeWithHostPort(":50051")
}

func (s *GrpcServer) ListenAndServeWithHostPort(hostPort string) error {
	// Start gRPC server
	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		return err
	}

	log.Printf("gRPC server is running on host port %s...\n", hostPort)
	if err := s.server.Serve(listener); err != nil {
		return err
	}
	log.Print("gRPC server is stopping...\n")
	return nil
}
