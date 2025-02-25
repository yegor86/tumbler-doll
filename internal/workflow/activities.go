package workflow

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/yegor86/tumbler-doll/plugins"
	"github.com/yegor86/tumbler-doll/plugins/docker"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

type StageActivities struct {
}

func (a *StageActivities) StageActivity(ctx context.Context, steps []*Step, agent Agent) ([]string, error) {
	var results []string

	// Get workflow information
	info := activity.GetInfo(ctx)
	ctx = context.WithValue(ctx, "workflowExecutionId", info.WorkflowExecution.ID)

	pluginManager := plugins.GetInstance()
	
	dockerPlugin, found := pluginManager.FindPlugin("docker").(*docker.DockerPlugin)
	if agent.Docker != nil && agent.Docker.Image != "" && found {
		ctx = context.WithValue(ctx, "imageName", string(agent.Docker.Image))
		err := dockerPlugin.Pull(ctx)
		if err != nil {
			log.Printf("Docker pull image failed: %v\n", err)
			return results, err
		}
		containerId, err := dockerPlugin.RunContainer(ctx)
		if err != nil {
			log.Printf("Docker run container failed: %v\n", err)
			return results, err
		}
		ctx = context.WithValue(ctx, "containerId", containerId)
		
		defer dockerPlugin.StopContainer(ctx, containerId)
	}

	for _, step := range steps {
		command, params := step.ToCommand()
		params["workflowExecutionId"] = info.WorkflowExecution.ID
		params["containerId"] = ctx.Value("containerId")

		pluginName, methodFunc, ok := pluginManager.GetPluginInfo(command)
		if !ok {
			err := temporal.NewNonRetryableApplicationError(
				"unexpected error",
				"plugin",
				fmt.Errorf("plugin is not registered for the command %s", command),
			)
			return nil, err
		}

		capitalizedCommand := strings.ToUpper(methodFunc[:1]) + strings.ToLower(methodFunc[1:])
		output, err := pluginManager.Execute(
			pluginName,
			capitalizedCommand,
			params)

		if err != nil {
			log.Printf("Command execution failed: %s", err)
			results = append(results, err.Error())
			return results, temporal.NewNonRetryableApplicationError(
				"command execution failed",
				"plugin",
				err,
			)
		} else if invokeResult, ok := output.(string); ok {
			results = append(results, invokeResult)
		} else {
			results = append(results, "(empty)")
		}
	}

	return results, nil
}
