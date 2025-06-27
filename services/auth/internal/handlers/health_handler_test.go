package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"afrikpay/services/auth/internal/models"
)

// TestHealthHandler_Health_Success tests health check endpoint
func TestHealthHandler_Health_Success(t *testing.T) {
	// Given
	handler := NewHealthHandler("1.0.0")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// When
	handler.Health(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response.Status)
	assert.Equal(t, "1.0.0", response.Version)
	assert.NotEmpty(t, response.Uptime)
	assert.False(t, response.Timestamp.IsZero())
	assert.True(t, response.Timestamp.Before(time.Now().Add(time.Second)))
}

// TestHealthHandler_Ready_Success tests readiness check endpoint
func TestHealthHandler_Ready_Success(t *testing.T) {
	// Given
	handler := NewHealthHandler("1.0.0")
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	// When
	handler.Ready(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "ready", response["status"])
	assert.NotNil(t, response["timestamp"])
}

// TestHealthHandler_Live_Success tests liveness check endpoint
func TestHealthHandler_Live_Success(t *testing.T) {
	// Given
	handler := NewHealthHandler("1.0.0")
	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	w := httptest.NewRecorder()

	// When
	handler.Live(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "alive", response["status"])
	assert.NotNil(t, response["timestamp"])
}

// TestHealthHandler_NewHealthHandler tests handler creation
func TestHealthHandler_NewHealthHandler(t *testing.T) {
	// Given
	version := "2.0.0"

	// When
	handler := NewHealthHandler(version)

	// Then
	assert.NotNil(t, handler)
	// Note: startTime and version are private fields, so we can't test them directly
	// We test them indirectly through the Health endpoint
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handler.Health(w, req)

	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, version, response.Version)
}

// TestHealthHandler_ContentType tests content type headers
func TestHealthHandler_ContentType(t *testing.T) {
	// Given
	handler := NewHealthHandler("1.0.0")
	
	testCases := []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
		path    string
	}{
		{"Health", handler.Health, "/health"},
		{"Ready", handler.Ready, "/ready"},
		{"Live", handler.Live, "/live"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			w := httptest.NewRecorder()

			// When
			tc.handler(w, req)

			// Then
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}

// TestHealthHandler_JSONResponse tests that all responses are valid JSON
func TestHealthHandler_JSONResponse(t *testing.T) {
	// Given
	handler := NewHealthHandler("1.0.0")
	
	testCases := []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
		path    string
	}{
		{"Health", handler.Health, "/health"},
		{"Ready", handler.Ready, "/ready"},
		{"Live", handler.Live, "/live"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			w := httptest.NewRecorder()

			// When
			tc.handler(w, req)

			// Then
			assert.Equal(t, http.StatusOK, w.Code)
			
			// Verify response is valid JSON
			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON")
		})
	}
}
