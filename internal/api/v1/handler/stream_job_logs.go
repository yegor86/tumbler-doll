package handler

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	temporal "go.temporal.io/sdk/client"

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
		jobId := r.URL.Query().Get("workflowId")

		inPath := filepath.Join(os.Getenv("JENKINS_HOME"), jobPath, "builds", jobId, "log")
		err := os.ErrNotExist
		for err != nil {
			_, err = os.Stat(inPath)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		
		for {
			err1 := openAndRead(w, inPath)
			state, err2 := getState(wfClient, jobId)
			if err1 != nil || err2 != nil {
				http.Error(w, errors.Join(err1, err2).Error(), http.StatusInternalServerError)
				return
			}
			if state == workflow.Done {
				break
			}
		}
		
		fmt.Fprintf(w, "Completed job: WorkflowID=%s", jobId)
	}
}

func openAndRead(w http.ResponseWriter, inPath string) error {
	wFlusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return errors.New("streaming not supported")
	}
	// make a read buffer
	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	
	for scanner.Scan() {
		formatedMsg := fmt.Sprintf("data: %s\n\n", scanner.Text())
		_, err := fmt.Fprint(w, formatedMsg)
		if err != nil {
			log.Println("Error writing to response:", err)
			break
		}
		wFlusher.Flush()
	}
	return scanner.Err()
}

func getState(wfClient temporal.Client, workflowId string) (workflow.State, error) {
	
	var queryResult workflow.State
	msgEncoded, err := wfClient.QueryWorkflow(context.Background(), workflowId, "", "state")
	if err != nil {
		return workflow.Undefined, nil
	}
	msgEncoded.Get(&queryResult)
	return queryResult, nil
}