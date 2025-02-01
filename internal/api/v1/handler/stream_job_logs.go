package handler

import (
    "bufio"
    "fmt"
    "net/http"
    "os/exec"
    "time"
)

func StreamLogs(w http.ResponseWriter, r *http.Request) {
    // Set headers for SSE
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    // Run a command and stream its output
    cmd := exec.Command("ping", "-c", "5", "google.com") // Replace with your command
    stdout, _ := cmd.StdoutPipe()
    cmd.Start()

    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        fmt.Fprintf(w, "data: %s\n\n", scanner.Text())
        w.(http.Flusher).Flush() // Flush output to client
        time.Sleep(500 * time.Millisecond) // Simulate delay
    }
    cmd.Wait()
}
