package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"afrikpay/services/auth/internal/models"
	"afrikpay/services/auth/internal/services"
)

// ClaimsContextKey is the key used to store claims in the request context
type ClaimsContextKey string

const (
	// ClaimsKey is the context key for JWT claims
	ClaimsKey ClaimsContextKey = "claims"
)

// AuthMiddleware handles JWT authentication for protected routes
type AuthMiddleware struct {
	jwtService services.JWTService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// Authenticate is a middleware that validates JWT tokens
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.writeErrorResponse(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
			return
		}

		// Extract token from "Bearer <token>" format
		const bearerPrefix = "Bearer "
		if len(authHeader) < len(bearerPrefix) || !strings.HasPrefix(authHeader, bearerPrefix) {
			m.writeErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Invalid authorization header format")
			return
		}

		token := authHeader[len(bearerPrefix):]

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			m.writeErrorResponse(w, http.StatusUnauthorized, "INVALID_TOKEN", "Token validation failed: "+err.Error())
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		r = r.WithContext(ctx)

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// RequireRole is a middleware that requires specific roles
func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get claims from context
			claims, ok := r.Context().Value(ClaimsKey).(*models.CustomClaims)
			if !ok {
				m.writeErrorResponse(w, http.StatusUnauthorized, "MISSING_CLAIMS", "Claims not found in context")
				return
			}

			// Check if user has any of the required roles
			hasRole := false
			for _, requiredRole := range roles {
				for _, userRole := range claims.Roles {
					if userRole == requiredRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				m.writeErrorResponse(w, http.StatusForbidden, "INSUFFICIENT_PERMISSIONS", "User does not have required role")
				return
			}

			// Call next handler
			next.ServeHTTP(w, r)
		})
	}
}

// GetClaimsFromContext extracts claims from request context
func GetClaimsFromContext(ctx context.Context) (*models.CustomClaims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*models.CustomClaims)
	return claims, ok
}

// writeErrorResponse writes an error response to the HTTP response writer
func (m *AuthMiddleware) writeErrorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	errorResp := models.NewErrorResponse(code, message, "")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResp)
}
