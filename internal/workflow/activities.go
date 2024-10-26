package workflow

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"go.temporal.io/sdk/activity"
)

type SampleActivities struct {
}

func (a *SampleActivities) SampleActivity(ctx context.Context, commands []string) ([]string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Run %s with command %v \n", name, commands)
	
	var results []string
	for _, command := range commands {
		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Command execution failed: %s", err)
			return results, err
		}
		fmt.Printf("Command Output: %s\n", output)
	}
	
	return results, nil
}

