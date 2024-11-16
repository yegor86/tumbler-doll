package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	cli "github.com/spf13/cobra"
)

func init() {
	// Parse defaults, config file and environment.
	_, _, err := Load()
	if err != nil {
		log.Fatalf("could not parse YAML config: %v", err)
	}
}

var (
	pidFile string

	// The Root Cli Handler
	rootCmd = &cli.Command{
		Version: GitVersion,
		Use:     Executable,
		PersistentPreRunE: func(cmd *cli.Command, args []string) error {

			// Load the metrics server
			if config.Metrics.Enabled {
				hostPort := net.JoinHostPort(config.Metrics.Host, strconv.Itoa(config.Metrics.Port))
				r := http.NewServeMux()
				r.Handle("/metrics", promhttp.Handler())
				if config.Profiler.Enabled {
					r.HandleFunc("/debug/pprof/", pprof.Index)
					r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
					r.HandleFunc("/debug/pprof/profile", pprof.Profile)
					r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
					r.HandleFunc("/debug/pprof/trace", pprof.Trace)
					log.Printf("Profiler enabled, profiler_path = %s", fmt.Sprintf("http://%s/debug/pprof/", hostPort))
				}
				go func() {
					if err := http.ListenAndServe(hostPort, r); err != nil {
						log.Printf("Metrics server error: %v", err)
					}
				}()
				log.Printf("Metrics enabled for address: %s", hostPort)
			}

			// Create Pid File
			pidFile = config.Profiler.Pidfile
			if pidFile != "" {
				file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
				if err != nil {
					return fmt.Errorf("could not create pid file: %s error:%v", pidFile, err)
				}
				defer file.Close()
				_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
				if err != nil {
					return fmt.Errorf("could not create pid file: %s error:%v", pidFile, err)
				}
			}
			return nil
		},
		PersistentPostRun: func(cmd *cli.Command, args []string) {
			// Remove Pid file
			if pidFile != "" {
				os.Remove(pidFile)
			}
		},
	}
)

// Execute starts the program
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
