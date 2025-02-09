package shared

import (
	"errors"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Scm interface {
	Checkout(args map[string]interface{}) (string, error)
}

// Here is an implementation that talks over RPC
type ScmRPCClient struct{ client *rpc.Client }

func (g *ScmRPCClient) Checkout(args map[string]interface{}) (string, error) {
	var resp []string
	err := g.client.Call("Plugin.Checkout", args, &resp)
	if err != nil{
		return "", err
	}
	result := resp[0]
	errMessage := resp[1]
	if errMessage != "" {
		return "", errors.New(errMessage)
	}

	return result, nil
}

type ScmRPCServer struct {
	// This is the real implementation
	Impl Scm
}

func (s *ScmRPCServer) Checkout(args map[string]interface{}, resp *[]string) error {
	r, err := s.Impl.Checkout(args)
	
	var errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	*resp = []string{r, errMessage}
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
