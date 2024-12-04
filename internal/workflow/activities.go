package workflow

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
	"github.com/yegor86/tumbler-doll/plugins"
)

type StageActivities struct {
}

func (a *StageActivities) StageActivity(ctx context.Context, steps []*Step, agent Agent) ([]string, error) {
	// name := activity.GetInfo(ctx).ActivityType.Name
	// fmt.Printf("Run %s...\n", name)
	
	if agent.Docker != nil && agent.Docker.Image != "" {
		return DockerActivity(ctx, string(agent.Docker.Image), steps)
	}
	
	pluginManager := plugins.GetInstance()
	var results []string
	for _, step := range steps {
		command, params := step.ToCommand()
		
		pluginName := pluginManager.GetPluginName(command)
		methodFunc := pluginManager.GetFunctionByMethod(command)
		
		capitalizedCommand := strings.ToUpper(methodFunc[:1]) + strings.ToLower(methodFunc[1:])
		output, err := pluginManager.Execute(
			pluginName,
			capitalizedCommand,
			params)
		
		if err != nil {
			log.Printf("Command execution failed: %s", err)
			results = append(results, err.Error())
			return results, err
		}
		results = append(results, output.(string))
		fmt.Printf("Command Output: %s\n", output)
	}

	return results, nil
}

// DockerActivity starts a Docker container with a specified image
func DockerActivity(ctx context.Context, imageName string, commands []*Step) ([]string, error) {

	// Create a Docker client
	docker, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithVersion("1.46"))
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer docker.Close()

	// Pull the image (if not present locally)
	pullOut, err := docker.ImagePull(ctx, buildImageWithTag(imageName), image.PullOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to pull Docker image: %w", err)
	}
	defer pullOut.Close()
	// todo: Stream docker pull output on to the progress view
	io.Copy(io.Discard, pullOut)

	// Create the container
	resp, err := docker.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		// Entrypoint: []string{"sh"},
		Tty:   true,
	}, &container.HostConfig{}, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker container: %w", err)
	}

	containerID := resp.ID

	// Start the container
	if err := docker.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start Docker container: %w", err)
	}

	// Execute each command separately inside the container
	results := []string{}
	for _, cmd := range commands {
		cmd_, _ := cmd.ToCommand()
		execResp, err := docker.ContainerExecCreate(ctx, containerID, container.ExecOptions{
			Cmd:          []string{cmd_},
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
		defer execAttachResp.Close()
		execAttachResp.CloseWrite()

		// Capture the output
		output, err := io.ReadAll(execAttachResp.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read output for command '%v': %w", cmd, err)
		}
		results = append(results, removeControlChars(output))
	}

	// Stop and remove the container after all commands are executed
	if err := docker.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		return nil, fmt.Errorf("failed to stop container: %w", err)
	}
	if err := docker.ContainerRemove(ctx, containerID, container.RemoveOptions{}); err != nil {
		log.Printf("Failed to remove container: %v", err)
	}

	return results, nil
}

// RemoveControlChars removes non-printable ASCII characters from byte array and return human readble string.
func removeControlChars(input []byte) string {
    return string(bytes.Map(func(r rune) rune {
        if unicode.IsControl(r) {
            return -1
        }
        return r
    }, input))
}

func buildImageWithTag(imageName string) string {
	imageTag := strings.Split(imageName, ":")
	if len(imageTag) > 1 {
		return imageTag[0] + ":" + imageTag[1]
	}
	return imageTag[0] + ":" + "latest"
}