package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"afrikpay/services/auth/internal/models"
	"afrikpay/services/auth/internal/services"
)

// JWTAuthMiddleware creates a middleware function that validates JWT tokens
func JWTAuthMiddleware(jwtService services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeErrorResponse(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required", "")
				return
			}

			// Extract token from "Bearer <token>" format
			const bearerPrefix = "Bearer "
			if len(authHeader) < len(bearerPrefix) || !strings.HasPrefix(authHeader, bearerPrefix) {
				writeErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Invalid authorization header format", "")
				return
			}

			token := authHeader[len(bearerPrefix):]

			// Validate token
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				writeErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Token validation failed", err.Error())
				return
			}

			// Add claims to request context
			ctx := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctx)

			// Call next handler
			next.ServeHTTP(w, r)
		})
	}
}

// SimpleLoggingMiddleware creates a simple logging middleware
func SimpleLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Call next handler
		next.ServeHTTP(w, r)
		
		// Log the request
		duration := time.Since(start)
		log.Printf("%s %s %v", r.Method, r.URL.Path, duration)
	})
}

// writeErrorResponse writes an error response to the HTTP response writer
func writeErrorResponse(w http.ResponseWriter, statusCode int, code, message, details string) {
	errorResp := models.NewErrorResponse(code, message, details)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResp)
}
