package shared

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
)

type ContainerId string

type DockerClient interface {
	Pull(ctx context.Context, imageName string) (io.ReadCloser, error)
	RunContainer(ctx context.Context, imageName string) (ContainerId, error)
	ExecContainer(ctx context.Context, containerId string, cmd []string) (*types.HijackedResponse, error)
	StopContainer(ctx context.Context, containerId ContainerId) error
	Stop() error
}

type DockerClientImpl struct {
	docker *dockerClient.Client
}

var (
	dockerClientVersion string = "1.46"
)

func NewDockerClient(ctx context.Context) (DockerClient, error) {

	// Create a Docker client
	docker, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithVersion(dockerClientVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &DockerClientImpl{
		docker: docker,
	}, nil
}

func (p *DockerClientImpl) Pull(ctx context.Context, imageName string) (io.ReadCloser, error) {
	return p.docker.ImagePull(ctx, buildImageWithTag(imageName), image.PullOptions{})
}

// RunContainer: same as `docker run`
func (p *DockerClientImpl) RunContainer(ctx context.Context, imageName string) (ContainerId, error) {

	// Create the container
	resp, err := p.docker.ContainerCreate(ctx, &container.Config{
		Image:      imageName,
		Entrypoint: []string{"sh"},
		Tty:        true,
	}, &container.HostConfig{}, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return ContainerId(""), fmt.Errorf("failed to create Docker container: %w", err)
	}

	// Start the container
	if err := p.docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return ContainerId(""), fmt.Errorf("failed to start Docker container: %w", err)
	}

	return ContainerId(resp.ID), nil
}

// ExecContainer: same as `docker exec`
func (p *DockerClientImpl) ExecContainer(ctx context.Context, containerId string, cmd []string) (*types.HijackedResponse, error) {
	docker, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithVersion(dockerClientVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer docker.Close()
	
	execResp, err := docker.ContainerExecCreate(ctx, containerId, container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create exec instance for command '%v': %w", cmd, err)
	}

	// Start the command execution
	execAttachResp, err := docker.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to start exec instance for command '%v': %w", cmd, err)
	}
	return &execAttachResp, nil
}

func (p *DockerClientImpl) StopContainer(ctx context.Context, containerId ContainerId) error {
	// Stop and remove the container after all commands are executed
	if err := p.docker.ContainerStop(ctx, string(containerId), container.StopOptions{}); err != nil {
		return err
	}
	if err := p.docker.ContainerRemove(ctx, string(containerId), container.RemoveOptions{}); err != nil {
		return err
	}
	return nil
}

func (p *DockerClientImpl) Stop() error {
	return p.docker.Close()
}

// Append tag 'latest' to image without tag
func buildImageWithTag(imageName string) string {
	imageTag := strings.Split(imageName, ":")
	if len(imageTag) > 1 {
		return imageTag[0] + ":" + imageTag[1]
	}
	return imageTag[0] + ":latest"
}
