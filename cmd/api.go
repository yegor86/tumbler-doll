package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

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

			var wg sync.WaitGroup
        	wg.Add(2)

			// Create a HTTP server
			httpServer := http.Server{
				Addr:    net.JoinHostPort(config.Server.Host, config.Server.Port),
				Handler: router,
			}
			go func() {
                defer wg.Done()
				if err = httpServer.ListenAndServe(); err != nil {
					log.Fatalf("HTTP server error: %v", err)
				}
			}()

			// Create GRPC server
			grpcServer := grpc.NewServer()
			go func() {
				defer wg.Done()
				if err := grpcServer.ListenAndServe(); err != nil {
					log.Fatalf("GRPC server error: %v", err)
				}
			}()
			
			fmt.Printf("Servers started on ports %s and %s\n", httpServer.Addr, grpcServer.Addr)

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			<-quit

			fmt.Println("Shutting down servers...")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
					log.Fatalf("HTTP server shutdown failed: %v", err)
			}

			if err := grpcServer.Shutdown(ctx); err != nil {
					log.Fatalf("GRPC server shutdown failed: %v", err)
			}

			wg.Wait() // Wait for server goroutines to exit.
			fmt.Println("Servers gracefully stopped.")
		},
	}
)

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
