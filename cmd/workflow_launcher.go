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

	_ "github.com/jackc/pgx/v4/stdlib"
	cli "github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workflowCmd)
}

var (
	wkflLauncherLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	workflowCmd = &cli.Command{
		Use:   "wf",
		Short: "Start Workflow",
		Long:  `Start Workflow`,
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
