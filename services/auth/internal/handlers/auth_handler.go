package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"afrikpay/services/auth/internal/models"
	"afrikpay/services/auth/internal/services"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	jwtService services.JWTService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(jwtService services.JWTService) *AuthHandler {
	return &AuthHandler{
		jwtService: jwtService,
	}
}

// Login handles user login requests
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := loginReq.Validate(); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// In a real implementation, you would validate credentials against a database
	// For now, we'll use a simple mock validation
	if !h.validateCredentials(loginReq.Email, loginReq.Password) {
		h.writeErrorResponse(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", "")
		return
	}

	// Create user object (in real implementation, fetch from database)
	user := models.User{
		ID:        "user123", // This would come from database
		Email:     loginReq.Email,
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Generate token pair
	tokenPair, err := h.jwtService.GenerateTokenPair(user.ID, user.Email, user.Roles)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "TOKEN_GENERATION_ERROR", "Failed to generate tokens", err.Error())
		return
	}

	// Create login response
	response := models.NewLoginResponse(
		tokenPair.AccessToken,
		tokenPair.RefreshToken,
		int64(24*time.Hour.Seconds()), // 24 hours
		user,
	)

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var refreshReq models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := refreshReq.Validate(); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Refresh token
	newTokenPair, err := h.jwtService.RefreshToken(refreshReq.RefreshToken)
	if err != nil {
		h.writeErrorResponse(w, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Failed to refresh token", err.Error())
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newTokenPair)
}

// VerifyToken handles token verification requests
func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required", "")
		return
	}

	// Extract token from "Bearer <token>" format
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		h.writeErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Invalid authorization header format", "")
		return
	}

	token := authHeader[len(bearerPrefix):]

	// Validate token
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		h.writeErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Token validation failed", err.Error())
		return
	}

	// Return claims
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(claims)
}

// validateCredentials validates user credentials (mock implementation)
func (h *AuthHandler) validateCredentials(email, password string) bool {
	// In a real implementation, this would check against a database
	// For testing purposes, we'll accept any valid email with password "password123"
	return email != "" && password == "password123"
}

// writeErrorResponse writes an error response to the HTTP response writer
func (h *AuthHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, code, message, details string) {
	errorResp := models.NewErrorResponse(code, message, details)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResp)
}
