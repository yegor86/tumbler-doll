package handler

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/yegor86/tumbler-doll/internal/dsl"
	"github.com/yegor86/tumbler-doll/internal/workflow"
	temporal "go.temporal.io/sdk/client"
)

// uploadForm serves the file upload form from an HTML template file
func UploadForm(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("web", "templates", "upload.html")

	// Parse the template file
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		log.Printf("Template loading error: %v", err)
		return
	}

	// Render the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		log.Printf("Template rendering error: %v", err)
		return
	}
}

// uploadFile handles the file upload
func UploadFile(wfClient temporal.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse the form data
		err := r.ParseMultipartForm(10 << 20) // Limit file size to 10MB
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusInternalServerError)
			return
		}

		// Retrieve the file from data
		file, file_header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to retrieve the file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Read the content of the file into a byte array
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading form data", http.StatusInternalServerError)
			return
		}

		var dslParser dsl.DslParser
		pipeline, err := dslParser.Parse(string(data))
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			log.Printf("Error reading file %v", file_header.Filename)
			return
		}

		workflowOptions := temporal.StartWorkflowOptions{
			ID:        "dsl/" + uuid.New().String(),
			TaskQueue: "dslQueue",
		}
		we, err := wfClient.ExecuteWorkflow(context.Background(), workflowOptions, workflow.GroovyDSLWorkflow, *pipeline)
		if err != nil {
			http.Error(w, "Error executing workflow", http.StatusInternalServerError)
			log.Printf("Unable to execute workflow %v", err)
			return
		}

		// Provide feedback to the user
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
		fmt.Fprintf(w, "Started workflow: WorkflowID=%s, RunID=%s", we.GetID(), we.GetRunID())
	}
}
