package services

import (
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJWTService_GenerateToken_ValidClaims tests JWT token generation with valid claims
func TestJWTService_GenerateToken_ValidClaims(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	userID := "user123"
	email := "test@example.com"
	roles := []string{"user"}

	// When
	token, err := service.GenerateToken(userID, email, roles)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Contains(t, token, ".")
}

// TestJWTService_GenerateToken_EmptyUserID tests JWT generation with empty user ID
func TestJWTService_GenerateToken_EmptyUserID(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	// When
	token, err := service.GenerateToken("", "test@example.com", []string{"user"})

	// Then
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "user ID cannot be empty")
}

// TestJWTService_GenerateToken_EmptyEmail tests JWT generation with empty email
func TestJWTService_GenerateToken_EmptyEmail(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	// When
	token, err := service.GenerateToken("user123", "", []string{"user"})

	// Then
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "email cannot be empty")
}

// TestJWTService_ValidateToken_ValidToken tests JWT token validation with valid token
func TestJWTService_ValidateToken_ValidToken(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	userID := "user123"
	email := "test@example.com"
	roles := []string{"user", "admin"}

	token, err := service.GenerateToken(userID, email, roles)
	require.NoError(t, err)

	// When
	claims, err := service.ValidateToken(token)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, "test-issuer", claims.Issuer)
	assert.Equal(t, "test-audience", claims.Audience[0])
}

// TestJWTService_ValidateToken_InvalidToken tests JWT validation with invalid token
func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	invalidToken := "invalid.jwt.token"

	// When
	claims, err := service.ValidateToken(invalidToken)

	// Then
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestJWTService_ValidateToken_ExpiredToken tests JWT validation with expired token
func TestJWTService_ValidateToken_ExpiredToken(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", -1*time.Hour) // Expired immediately
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	token, err := service.GenerateToken("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	// Wait a moment to ensure expiration
	time.Sleep(10 * time.Millisecond)

	// When
	claims, err := service.ValidateToken(token)

	// Then
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token is expired")
}

// TestJWTService_LoadKeys_ValidKeys tests loading valid RSA keys
func TestJWTService_LoadKeys_ValidKeys(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)

	// When
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, service.GetPrivateKey())
	assert.NotNil(t, service.GetPublicKey())
	assert.IsType(t, &rsa.PrivateKey{}, service.GetPrivateKey())
	assert.IsType(t, &rsa.PublicKey{}, service.GetPublicKey())
}

// TestJWTService_LoadKeys_InvalidPrivateKey tests loading invalid private key
func TestJWTService_LoadKeys_InvalidPrivateKey(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)

	// When
	err := service.LoadKeys("nonexistent.pem", "../../../../config/keys/public.pem")

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read private key file")
}

// TestJWTService_LoadKeys_InvalidPublicKey tests loading invalid public key
func TestJWTService_LoadKeys_InvalidPublicKey(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)

	// When
	err := service.LoadKeys("../../../../config/keys/private.pem", "nonexistent.pem")

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read public key file")
}

// TestJWTService_RefreshToken_ValidToken tests token refresh with valid token
func TestJWTService_RefreshToken_ValidToken(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	// Generate a token pair first to get a refresh token
	originalTokenPair, err := service.GenerateTokenPair("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	// When - Use the refresh token to get a new token pair
	newTokenPair, err := service.RefreshToken(originalTokenPair.RefreshToken)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, newTokenPair)
	assert.NotEmpty(t, newTokenPair.AccessToken)
	assert.NotEmpty(t, newTokenPair.RefreshToken)
	assert.NotEqual(t, originalTokenPair.AccessToken, newTokenPair.AccessToken)

	// Validate new access token
	claims, err := service.ValidateToken(newTokenPair.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claims.UserID)
}

// TestJWTService_RefreshToken_InvalidToken tests token refresh with invalid token
func TestJWTService_RefreshToken_InvalidToken(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	// When
	newTokenPair, err := service.RefreshToken("invalid.token")

	// Then
	assert.Error(t, err)
	assert.Nil(t, newTokenPair)
	assert.Contains(t, err.Error(), "failed to parse")
}

// TestJWTService_GetClaims_ValidToken tests extracting claims from valid token
func TestJWTService_GetClaims_ValidToken(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	userID := "user123"
	email := "test@example.com"
	roles := []string{"user", "admin"}

	token, err := service.GenerateToken(userID, email, roles)
	require.NoError(t, err)

	// When
	claims, err := service.GetClaims(token)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, roles, claims.Roles)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

// TestJWTService_IsTokenExpired_ExpiredToken tests checking if token is expired
func TestJWTService_IsTokenExpired_ExpiredToken(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", -1*time.Hour)
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	token, err := service.GenerateToken("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	// When
	isExpired := service.IsTokenExpired(token)

	// Then
	assert.True(t, isExpired)
}

// TestJWTService_IsTokenExpired_ValidToken tests checking if valid token is not expired
func TestJWTService_IsTokenExpired_ValidToken(t *testing.T) {
	// Given
	service := NewJWTService("test-issuer", "test-audience", 24*time.Hour)
	err := service.LoadKeys("../../../../config/keys/private.pem", "../../../../config/keys/public.pem")
	require.NoError(t, err)

	token, err := service.GenerateToken("user123", "test@example.com", []string{"user"})
	require.NoError(t, err)

	// When
	isExpired := service.IsTokenExpired(token)

	// Then
	assert.False(t, isExpired)
}
