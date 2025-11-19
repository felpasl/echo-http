package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Simulate latency
		if delay := os.Getenv("DELAY"); delay != "" {
			if ms, err := strconv.Atoi(delay); err == nil {
				time.Sleep(time.Duration(ms) * time.Millisecond)
			}
		}

		// Set status code
		status := 200
		if s := os.Getenv("STATUS_CODE"); s != "" {
			if code, err := strconv.Atoi(s); err == nil {
				status = code
			}
		}
		w.WriteHeader(status)

		// Read body
		body, _ := io.ReadAll(r.Body)

		// Build response
		var response strings.Builder
		response.WriteString(fmt.Sprintf("Method: %s\n\n", r.Method))
		response.WriteString(fmt.Sprintf("Path: %s\n\n", r.URL.Path))
		response.WriteString("Headers:\n")
		var keys []string
		for k := range r.Header {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := r.Header[k]
			response.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(v, ", ")))
		}
		response.WriteString("\n")
		if len(body) > 0 {
			response.WriteString(fmt.Sprintf("Body: %s\n", string(body)))
		}

		// Add padding if RESPONSE_SIZE is set
		if sizeStr := os.Getenv("RESPONSE_SIZE"); sizeStr != "" {
			if size, err := strconv.Atoi(sizeStr); err == nil && size > len(response.String()) {
				padding := strings.Repeat("x", size-len(response.String()))
				response.WriteString(padding)
			}
		}

		w.Write([]byte(response.String()))

		duration := time.Since(start)
		ms := duration.Seconds() * 1000
		fmt.Printf("Method: %s Path:%s status: %d time: %.3f ms\n", r.Method, r.URL.Path, status, ms)
	})

	fmt.Printf("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
