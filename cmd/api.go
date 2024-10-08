package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	cli "github.com/spf13/cobra"
	// "github.com/snowzach/gorestapi/gorestapi/mainrpc"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func init() {
	// Parse defaults, config file and environment.
	_, _, err := Load()
	if err != nil {
		Logger.Error(fmt.Sprintf("could not parse YAML config: %v", err))
		os.Exit(1)
	}
	rootCmd.AddCommand(apiCmd)
}

var (

	apiLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // Enables logging the file and line number
	}))

	apiCmd = &cli.Command{
		Use:   "api",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the databse

			// Register signal handler and wait
			signalChannel := make(chan os.Signal, 1)
			signal.Notify(signalChannel, []os.Signal{syscall.SIGINT, syscall.SIGTERM}...)
			var err error

			// Create the router and server config
			router, err := newRouter()
			if err != nil {
				apiLogger.Error("router config error: %v", err)
				close(signalChannel)
			}

			// Create the database
			// db, err := newDatabase()
			_, err = newDatabase()
			if err != nil {
				apiLogger.Error("database config error: %v", err)
				close(signalChannel)
			}

			// Version endpoint
			router.Get("/version", func(w http.ResponseWriter, r *http.Request) {
				v := struct {
					Version string `json:"version"`
				}{
					Version: GitVersion,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(v)
			})

			// MainRPC
			// if err = mainrpc.Setup(router, db); err != nil {
			// 	log.Fatalf("Could not setup mainrpc: %v", err)
			// }

			// Create a server
			s := http.Server{
				Addr:    net.JoinHostPort(config.Server.Host, config.Server.Port),
				Handler: router,
			}

			// Start the listener and service connections.
			go func() {
				if err = s.ListenAndServe(); err != nil {
					slog.Error(fmt.Sprintf("Server error: %v", err))

					close(signalChannel)
				}
			}()
			slog.Info(fmt.Sprintf("API listening on %s"), s.Addr)

			<-signalChannel
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

func newDatabase() (*sqlx.DB, error) {

	var err error
	var db *sqlx.DB

	// Create database
	db, err = sqlx.Connect("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Database))
	// db, err = sqlx.Connect("pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Password, config.Database.Database))
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database %w", err)
	}
	db.SetMaxOpenConns(config.Database.MaxConnections)

	return db, nil
}
