package grpc

import (
	"io"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func Redirect[Req, Resp any](out grpc.ServerStreamingClient[Resp], in grpc.ClientStreamingClient[Req, Resp], toReq func(r *Resp) *Req) error {

	// Create a channel for event transfer
	eventChannel := make(chan *Resp, 100)
	var eg errgroup.Group

	// Goroutine to receive events from server streaming
	eg.Go(func() error {
		defer close(eventChannel) // Close when server stream ends

		for {
			event, err :=  out.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			eventChannel <- event
		}
		return nil
	})

	// Goroutine to send events to client streaming
	eg.Go(func() error {
		for event := range eventChannel {
			if err := in.Send(toReq(event)); err != nil {
				return err
			}
		}

		_, err := in.CloseAndRecv()
		return err
	})

	return eg.Wait()
}
