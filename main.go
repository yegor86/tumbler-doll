package main

import (
	"fmt"

	"github.com/yegor86/tumbler-doll/cmd"
	"github.com/yegor86/tumbler-doll/plugins/scm/shared"
	"github.com/yegor86/tumbler-doll/plugins/scm"
	"github.com/yegor86/tumbler-doll/plugins"
)

func main() {

	pluginManager := plugins.NewPluginManager()
	defer pluginManager.UnregisterAll()

	pluginManager.Register("scm", &scm.ScmPlugin{})

	args := shared.CheckoutArgs{
		Url:           "http://testurl.com",
		Branch:        "master",
		CredentialsId: "",
	}
	res, err := pluginManager.Execute("scm", "Checkout", args)
	if err != nil {
		fmt.Printf("Error executing scm.Checkout %v", err)
	}
	fmt.Println(res)

	cmd.Execute()
}
