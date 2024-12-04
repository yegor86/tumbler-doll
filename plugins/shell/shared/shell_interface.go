package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Shell interface {
	Echo(args map[string]interface{}) string
	Sh(args map[string]interface{}) string
}

// Here is an implementation that talks over RPC
type ShellRPCClient struct{ client *rpc.Client }

func (g *ShellRPCClient) Echo(args map[string]interface{}) string {
	var resp string
	err := g.client.Call("Plugin.Echo", args, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

func (g *ShellRPCClient) Sh(args map[string]interface{}) string {
	var resp string
	err := g.client.Call("Plugin.Sh", args, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

type ShellRPCServer struct {
	// This is the real implementation
	Impl Shell
}

func (s *ShellRPCServer) Echo(args map[string]interface{}, resp *string) error {
	*resp = s.Impl.Echo(args)
	return nil
}

func (s *ShellRPCServer) Sh(args map[string]interface{}, resp *string) error {
	*resp = s.Impl.Sh(args)
	return nil
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
