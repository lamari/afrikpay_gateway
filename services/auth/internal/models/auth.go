package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Validate validates the login request
func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	
	if !isValidEmail(r.Email) {
		return fmt.Errorf("invalid email format")
	}
	
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}
	
	if len(r.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	
	return nil
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         User   `json:"user"`
}

// NewLoginResponse creates a new login response
func NewLoginResponse(accessToken, refreshToken string, expiresIn int64, user User) *LoginResponse {
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         user,
	}
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Validate validates the refresh token request
func (r *RefreshTokenRequest) Validate() error {
	if r.RefreshToken == "" {
		return fmt.Errorf("refresh token is required")
	}
	return nil
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.HasRole("admin")
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// Validate validates the user data
func (u *User) Validate() error {
	if u.ID == "" {
		return fmt.Errorf("user ID is required")
	}
	
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	
	if !isValidEmail(u.Email) {
		return fmt.Errorf("invalid email format")
	}
	
	if len(u.Roles) == 0 {
		return fmt.Errorf("user must have at least one role")
	}
	
	return nil
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message, details string) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
}

// NewHealthResponse creates a new health response
func NewHealthResponse(status, version string) *HealthResponse {
	return &HealthResponse{
		Status:    status,
		Version:   version,
		Timestamp: time.Now(),
		Uptime:    time.Since(time.Now()).String(), // Will be calculated properly in actual implementation
	}
}

// AuthError represents an authentication error
type AuthError struct {
	Message string
	Code    string
}

// Error implements the error interface
func (e *AuthError) Error() string {
	return e.Message
}

// NewAuthError creates a new authentication error
func NewAuthError(message string) *AuthError {
	return &AuthError{
		Message: message,
		Code:    "AUTH_ERROR",
	}
}

// isValidEmail validates email format using regex
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// SanitizeInput sanitizes user input by trimming whitespace and normalizing
func SanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Replace multiple whitespace characters with single space
	re := regexp.MustCompile(`\s+`)
	input = re.ReplaceAllString(input, " ")
	
	return input
}
