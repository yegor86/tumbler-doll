package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yegor86/tumbler-doll/internal/jenkins/jobs"
	"github.com/yegor86/tumbler-doll/internal/workflow"
	temporal "go.temporal.io/sdk/client"
)

var (
	dslParser workflow.DslParser
)

type SubmitJobResponse struct {
	Status     string
	WorkflowID string
	RunId      string
}

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

		jobId := uuid.New().String()
		workflowOptions := temporal.StartWorkflowOptions{
			ID:        job.Name + "/" + jobId,
			TaskQueue: "JobQueue",
		}
		props := map[string]interface{}{
			"jobPath": jobPath,
			"jobId":   jobId,
			"ipAddress": os.Getenv("IP_ADDRESS"),
		}
		we, err := wfClient.ExecuteWorkflow(context.Background(), workflowOptions, workflow.GroovyDSLWorkflow, *pipeline, props)

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			http.Error(w, "Error executing workflow", http.StatusInternalServerError)
			log.Printf("Unable to execute workflow %v", err)
			return
		}

		if err := json.NewEncoder(w).Encode(SubmitJobResponse{
			Status:     "Started workflow: WorkflowID=%s, RunID=%s",
			WorkflowID: we.GetID(),
			RunId:      we.GetRunID(),
		}); err != nil {
			http.Error(w, "Failed to encode jobs as JSON", http.StatusInternalServerError)
		}
	}
}
