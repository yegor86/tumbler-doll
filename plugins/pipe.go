package plugins

import (
	"bufio"
	"io"

	"google.golang.org/grpc"
)

func RedirectGrpcToGrpc[EventResp, LogReq, LogResp any](in grpc.ServerStreamingClient[EventResp], out grpc.ClientStreamingClient[LogReq, LogResp], toReq func(r *EventResp) *LogReq) error {

	for {
		event, err :=  in.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := out.Send(toReq(event)); err != nil {
			return err
		}
	}
	return nil
}

func RedirectIoReaderToGrpc[LogReq, LogResp any](in io.Reader, out grpc.ClientStreamingClient[LogReq, LogResp], toReq func(r string) *LogReq) error {

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		msg := scanner.Text()
		
		if err := out.Send(toReq(msg)); err != nil {
			return err
		}
	}
	return scanner.Err()
}