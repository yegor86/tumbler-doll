package grpc

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	pb "github.com/yegor86/tumbler-doll/internal/grpc/proto"
)

type GrpcServer struct {
	Addr net.Addr
	
	server *grpc.Server
	pb.UnimplementedLogStreamingServiceServer
	onReceived func(logEvent *pb.LogRequest)
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
		logEvent, err := stream.Recv()
		
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
			s.onReceived(logEvent)
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (s *GrpcServer) ListenAndServe(onReceived func(req *pb.LogRequest)) error {
	s.onReceived = onReceived
	return s.ListenAndServeWithHostPort(":50051")
}

func (s *GrpcServer) ListenAndServeWithHostPort(hostPort string) error {
	// Start gRPC server
	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		return err
	}
	s.Addr = listener.Addr()

	log.Printf("gRPC server is running on host port %s...\n", hostPort)
	if err := s.server.Serve(listener); err != nil {
		return err
	}
	log.Print("gRPC server is stopping...\n")
	return nil
}

func (s *GrpcServer) Shutdown(ctx context.Context) error {
	s.server.Stop()
	return nil
}
