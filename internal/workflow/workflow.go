package workflow

import (
	"context"
	"fmt"
	"os"
	"time"

	temporalClient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	Undefined State = -1
	Pending   State = 0
	Started   State = 1
	Running   State = 2
	Done      State = 3
)

type (
	State int64
	
	executable interface {
		execute(ctx workflow.Context, variables map[string]string, results map[string]any) error
	}
)

var (
	currentState State = Undefined
)

func GroovyDSLWorkflow(ctx workflow.Context, pipeline Pipeline, properties map[string]interface{}) (map[string]any, error) {
	currentState = Started

	logger := workflow.GetLogger(ctx)
	// setup query handler for query type "state"
	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (State, error) {
		return currentState, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
		return nil, err
	}

	variables := make(map[string]string)
	results := make(map[string]any)

	fmt.Printf("Temporal address: %s\n", os.Getenv("TEMPORAL_ADDRESS"))
	
	currentState = Running
	for _, stage := range pipeline.Stages {
		if err := stage.execute(ctx, variables, results); err != nil {
			// return nil, err
			logger.Error(err.Error())
			break
		}
	}

	logger.Info("Groovy Workflow completed.")
	currentState = Done
	return results, nil
}

func GetState(wfClient temporalClient.Client, workflowId string) (State, error) {
	msgEncoded, err := wfClient.QueryWorkflow(context.Background(), workflowId, "", "state")
	if err != nil {
		return Undefined, err
	}
	
	var queryResult State
	msgEncoded.Get(&queryResult)
	return queryResult, nil
}

func (stage *Stage) execute(ctx workflow.Context, variables map[string]string, results map[string]any) error {
	if len(stage.Parallel) > 0 {
		parallelResults := make(map[string]any)
		err := stage.Parallel.execute(ctx, variables, parallelResults)
		if err != nil {
			return err
		}
		results[string(stage.Name)] = parallelResults
	}

	if err := stage.executeSteps(ctx, variables, results); err != nil {
		return err
	}
	return nil
}

func (stage *Stage) executeSteps(ctx workflow.Context, variables map[string]string, results map[string]any) error {
	if len(stage.Steps) == 0 {
		return nil
	}
	var result []string

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, "StageActivity", stage.Steps, stage.Agent).Get(ctx, &result)
	if err != nil {
		return err
	}
	results[string(stage.Name)] = result
	return nil
}

func (p Parallel) execute(ctx workflow.Context, variables map[string]string, results map[string]any) error {
	//
	// You can use the context passed in to activity as a way to cancel the activity like standard GO way.
	// Cancelling a parent context will cancel all the derived contexts as well.
	//

	// In the parallel block, we want to execute all of them in parallel and wait for all of them.
	// if one activity fails then we want to cancel all the rest of them as well.
	childCtx, cancelHandler := workflow.WithCancel(ctx)
	selector := workflow.NewSelector(ctx)
	var activityErr error
	for _, s := range p {
		f := executeAsync(s, childCtx, variables, results)
		selector.AddFuture(f, func(f workflow.Future) {
			err := f.Get(ctx, nil)
			if err != nil {
				// cancel all pending activities
				cancelHandler()
				activityErr = err
			}
		})
	}

	for i := 0; i < len(p); i++ {
		selector.Select(ctx) // this will wait for one branch
		if activityErr != nil {
			return activityErr
		}
	}

	return nil
}

func executeAsync(exe executable, ctx workflow.Context, variables map[string]string, results map[string]any) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := exe.execute(ctx, variables, results)
		settable.Set(nil, err)
	})
	return future
}

func (step *Step) Name() string {
	if step.SingleKV != nil {
		return step.SingleKV.Command
	} else if step.MultiKV != nil {
		return step.MultiKV.Command
	}
	return "Unknown"
}

func (step *Step) ToCommand() (string, map[string]interface{}) {
	if step.SingleKV == nil && step.MultiKV == nil {
		return "", nil
	}
	params := make(map[string]interface{})
	if step.SingleKV != nil {
		params["text"] = string(step.SingleKV.Value)
		return step.SingleKV.Command, params
	}
	for _, p := range step.MultiKV.Params {
		params[p.Key] = string(p.Value)
	}
	return step.MultiKV.Command, params
}
