package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	cli "github.com/spf13/cobra"

	wfClient "go.temporal.io/sdk/client"

	"github.com/yegor86/tumbler-doll/internal/api/v1/handler"
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

			client, err := wfClient.Dial(wfClient.Options{})
			if err != nil {
				log.Fatalln("Unable to create Workflow client", err)
			}
			defer client.Close()

			// Create the router and server config
			router, err := newRouter()
			if err != nil {
				log.Fatalf("Router config error: %v", err)
			}

			router.Get("/upload", handler.UploadForm)
			router.Post("/uploadfile", handler.UploadFile(client))

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