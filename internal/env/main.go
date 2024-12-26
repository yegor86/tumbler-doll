package env

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func LoadEnvVars() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}
	
	appDir := os.Getenv("JENKINS_HOME")
	if appDir == "" {
		appDir = homeDir
		os.Setenv("JENKINS_HOME", appDir)
	}
	if os.Getenv("WORKSPACE") == "" {
		os.Setenv("WORKSPACE", filepath.Join(appDir, "workspace"))
	}

	printEnv()
}

func printEnv() {
	for _, e := range os.Environ() {
        fmt.Println(e)
    }
}
