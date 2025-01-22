package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yegor86/tumbler-doll/internal/jenkins/jobs"
)

// Handler function for GET /jobs/{jobpath}
func ListJobs(defaultPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		jobPath := chi.URLParam(r, "*")
		
		jobDB := jobs.GetInstance()
		jobs := jobDB.ListJobs("/jobs/" + jobPath)
		
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jobs); err != nil {
			http.Error(w, "Failed to encode jobs as JSON", http.StatusInternalServerError)
		}
		
	}
}