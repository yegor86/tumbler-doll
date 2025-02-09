package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	temporal "go.temporal.io/sdk/client"
)

func StreamLogs(wfClient temporal.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		workflowId := r.URL.Query().Get("workflowId")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		for {
			var msg string
			msgEncoded, err := wfClient.QueryWorkflow(r.Context(), workflowId, "", "getLogs")
			if err != nil {
				log.Println("Error querying workflow:", err)
				break
			}

			msgEncoded.Get(&msg)
			if msg != "" {
                formatedMsg := fmt.Sprintf("data: %s\n\n", msg)
		        _, err := fmt.Fprint(w, formatedMsg)
				if err != nil {
					log.Println("Error writing to response:", err)
					break
				}
				flusher.Flush()
			}
            time.Sleep(100 * time.Millisecond)
		}

		fmt.Fprintf(w, "Completed workflow: WorkflowID=%s", workflowId)
	}
}
