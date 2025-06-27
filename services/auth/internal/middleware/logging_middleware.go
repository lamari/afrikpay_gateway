package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware handles request logging
type LoggingMiddleware struct {
	logger *log.Logger
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n
	return n, err
}

// LogRequests is a middleware that logs HTTP requests
func (m *LoggingMiddleware) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     0,
		}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Log the request
		duration := time.Since(start)
		m.logger.Printf(
			"%s %s %d %d bytes %v %s",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			wrapped.written,
			duration,
			r.RemoteAddr,
		)
	})
}

// LogErrors is a middleware that logs errors
func (m *LoggingMiddleware) LogErrors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     0,
		}

		next.ServeHTTP(wrapped, r)

		// Log errors (4xx and 5xx status codes)
		if wrapped.statusCode >= 400 {
			m.logger.Printf(
				"ERROR: %s %s returned %d",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
			)
		}
	})
}
