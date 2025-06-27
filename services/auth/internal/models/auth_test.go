package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestLoginRequest_Validate_ValidData tests login request validation with valid data
func TestLoginRequest_Validate_ValidData(t *testing.T) {
	// Given
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// When
	err := req.Validate()

	// Then
	assert.NoError(t, err)
}

// TestLoginRequest_Validate_EmptyEmail tests login request validation with empty email
func TestLoginRequest_Validate_EmptyEmail(t *testing.T) {
	// Given
	req := LoginRequest{
		Email:    "",
		Password: "password123",
	}

	// When
	err := req.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email is required")
}

// TestLoginRequest_Validate_InvalidEmail tests login request validation with invalid email
func TestLoginRequest_Validate_InvalidEmail(t *testing.T) {
	// Given
	req := LoginRequest{
		Email:    "invalid-email",
		Password: "password123",
	}

	// When
	err := req.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email format")
}

// TestLoginRequest_Validate_EmptyPassword tests login request validation with empty password
func TestLoginRequest_Validate_EmptyPassword(t *testing.T) {
	// Given
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "",
	}

	// When
	err := req.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password is required")
}

// TestLoginRequest_Validate_ShortPassword tests login request validation with short password
func TestLoginRequest_Validate_ShortPassword(t *testing.T) {
	// Given
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "123",
	}

	// When
	err := req.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password must be at least 6 characters")
}

// TestLoginResponse_NewLoginResponse tests creating new login response
func TestLoginResponse_NewLoginResponse(t *testing.T) {
	// Given
	token := "jwt.token.here"
	refreshToken := "refresh.token.here"
	expiresIn := int64(3600)
	user := User{
		ID:    "user123",
		Email: "test@example.com",
		Roles: []string{"user"},
	}

	// When
	response := NewLoginResponse(token, refreshToken, expiresIn, user)

	// Then
	assert.Equal(t, token, response.AccessToken)
	assert.Equal(t, refreshToken, response.RefreshToken)
	assert.Equal(t, "Bearer", response.TokenType)
	assert.Equal(t, expiresIn, response.ExpiresIn)
	assert.Equal(t, user, response.User)
}

// TestRefreshTokenRequest_Validate_ValidToken tests refresh token request validation
func TestRefreshTokenRequest_Validate_ValidToken(t *testing.T) {
	// Given
	req := RefreshTokenRequest{
		RefreshToken: "valid.refresh.token",
	}

	// When
	err := req.Validate()

	// Then
	assert.NoError(t, err)
}

// TestRefreshTokenRequest_Validate_EmptyToken tests refresh token request with empty token
func TestRefreshTokenRequest_Validate_EmptyToken(t *testing.T) {
	// Given
	req := RefreshTokenRequest{
		RefreshToken: "",
	}

	// When
	err := req.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refresh token is required")
}

// TestUser_IsAdmin_AdminRole tests user admin role check
func TestUser_IsAdmin_AdminRole(t *testing.T) {
	// Given
	user := User{
		ID:    "user123",
		Email: "admin@example.com",
		Roles: []string{"user", "admin"},
	}

	// When
	isAdmin := user.IsAdmin()

	// Then
	assert.True(t, isAdmin)
}

// TestUser_IsAdmin_NoAdminRole tests user without admin role
func TestUser_IsAdmin_NoAdminRole(t *testing.T) {
	// Given
	user := User{
		ID:    "user123",
		Email: "user@example.com",
		Roles: []string{"user"},
	}

	// When
	isAdmin := user.IsAdmin()

	// Then
	assert.False(t, isAdmin)
}

// TestUser_HasRole_ExistingRole tests user role check with existing role
func TestUser_HasRole_ExistingRole(t *testing.T) {
	// Given
	user := User{
		ID:    "user123",
		Email: "user@example.com",
		Roles: []string{"user", "moderator"},
	}

	// When
	hasRole := user.HasRole("moderator")

	// Then
	assert.True(t, hasRole)
}

// TestUser_HasRole_NonExistingRole tests user role check with non-existing role
func TestUser_HasRole_NonExistingRole(t *testing.T) {
	// Given
	user := User{
		ID:    "user123",
		Email: "user@example.com",
		Roles: []string{"user"},
	}

	// When
	hasRole := user.HasRole("admin")

	// Then
	assert.False(t, hasRole)
}

// TestUser_Validate_ValidUser tests user validation with valid data
func TestUser_Validate_ValidUser(t *testing.T) {
	// Given
	user := User{
		ID:        "user123",
		Email:     "test@example.com",
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// When
	err := user.Validate()

	// Then
	assert.NoError(t, err)
}

// TestUser_Validate_EmptyID tests user validation with empty ID
func TestUser_Validate_EmptyID(t *testing.T) {
	// Given
	user := User{
		ID:        "",
		Email:     "test@example.com",
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// When
	err := user.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID is required")
}

// TestUser_Validate_InvalidEmail tests user validation with invalid email
func TestUser_Validate_InvalidEmail(t *testing.T) {
	// Given
	user := User{
		ID:        "user123",
		Email:     "invalid-email",
		Roles:     []string{"user"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// When
	err := user.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email format")
}

// TestUser_Validate_EmptyRoles tests user validation with empty roles
func TestUser_Validate_EmptyRoles(t *testing.T) {
	// Given
	user := User{
		ID:        "user123",
		Email:     "test@example.com",
		Roles:     []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// When
	err := user.Validate()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user must have at least one role")
}

// TestErrorResponse_NewErrorResponse tests creating error response
func TestErrorResponse_NewErrorResponse(t *testing.T) {
	// Given
	code := "AUTH_001"
	message := "Authentication failed"
	details := "Invalid credentials provided"

	// When
	response := NewErrorResponse(code, message, details)

	// Then
	assert.Equal(t, code, response.Code)
	assert.Equal(t, message, response.Message)
	assert.Equal(t, details, response.Details)
	assert.False(t, response.Timestamp.IsZero())
}

// TestHealthResponse_NewHealthResponse tests creating health response
func TestHealthResponse_NewHealthResponse(t *testing.T) {
	// Given
	status := "healthy"
	version := "1.0.0"

	// When
	response := NewHealthResponse(status, version)

	// Then
	assert.Equal(t, status, response.Status)
	assert.Equal(t, version, response.Version)
	assert.False(t, response.Timestamp.IsZero())
	assert.NotEmpty(t, response.Uptime)
}
