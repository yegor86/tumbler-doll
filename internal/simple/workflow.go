package simple

import (
	"strings"
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	// Workflow is the type used to express the workflow definition. Variables are a map of valuables. Variables can be
	// used as input to Activity.
	Workflow struct {
		Variables map[string]string
		Root      Statement
	}

	// Statement is the building block of dsl workflow. A Statement can be a simple ActivityInvocation or it
	// could be a Sequence or Parallel.
	Statement struct {
		Activity *ActivityInvocation
		Sequence *Sequence
		Parallel *Parallel
	}

	// Sequence consist of a collection of Statements that runs in sequential.
	Sequence struct {
		Elements []*Statement
	}

	// Parallel can be a collection of Statements that runs in parallel.
	Parallel struct {
		Branches []*Statement
	}

	// ActivityInvocation is used to express invoking an Activity. The Arguments defined expected arguments as input to
	// the Activity, the result specify the name of variable that it will store the result as which can then be used as
	// arguments to subsequent ActivityInvocation.
	ActivityInvocation struct {
		Name           string
		Arguments      []string
		Result         string
		Commands       []string
		ContainerImage string `yaml:"container_image"`
	}

	executable interface {
		execute(ctx workflow.Context, variables map[string]string, results map[string]any) error
	}
)

// SimpleDSLWorkflow workflow definition
func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow Workflow) ([]byte, error) {
	variables := make(map[string]string)
	//workflowcheck:ignore Only iterates for building another map
	for k, v := range dslWorkflow.Variables {
		variables[k] = v
	}
	results := make(map[string]any)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	err := dslWorkflow.Root.execute(ctx, variables, results)
	if err != nil {
		logger.Error("DSL Workflow failed.", "Error", err)
		return nil, err
	}

	logger.Info("DSL Workflow completed.")
	return nil, err
}

func (b *Statement) execute(ctx workflow.Context, variables map[string]string, results map[string]any) error {
	if b.Parallel != nil {
		err := b.Parallel.execute(ctx, variables, results)
		if err != nil {
			return err
		}
	}
	if b.Sequence != nil {
		err := b.Sequence.execute(ctx, variables, results)
		if err != nil {
			return err
		}
	}
	if b.Activity != nil {
		err := b.Activity.execute(ctx, variables, results)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a ActivityInvocation) execute(ctx workflow.Context, variables map[string]string, results map[string]any) error {
	inputParam := makeInput(a.Commands, a.Arguments, variables)
	var result []string

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5, // Adjust the timeout as needed
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, a.Name, inputParam, a.ContainerImage).Get(ctx, &result)
	if err != nil {
		return err
	}
	if a.Result != "" {
		results[a.Result] = result
	}
	return nil
}

func (s Sequence) execute(ctx workflow.Context, variables map[string]string, results map[string]any) error {
	for _, a := range s.Elements {
		err := a.execute(ctx, variables, results)
		if err != nil {
			return err
		}
	}
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
	for _, s := range p.Branches {
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

	for i := 0; i < len(p.Branches); i++ {
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

func makeInput(commands []string, argNames []string, argsMap map[string]string) []string {
	var results []string
	for _, command := range commands {
		for _, arg := range argNames {
			command = strings.ReplaceAll(command, "$"+arg, argsMap[arg])
		}
		results = append(results, command)
	}
	return results
}
