package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type CheckoutArgs struct {
	Url           string
	Branch        string
	CredentialsId string
}

type Scm interface {
	Checkout(args CheckoutArgs) string
}

// Here is an implementation that talks over RPC
type ScmRPCClient struct{ client *rpc.Client }

func (g *ScmRPCClient) Checkout(args CheckoutArgs) string {
	var resp string
	err := g.client.Call("Plugin.Checkout", args, &resp)
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

func (s *ScmRPCServer) Checkout(args CheckoutArgs, resp *string) error {
	*resp = s.Impl.Checkout(args)
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
	return &ScmRPCClient{client: c}, nil
}
