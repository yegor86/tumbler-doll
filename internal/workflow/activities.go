package workflow

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/yegor86/tumbler-doll/plugins"
)

type StageActivities struct {
}

func (a *StageActivities) StageActivity(ctx context.Context, steps []*Step, agent Agent) ([]string, error) {
	// name := activity.GetInfo(ctx).ActivityType.Name
	// fmt.Printf("Run %s...\n", name)
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
	
	pluginManager := plugins.GetInstance()
	for _, step := range steps {
		command, params := step.ToCommand()
		if dockerContainer != nil {
			params["containerId"] = dockerContainer.ContainerId
		}
		
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