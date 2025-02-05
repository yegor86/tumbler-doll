package main

import (
	"io"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/yegor86/tumbler-doll/plugins/shell/proto"
	"google.golang.org/grpc"
)

type DummyResponse struct {
	grpc.ServerStream
}

func (r *DummyResponse) Send(resp *proto.LogResponse) error {
	return nil
}

func Test_shell_command(t *testing.T) {

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Debug,
		Output:     os.Stdout,
		JSONFormat: true,
	})

	shellImpl := &ShellPluginImpl{
		logger: logger,
	}

	err := shellImpl.Sh(&proto.LogRequest{
		Command:     "mvn --version",
		ContainerId: "",
	}, &DummyResponse{})

	if err != nil && err != io.EOF {
		t.Fatalf("Error executing plugin: %v", err)
	}

	err = shellImpl.Sh(&proto.LogRequest{
		Command:     "java -version",
		ContainerId: "",
	}, &DummyResponse{})

	if err != nil && err != io.EOF {
		t.Fatalf("Error executing plugin: %v", err)
	}
}
