package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yegor86/tumbler-doll/internal/jenkins/jobs"
)

func ListJobs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		jobDb := jobs.GetInstance()

		if err := json.NewEncoder(w).Encode(jobDb.ListJobs()); err != nil {
			http.Error(w, "Failed to encode jobs as JSON", http.StatusInternalServerError)
		}
		
	}
}