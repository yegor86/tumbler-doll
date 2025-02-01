package shared

import (
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// StreamLogsReply is the response struct (log line chunk)
type StreamLogsReply struct {
	Chunk string
}

type Shell interface {
	Echo(args map[string]interface{}, reply *StreamLogsReply) error
	Sh(args map[string]interface{}, reply *StreamLogsReply) error
}

// Here is an implementation that talks over RPC
type ShellRPCClient struct{ client *rpc.Client }

func (g *ShellRPCClient) Echo(args map[string]interface{}, reply *StreamLogsReply) error {
	// err := g.client.Call("Plugin.Echo", args, reply)
	var err error
	for {
		err = g.client.Call("Plugin.Echo", args, reply)
		if err != nil {
			break
		}
		fmt.Print(reply.Chunk) // Print log chunks in real time
	}
	return err
}

func (g *ShellRPCClient) Sh(args map[string]interface{}, reply *StreamLogsReply) error {
	// err := g.client.Call("Plugin.Sh", args, reply)
	var err error
	for {
		err = g.client.Call("Plugin.Sh", args, reply)
		if err != nil {
			break
		}
		fmt.Print(reply.Chunk) // Print log chunks in real time
	}
	return err
}

type ShellRPCServer struct {
	// This is the real implementation
	Impl Shell
}

func (s *ShellRPCServer) Echo(args map[string]interface{}, reply *StreamLogsReply) error {
	return s.Impl.Echo(args, reply)
}

func (s *ShellRPCServer) Sh(args map[string]interface{}, reply *StreamLogsReply) error {
	return s.Impl.Sh(args, reply)
}

type ShellPlugin struct {
	// Impl Injection
	Impl Shell
}

func (p *ShellPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ShellRPCServer{Impl: p.Impl}, nil
}

func (ShellPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ShellRPCClient{client: c}, nil
}
