package plugins

import (
	"io"

	"google.golang.org/grpc"
)

func Redirect[EventResp, LogReq, LogResp any](out grpc.ServerStreamingClient[EventResp], in grpc.ClientStreamingClient[LogReq, LogResp], toReq func(r *EventResp) *LogReq) error {

	for {
		event, err :=  out.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := in.Send(toReq(event)); err != nil {
			return err
		}
	}
	return nil
}
