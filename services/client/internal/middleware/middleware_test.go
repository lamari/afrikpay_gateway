package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "GET request",
			method:         "GET",
			path:           "/test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request",
			method:         "POST",
			path:           "/api/test",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Error response",
			method:         "GET",
			path:           "/error",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.Path, "error") {
					w.WriteHeader(http.StatusInternalServerError)
				} else if r.Method == "POST" {
					w.WriteHeader(http.StatusCreated)
				} else {
					w.WriteHeader(http.StatusOK)
				}
				w.Write([]byte("test response"))
			})

			middleware := LoggingMiddleware(handler)
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("User-Agent", "test-agent")
			w := httptest.NewRecorder()

			// When
			middleware.ServeHTTP(w, req)

			// Then
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "test response", w.Body.String())
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   bool
	}{
		{
			name:           "Regular request",
			method:         "GET",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "OPTIONS preflight request",
			method:         "OPTIONS",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "POST request",
			method:         "POST",
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("test response"))
			})

			middleware := CORSMiddleware(handler)
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			// When
			middleware.ServeHTTP(w, req)

			// Then
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkHeaders {
				assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "Content-Type, Authorization, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
			}

			// OPTIONS requests should not call the next handler
			if tt.method == "OPTIONS" {
				assert.Empty(t, w.Body.String())
			} else {
				assert.Equal(t, "test response", w.Body.String())
			}
		})
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	middleware := SecurityHeadersMiddleware(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// When
	middleware.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test response", w.Body.String())

	// Check security headers
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
}

func TestRateLimitMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		requestsPerMin  int
		numRequests     int
		expectedStatus  []int
		description     string
	}{
		{
			name:           "Within rate limit",
			requestsPerMin: 5,
			numRequests:    3,
			expectedStatus: []int{http.StatusOK, http.StatusOK, http.StatusOK},
			description:    "All requests should pass",
		},
		{
			name:           "Exceed rate limit",
			requestsPerMin: 2,
			numRequests:    3,
			expectedStatus: []int{http.StatusOK, http.StatusOK, http.StatusTooManyRequests},
			description:    "Third request should be rate limited",
		},
		{
			name:           "Single request limit",
			requestsPerMin: 1,
			numRequests:    2,
			expectedStatus: []int{http.StatusOK, http.StatusTooManyRequests},
			description:    "Second request should be rate limited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("test response"))
			})

			middleware := RateLimitMiddleware(tt.requestsPerMin)(handler)

			// When & Then
			for i := 0; i < tt.numRequests; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "127.0.0.1:12345" // Same IP for all requests
				w := httptest.NewRecorder()

				middleware.ServeHTTP(w, req)

				assert.Equal(t, tt.expectedStatus[i], w.Code, 
					"Request %d: expected status %d, got %d", i+1, tt.expectedStatus[i], w.Code)

				if tt.expectedStatus[i] == http.StatusOK {
					assert.Equal(t, "test response", w.Body.String())
				} else if tt.expectedStatus[i] == http.StatusTooManyRequests {
					assert.Contains(t, w.Body.String(), "Rate limit exceeded")
				}
			}
		})
	}
}

func TestRateLimitMiddleware_DifferentIPs(t *testing.T) {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	middleware := RateLimitMiddleware(1)(handler) // Only 1 request per minute

	// When & Then - Different IPs should have separate rate limits
	ips := []string{"127.0.0.1:12345", "192.168.1.1:54321", "10.0.0.1:9876"}
	
	for _, ip := range ips {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request from IP %s should succeed", ip)
		assert.Equal(t, "test response", w.Body.String())
	}
}

func TestRateLimitMiddleware_TimeWindow(t *testing.T) {
	// This test is more complex as it involves time manipulation
	// For now, we'll test the basic functionality
	// In a real scenario, you might want to use a time mocking library

	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	middleware := RateLimitMiddleware(2)(handler)

	// When - Make 2 requests (should succeed)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// When - Make 3rd request (should fail)
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRecoveryMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		shouldPanic bool
		panicValue  interface{}
	}{
		{
			name:        "No panic",
			shouldPanic: false,
		},
		{
			name:        "String panic",
			shouldPanic: true,
			panicValue:  "test panic",
		},
		{
			name:        "Error panic",
			shouldPanic: true,
			panicValue:  assert.AnError,
		},
		{
			name:        "Nil panic",
			shouldPanic: true,
			panicValue:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.shouldPanic {
					panic(tt.panicValue)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("test response"))
			})

			middleware := RecoveryMiddleware(handler)
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// When
			middleware.ServeHTTP(w, req)

			// Then
			if tt.shouldPanic {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Contains(t, w.Body.String(), "Internal Server Error")
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, "test response", w.Body.String())
			}
		})
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	// Given
	w := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	// When
	rw.WriteHeader(http.StatusCreated)

	// Then
	assert.Equal(t, http.StatusCreated, rw.statusCode)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		expected   string
	}{
		{
			name: "X-Forwarded-For header",
			headers: map[string]string{
				"X-Forwarded-For": "192.168.1.1",
			},
			remoteAddr: "127.0.0.1:12345",
			expected:   "192.168.1.1",
		},
		{
			name: "X-Real-IP header",
			headers: map[string]string{
				"X-Real-IP": "10.0.0.1",
			},
			remoteAddr: "127.0.0.1:12345",
			expected:   "10.0.0.1",
		},
		{
			name: "X-Forwarded-For takes precedence",
			headers: map[string]string{
				"X-Forwarded-For": "192.168.1.1",
				"X-Real-IP":       "10.0.0.1",
			},
			remoteAddr: "127.0.0.1:12345",
			expected:   "192.168.1.1",
		},
		{
			name:       "Fall back to RemoteAddr",
			headers:    map[string]string{},
			remoteAddr: "127.0.0.1:12345",
			expected:   "127.0.0.1:12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tt.remoteAddr
			
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			// When
			result := getClientIP(req)

			// Then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMiddlewareChaining(t *testing.T) {
	// Given
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Chain multiple middlewares
	middleware := RecoveryMiddleware(
		SecurityHeadersMiddleware(
			CORSMiddleware(
				LoggingMiddleware(handler),
			),
		),
	)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// When
	middleware.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test response", w.Body.String())

	// Check that all middleware headers are set
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
}
