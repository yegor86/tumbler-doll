package docker

import (
	"context"
	"errors"
	"fmt"

	"github.com/yegor86/tumbler-doll/internal/grpc"
	"github.com/yegor86/tumbler-doll/plugins"

	logstream "github.com/yegor86/tumbler-doll/internal/grpc/proto"
	"github.com/yegor86/tumbler-doll/plugins/docker/shared"
)

type DockerPlugin struct {
	dockerClient shared.DockerClient
	streamClient *grpc.GrpcClient
	ctx context.Context
}

func (p *DockerPlugin) Start(ctx context.Context) error {
	temporalHostPort, ok := ctx.Value("temporalHostport").(string)
	if !ok {
		return fmt.Errorf("failed to extract TEMPORAL_ADDRESS from context: %v", ctx.Value("temporalHostport"))
	}

	dockerClient, err := shared.NewDockerClient(ctx)
	if err != nil {
		return err
	}

	streamClient, err := grpc.NewClient(temporalHostPort)
	if err != nil {
		return err
	}
	
	p.dockerClient = dockerClient
	p.streamClient = streamClient
	p.ctx = ctx
	return nil
}

func (p *DockerPlugin) Stop() error {
	err := p.streamClient.CloseStream()
	if err != nil {
		return err
	}
	return p.dockerClient.Stop()
}

func (p *DockerPlugin) ListMethods() map[string]string {
	return map[string]string{}
}

func (p *DockerPlugin) Pull(ctx context.Context) error {
	workflowExecutionId, ok := ctx.Value("workflowExecutionId").(string)
	if !ok {
		return errors.New("unable to redirect DockerPlugin.Pull output. 'workflowExecutionId' not found")
	}

	imageName, ok := ctx.Value("imageName").(string)
	if !ok {
		return fmt.Errorf("docker image type is wrong %v", imageName)
	}

	ioReader, err := p.dockerClient.Pull(p.ctx, imageName)
	if err != nil {
		return err
	}
	return plugins.RedirectIoReaderToGrpc(ioReader, p.streamClient.Stream, func(resp string) *logstream.LogRequest {
		return &logstream.LogRequest{
			Message: resp,
			WorkflowId: workflowExecutionId,
		}
	})
}

func (p *DockerPlugin) RunContainer(ctx context.Context) (string, error) {
	imageName, ok := ctx.Value("imageName").(string)
	if !ok {
		return "", fmt.Errorf("docker image type is wrong %v", imageName)
	}
	return p.dockerClient.RunContainer(ctx, imageName)
}

func (p *DockerPlugin) StopContainer(ctx context.Context, containerId string) error {
	return p.dockerClient.StopContainer(ctx, containerId)
}
