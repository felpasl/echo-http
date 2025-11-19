package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestEchoHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		headers        map[string]string
		body           string
		expectedStatus int
		expectedBody   string
		env            map[string]string
	}{
		{
			name:           "GET request",
			method:         "GET",
			path:           "/test",
			headers:        map[string]string{"User-Agent": "test"},
			body:           "",
			expectedStatus: 200,
			expectedBody:   "Method: GET\n\nPath: /test\n\nHeaders:\nUser-Agent: test\n\n",
			env:            map[string]string{},
		},
		{
			name:           "POST request with body",
			method:         "POST",
			path:           "/api",
			headers:        map[string]string{"Content-Type": "application/json"},
			body:           `{"key": "value"}`,
			expectedStatus: 200,
			expectedBody:   "Method: POST\n\nPath: /api\n\nHeaders:\nContent-Type: application/json\n\nBody: {\"key\": \"value\"}\n",
			env:            map[string]string{},
		},
		{
			name:           "Custom status code",
			method:         "GET",
			path:           "/",
			headers:        map[string]string{},
			body:           "",
			expectedStatus: 404,
			expectedBody:   "Method: GET\n\nPath: /\n\nHeaders:\n\n",
			env:            map[string]string{"STATUS_CODE": "404"},
		},
		{
			name:           "Response size padding",
			method:         "GET",
			path:           "/",
			headers:        map[string]string{},
			body:           "",
			expectedStatus: 200,
			expectedBody:   "Method: GET\n\nPath: /\n\nHeaders:\n\n" + strings.Repeat("x", 100-len("Method: GET\n\nPath: /\n\nHeaders:\n\n")),
			env:            map[string]string{"RESPONSE_SIZE": "100"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.env {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				// Simulate latency
				if delay := os.Getenv("DELAY"); delay != "" {
					if d, err := time.ParseDuration(delay); err == nil {
						time.Sleep(d)
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
				for k, v := range r.Header {
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
				us := duration.Microseconds() % 1000
				fmt.Printf("Method: %s Path:%s status: %d time: %.3fms %dµs\n", r.Method, r.URL.Path, status, ms, us)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestDelay(t *testing.T) {
	os.Setenv("DELAY", "10")
	defer os.Unsetenv("DELAY")

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	start := time.Now()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		for k, v := range r.Header {
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
		us := duration.Microseconds() % 1000
		fmt.Printf("Method: %s Path:%s status: %d time: %.3fms %dµs\n", r.Method, r.URL.Path, status, ms, us)
	})

	handler.ServeHTTP(w, req)
	elapsed := time.Since(start)

	if elapsed < 10*time.Millisecond {
		t.Errorf("expected delay of at least 10ms, got %v", elapsed)
	}
}
