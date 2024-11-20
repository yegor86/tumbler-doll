package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Scm interface {
	Checkout(url string, branch string, credentialsId string) string
}

// Here is an implementation that talks over RPC
type ScmRPC struct{ client *rpc.Client }

func (g *ScmRPC) Checkout(url string, branch string, credentialsId string) string {
	var resp string
	err := g.client.Call("Plugin.Checkout", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

type ScmRPCServer struct {
	// This is the real implementation
	Impl Scm
}

func (s *ScmRPCServer) Checkout(url string, branch string, credentialsId string, resp *string) error {
	*resp = s.Impl.Checkout(url, branch, credentialsId)
	return nil
}

type ScmPlugin struct {
	// Impl Injection
	Impl Scm
}

func (p *ScmPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ScmRPCServer{Impl: p.Impl}, nil
}

func (ScmPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ScmRPC{client: c}, nil
}