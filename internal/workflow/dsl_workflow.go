package workflow

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	// Pipeline represents the main Jenkins pipeline structure
	Pipeline struct {
		Agent  *Agent   `"pipeline" "{" "agent" @@`
		Stages []*Stage `"stages" "{" @@+ "}"`
		Close  string   `"}"`
	}

	// Agent represents the agent block in a Jenkinsfile
	Agent struct {
		None   bool    `( "none" )?`
		Docker *Docker `( "{" "docker" @@ "}" )?`
	}

	Docker struct {
		Image QuotedString `@String`
	}

	Parallel []*Stage

	// Stage represents a stage block within stages
	Stage struct {
		Name     QuotedString `"stage" "(" @String ")" "{"`
		Agent    *Agent       `( "agent" @@ )?`
		Steps    []*Step      `( "steps" "{" @@+ "}" )?`
		FailFast *bool        `( "failFast" @Bool )?`
		Parallel Parallel     `( "parallel" "{" @@+ "}" )?`
		Close    string       `"}"`
	}

	// Step represents individual steps within a stage
	Step struct {
		SingleKV *SingleKVCommand `@@ |`
		MultiKV  *MultiKVCommand  `@@`
	}

	SingleKVCommand struct {
		Command string       `@Ident`
		Value   QuotedString `@String`
	}

	MultiKVCommand struct {
		Command string  `@Ident`
		Params  []Param `@@ ("," @@)*`
	}

	Param struct {
		Key   string       `@Ident ":"`
		Value QuotedString `@String`
	}

	executable interface {
		execute(ctx workflow.Context, variables map[string]string, results map[string]any) error
	}
)

type QuotedString string

// Capture method strips quotes from the Image field
func (o *QuotedString) Capture(values []string) error {
	*o = QuotedString(strings.Trim(values[0], `"'`))
	return nil
}

func GroovyDSLWorkflow(ctx workflow.Context, pipeline Pipeline, properties map[string]interface{}) (map[string]any, error) {
	logger := workflow.GetLogger(ctx)

	if err := dumpLogs(ctx, properties); err != nil {
		return nil, err
	}

	variables := make(map[string]string)
	results := make(map[string]any)
	for _, stage := range pipeline.Stages {
		if err := stage.execute(ctx, variables, results); err != nil {
			return nil, err
		}
	}

	logger.Info("Groovy Workflow completed.")
	return results, nil
}

func dumpLogs(ctx workflow.Context, props map[string]interface{}) error {
	
	jobPath, ok := props["jobPath"]
	if !ok {
		return errors.New("'jobPath' not found inside properties")
	}
	jobId, ok := props["jobId"]
	if !ok {
		return errors.New("'jobId' not found in workflow context")
	}

	outPath := filepath.Join(os.Getenv("JENKINS_HOME"), jobPath.(string), "builds", jobId.(string))
	err := os.MkdirAll(outPath, 0740)
	if err != nil {
		return err
	}

	outFile, err := os.Create(filepath.Join(outPath, "log"))
	if err != nil {
		return err
	}
    w := bufio.NewWriter(outFile)
	
	temporalChan := workflow.GetSignalChannel(ctx, "logs")
	workflow.Go(ctx, func(ctx workflow.Context) {
		
		for {
			var logChunk string
			if more := temporalChan.Receive(ctx, &logChunk); !more {
				break
			}
			// write a chunk
			if _, err := w.Write([]byte(logChunk)); err != nil {
				log.Printf("error when writing log %v. Failed chunk: %s", err, logChunk)
			}
			if err = w.Flush(); err != nil {
				log.Printf("error when flushing log %v. Failed chunk: %s", err, logChunk)
			}
		}
		if err := outFile.Close(); err != nil {
            panic(err)
        }
	})
	return nil
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
