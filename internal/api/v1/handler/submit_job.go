package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yegor86/tumbler-doll/internal/dsl"
	"github.com/yegor86/tumbler-doll/internal/jenkins/jobs"
	"github.com/yegor86/tumbler-doll/internal/workflow"
	temporal "go.temporal.io/sdk/client"
)

// Handler function for POST /job/{jobpath}
func SubmitJob(wfClient temporal.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		
		jobDB := jobs.GetInstance()
		
		jobPath := chi.URLParam(r, "*")
		job := jobDB.FindJobs(jobPath)

		var dslParser dsl.DslParser
		pipeline, err := dslParser.Parse(job.Script)
		if err != nil {
			http.Error(w, "Error reading script", http.StatusInternalServerError)
			log.Printf("Error reading script %s", jobPath)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if job == nil {
			http.Error(w, fmt.Sprintf("job %s not found", jobPath), http.StatusNotFound)
			return
		}

		workflowOptions := temporal.StartWorkflowOptions{
			ID:        "job/" + uuid.New().String(),
			TaskQueue: "JobQueue",
		}
		we, err := wfClient.ExecuteWorkflow(context.Background(), workflowOptions, workflow.GroovyDSLWorkflow, *pipeline)
		if err != nil {
			http.Error(w, "Error executing workflow", http.StatusInternalServerError)
			log.Printf("Unable to execute workflow %v", err)
			return
		}

		fmt.Fprintf(w, "Started workflow: WorkflowID=%s, RunID=%s", we.GetID(), we.GetRunID())
	}
}
