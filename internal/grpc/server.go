package grpc

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/yegor86/tumbler-doll/internal/grpc/proto"
	"google.golang.org/grpc"
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

// pipeLogEventToFile: pipe log event into a text file
func (s *GrpcServer) pipeLogEventToFile(req *pb.LogRequest) {
			
	workflowId, chunk := req.WorkflowId, req.Message
	
	delim := strings.LastIndex(workflowId, "/")
	jobPath, jobId := workflowId[:delim], workflowId[delim + 1:]
	opath := filepath.Join(os.Getenv("JENKINS_HOME"), jobPath, "builds", jobId)
	err := os.MkdirAll(opath, 0740)
	if err != nil {
		log.Printf("error creating dir %s: %v", opath, err)
		return
	}

	ofile, err := os.OpenFile(filepath.Join(opath, "log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error creating/opening file %s: %v", filepath.Join(opath, "log"), err)
		return
	}

	w := bufio.NewWriter(ofile)

	// write a chunk
	if _, err := w.Write([]byte(chunk + "\n")); err != nil {
		log.Printf("error when writing log %v. Failed chunk: %s", err, chunk)
	}
	if err = w.Flush(); err != nil {
		log.Printf("error when flushing log %v. Failed chunk: %s", err, chunk)
	}
}

func (s *GrpcServer) ListenAndServe() error {
	return s.ListenAndServeWithCallback(nil)
}

func (s *GrpcServer) ListenAndServeWithCallback(onReceived func(req *pb.LogRequest)) error {
	s.onReceived = onReceived
	if onReceived == nil {
		s.onReceived = s.pipeLogEventToFile
	}
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
