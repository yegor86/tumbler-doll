package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	temporal "go.temporal.io/sdk/client"

	"github.com/yegor86/tumbler-doll/internal/api/sse"
	"github.com/yegor86/tumbler-doll/internal/workflow"
)

func StreamLogs(wfClient temporal.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set headers for SSE
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		jobPath := chi.URLParam(r, "*")
		workflowId := r.URL.Query().Get("workflowId")
		jobId := workflowId[strings.LastIndex(workflowId, "/")+1:]
		state := workflow.Undefined

		inFilepath := filepath.Join(os.Getenv("JENKINS_HOME"), jobPath, "builds", jobId, "log")
		err := os.ErrNotExist
		for err != nil && state != workflow.Done {
			_, err = os.Stat(inFilepath)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}

		inFile, err := os.Open(inFilepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer inFile.Close()

		// Remember inFile read-offset to use it in case of inFile was not fully read at the first traversal
		var seekOffset int64 = 0

		for state != workflow.Done {
			_, err = inFile.Seek(seekOffset, io.SeekStart)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = sse.CopyAndFlush(w, inFile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			seekOffset, err = inFile.Seek(0, io.SeekCurrent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			state, err = workflow.GetState(wfClient, workflowId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}

		fmt.Fprintf(w, "Completed job: WorkflowID=%s", jobId)
	}
}
