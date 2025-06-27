package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"afrikpay/services/auth/internal/models"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
	version   string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// Health handles health check requests
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime)
	
	response := &models.HealthResponse{
		Status:    "healthy",
		Version:   h.version,
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Ready handles readiness check requests
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you would check dependencies like database connectivity
	response := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Live handles liveness check requests
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
