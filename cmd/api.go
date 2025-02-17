package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	cli "github.com/spf13/cobra"
	temporal "go.temporal.io/sdk/client"

	"github.com/yegor86/tumbler-doll/internal/api/v1/handler"
	"github.com/yegor86/tumbler-doll/internal/grpc"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var (
	apiCmd = &cli.Command{
		Use:   "api",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) {
			var err error

			wfClient, ok := cmd.Context().Value("wfClient").(temporal.Client)
			if !ok {
				log.Fatalf("Failed to obtain temporal client")
			}

			newGrpcServer()

			// Create the router and server config
			router, err := newRouter()
			if err != nil {
				log.Fatalf("Router config error: %v", err)
			}

			router.Get("/upload", handler.UploadForm)
			router.Get("/jobs", handler.ListJobs("/"))
			router.Get("/jobs/*", handler.ListJobs("/"))
			router.Post("/submit/*", handler.SubmitJob(wfClient))
			router.Post("/uploadfile", handler.UploadFile(wfClient))
			router.HandleFunc("/stream/*", handler.StreamLogs(wfClient))

			// Create a server
			s := http.Server{
				Addr:    net.JoinHostPort(config.Server.Host, config.Server.Port),
				Handler: router,
			}

			// Start the listener and service connections.
			if err = s.ListenAndServe(); err != nil {
				log.Fatalf("Server error: %v", err)
			}

			log.Printf("API listening on %s", s.Addr)
		},
	}
)

// newGrpcServer: Load a GRPC server
func newGrpcServer() {
	grpcServer := grpc.NewServer()
	go func() {
		if err := grpcServer.ListenAndServe(func(workflowId string, chunk string) {
			
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

		}); err != nil {
			log.Fatalf("GRPC server error: %v", err)
		}
	}()
}

func newRouter() (chi.Router, error) {

	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer, // Recover from panics
		middleware.RequestID, // Inject request-id
	)

	// Request logger
	if config.Server.Log.Enabled {
		// router.Use(logger.LoggerStandardMiddleware(log.Logger.With("context", "server"), loggerConfig))
	}

	// CORS handler
	if config.Server.CORS.Enabled {
		var corsOptions cors.Options
		if err := koanfConfig.Unmarshal("server.cors", &corsOptions); err != nil {
			return nil, fmt.Errorf("could not parser server.cors config: %w", err)
		}
		router.Use(cors.New(corsOptions).Handler)
	}

	return router, nil
}
