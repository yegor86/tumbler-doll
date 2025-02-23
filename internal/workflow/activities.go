package workflow

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/yegor86/tumbler-doll/plugins"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

type StageActivities struct {
}

func (a *StageActivities) StageActivity(ctx context.Context, steps []*Step, agent Agent) ([]string, error) {
	var results []string
	var dockerContainer *DockerContainer = nil
	var err error

	if agent.Docker != nil && agent.Docker.Image != "" {
		dockerContainer, err = NewDockerContainer(ctx, string(agent.Docker.Image))
		if err != nil {
			return results, err
		}
		defer dockerContainer.StopContainer(ctx, string(agent.Docker.Image))
	}

	// Get workflow information to send signals
	info := activity.GetInfo(ctx)

	pluginManager := plugins.GetInstance()
	for _, step := range steps {
		command, params := step.ToCommand()

		params["workflowExecutionId"] = info.WorkflowExecution.ID
		if dockerContainer != nil {
			params["containerId"] = dockerContainer.ContainerId
		}

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
			// return results, nil
			// return results, temporal.NewNonRetryableApplicationError(
			// 	"command execution failed",
			// 	"plugin",
			// 	err,
			// )
		} else if invokeResult, ok := output.(string); ok {
			results = append(results, invokeResult)
		} else {
			results = append(results, "(empty)")
		}
	}

	return results, nil
}
