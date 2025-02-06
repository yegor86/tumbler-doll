package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yegor86/tumbler-doll/internal/jenkins/jobs"
	"github.com/yegor86/tumbler-doll/internal/workflow"
	temporal "go.temporal.io/sdk/client"
)

var (
	dslParser workflow.DslParser
)

// Handler function for POST /submit/{jobpath}
func SubmitJob(wfClient temporal.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		jobDB := jobs.GetInstance()

		jobPath := chi.URLParam(r, "*")
		job := jobDB.FindJobs(jobPath)
		if job == nil {
			http.Error(w, fmt.Sprintf("job %s not found", jobPath), http.StatusNotFound)
			return
		}

		pipeline, err := dslParser.Parse(job.Script)
		if err != nil {
			http.Error(w, "Error reading script", http.StatusInternalServerError)
			log.Printf("Error reading script %s", jobPath)
			return
		}

		context := r.Context()
		workflowOptions := temporal.StartWorkflowOptions{
			ID:        job.Name + "/" + uuid.New().String(),
			TaskQueue: "JobQueue",
		}
		we, err := wfClient.ExecuteWorkflow(context, workflowOptions, workflow.GroovyDSLWorkflow, *pipeline)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			http.Error(w, "Error executing workflow", http.StatusInternalServerError)
			log.Printf("Unable to execute workflow %v", err)
			return
		}

		fmt.Fprintf(w, "Started workflow: WorkflowID=%s, RunID=%s", we.GetID(), we.GetRunID())
	}
}
