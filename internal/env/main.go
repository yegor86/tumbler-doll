package env

import (
	"fmt"
	"log"
	"net"
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

	ipAddress, err := getIPAddress()
	if err != nil {
		log.Fatalf("Error getting IP address: %v", err)
	}
	os.Setenv("TEMPORAL_HOSTPORT", fmt.Sprintf("%s:%s", ipAddress, "50051"))

	printEnv()
}

func printEnv() {
	for _, e := range os.Environ() {
        fmt.Println(e)
    }
}

// getIPAddress retrieves the system's IP address
func getIPAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to get IP addresses: %w", err)
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil { // Only get IPv4
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no valid IP address found")
}