package middleware

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logMutex sync.Mutex
)

// LoggingMiddleware creates middleware for logging.
func LoggingMiddleware(filename string) func(http.Handler) http.Handler {
	dir := filepath.Dir(filename)
	if dir != "" && dir != "." {
		os.MkdirAll(dir, 0755)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)
			duration := time.Since(start)

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			logEntry := fmt.Sprintf("[%s] %s %s - %d - %v\n",
				timestamp, r.Method, r.URL.Path, rw.statusCode, duration)

			logMutex.Lock()
			file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				file.WriteString(logEntry)
				file.Close()
			} else {
				fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
			}
			logMutex.Unlock()
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader writes header of a response.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
