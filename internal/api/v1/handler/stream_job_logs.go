package handler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	temporal "go.temporal.io/sdk/client"

	"github.com/yegor86/tumbler-doll/internal/api/sse"
	"github.com/yegor86/tumbler-doll/internal/workflow"

	pb "github.com/yegor86/tumbler-doll/internal/grpc/proto"
)

// ReadLogs: read event logs from text file into http writer
func ReadLogs(wfClient temporal.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set headers for SSE
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		jobPath := chi.URLParam(r, "*")
		workflowId := r.URL.Query().Get("workflowId")
		
		delim := strings.LastIndex(workflowId, "/")
		jobId := workflowId[delim + 1:]
		state := workflow.Undefined

		ipath := filepath.Join(os.Getenv("JENKINS_HOME"), jobPath, "builds", jobId, "log")
		err := os.ErrNotExist
		for err != nil && state != workflow.Done {
			_, err = os.Stat(ipath)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}

		ifile, err := os.Open(ipath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer ifile.Close()

		// Remember inFile read-offset to use it in case of inFile was not fully read at the first traversal
		var seekOffset int64 = 0

		for state != workflow.Done {
			_, err = ifile.Seek(seekOffset, io.SeekStart)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = sse.CopyAndFlush(w, ifile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			seekOffset, err = ifile.Seek(0, io.SeekCurrent)
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

// WriteLogs: pipe log event into a text file
func WriteLogs(req *pb.LogRequest) {
			
	workflowId, chunk := req.WorkflowId, req.Message
	
	delim := strings.LastIndex(workflowId, "/")
	jobPath, jobId := workflowId[:delim], workflowId[delim + 1:]
	opath := filepath.Join(os.Getenv("JENKINS_HOME"), jobPath, "builds", jobId)
	err := os.MkdirAll(opath, 0740)
	if err != nil {
		log.Printf("error creating dir %s: %v", opath, err)
		return
	}

	ofile, err := os.OpenFile(filepath.Join(opath, "log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error creating/opening file %s: %v", filepath.Join(opath, "log"), err)
		return
	}

	w := bufio.NewWriter(ofile)

	// write a chunk
	if _, err := w.Write([]byte(chunk + "\n")); err != nil {
		log.Printf("error when writing log %v. Failed chunk: %s", err, chunk)
	}
	if err = w.Flush(); err != nil {
		log.Printf("error when flushing log %v. Failed chunk: %s", err, chunk)
	}
}