package handler

import (
	"fmt"
	"net/http"
	"time"

	temporal "go.temporal.io/sdk/client"
)

func StreamLogs(wfClient temporal.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

        workflowId := r.URL.Query().Get("workflowId")
        
        for {
            var msg string
            msgEncoded, err := wfClient.QueryWorkflow(r.Context(), workflowId, "", "getLogs")
            if err != nil {
                break
            }

            msgEncoded.Get(&msg)
            if msg != "" {
                fmt.Fprintf(w, "%s\n", msg)
			    w.(http.Flusher).Flush()
            }
			time.Sleep(500 * time.Millisecond)
        }
        fmt.Fprintf(w, "Completed workflow: WorkflowID=%s", workflowId)
	}
}
