package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLoggingMiddleware_Success tests logging middleware with successful request
func TestLoggingMiddleware_Success(t *testing.T) {
	// Given
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", log.LstdFlags)

	middleware := NewLoggingMiddleware(logger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// When
	middleware.LogRequests(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "200")
}

// TestLoggingMiddleware_Error tests logging middleware with error response
func TestLoggingMiddleware_Error(t *testing.T) {
	// Given
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", log.LstdFlags)

	middleware := NewLoggingMiddleware(logger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	})

	req := httptest.NewRequest(http.MethodPost, "/error", nil)
	w := httptest.NewRecorder()

	// When
	middleware.LogRequests(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", w.Body.String())

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "POST")
	assert.Contains(t, logOutput, "/error")
	assert.Contains(t, logOutput, "500")
}

// TestLoggingMiddleware_LogErrors tests error logging middleware
func TestLoggingMiddleware_LogErrors(t *testing.T) {
	// Given
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", log.LstdFlags)

	middleware := NewLoggingMiddleware(logger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})

	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	w := httptest.NewRecorder()

	// When
	middleware.LogErrors(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "not found", w.Body.String())

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "ERROR")
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/notfound")
	assert.Contains(t, logOutput, "404")
}

// TestLoggingMiddleware_LogErrors_NoError tests error logging middleware with successful response
func TestLoggingMiddleware_LogErrors_NoError(t *testing.T) {
	// Given
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", log.LstdFlags)

	middleware := NewLoggingMiddleware(logger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	req := httptest.NewRequest(http.MethodGet, "/success", nil)
	w := httptest.NewRecorder()

	// When
	middleware.LogErrors(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())

	logOutput := logBuffer.String()
	// Should not contain ERROR for successful requests
	assert.NotContains(t, logOutput, "ERROR")
}
