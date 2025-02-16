package sse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CopyAndFlush: copies line by line from io.Reader (file) to io.Writer (http.ResponseWriter)
func CopyAndFlush(w io.Writer, r io.Reader) error {
	wFlusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("streaming not supported")
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		formatedMsg := fmt.Sprintf("data: %s\n\n", scanner.Text())
		_, err := fmt.Fprint(w, formatedMsg)
		if err != nil {
			return err
		}
		wFlusher.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	return scanner.Err()
}