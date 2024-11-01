package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"
	"github.com/yegor86/tumbler-doll/internal/dsl"
)

func GroovyDSLWorkflow(ctx workflow.Context, pipeline dsl.Pipeline) ([]byte, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)


	logger.Info("Grrovy Workflow completed.")
	return nil, nil
}