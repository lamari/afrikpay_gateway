package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"afrikpay/services/auth/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockJWTService is a mock implementation of JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID, email string, roles []string) (string, error) {
	args := m.Called(userID, email, roles)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateTokenPair(userID, email string, roles []string) (*models.TokenPair, error) {
	args := m.Called(userID, email, roles)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TokenPair), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*models.CustomClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CustomClaims), args.Error(1)
}

func (m *MockJWTService) RefreshToken(refreshToken string) (*models.TokenPair, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TokenPair), args.Error(1)
}

func (m *MockJWTService) LoadKeys(privateKeyPath, publicKeyPath string) error {
	args := m.Called(privateKeyPath, publicKeyPath)
	return args.Error(0)
}

func (m *MockJWTService) GetClaims(tokenString string) (*models.CustomClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CustomClaims), args.Error(1)
}

func (m *MockJWTService) IsTokenExpired(tokenString string) bool {
	args := m.Called(tokenString)
	return args.Bool(0)
}

func (m *MockJWTService) RevokeToken(tokenID string) error {
	args := m.Called(tokenID)
	return args.Error(0)
}

func (m *MockJWTService) GetPublicKey() *rsa.PublicKey {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*rsa.PublicKey)
}

func (m *MockJWTService) GetPrivateKey() *rsa.PrivateKey {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*rsa.PrivateKey)
}

// TestAuthMiddleware_ValidToken tests middleware with valid token
func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Given
	mockJWT := new(MockJWTService)
	middleware := NewAuthMiddleware(mockJWT)

	expectedClaims := &models.CustomClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{"user"},
	}

	mockJWT.On("ValidateToken", "valid.jwt.token").Return(expectedClaims, nil)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify claims are in context
		claims, ok := r.Context().Value(ClaimsKey).(*models.CustomClaims)
		assert.True(t, ok)
		assert.Equal(t, expectedClaims.UserID, claims.UserID)
		assert.Equal(t, expectedClaims.Email, claims.Email)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid.jwt.token")
	w := httptest.NewRecorder()

	// When
	middleware.Authenticate(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockJWT.AssertExpectations(t)
}

// TestAuthMiddleware_InvalidToken tests middleware with invalid token
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Given
	mockJWT := new(MockJWTService)
	middleware := NewAuthMiddleware(mockJWT)

	mockJWT.On("ValidateToken", "invalid.jwt.token").Return(nil, errors.New("invalid token"))

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")
	w := httptest.NewRecorder()

	// When
	middleware.Authenticate(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockJWT.AssertExpectations(t)
}

// TestAuthMiddleware_MissingAuthHeader tests middleware without auth header
func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	// Given
	mockJWT := new(MockJWTService)
	middleware := NewAuthMiddleware(mockJWT)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// When
	middleware.Authenticate(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddleware_InvalidAuthFormat tests middleware with invalid auth format
func TestAuthMiddleware_InvalidAuthFormat(t *testing.T) {
	// Given
	mockJWT := new(MockJWTService)
	middleware := NewAuthMiddleware(mockJWT)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Next handler should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	w := httptest.NewRecorder()

	// When
	middleware.Authenticate(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddleware_RequireRole_ValidRole tests role requirement with valid role
func TestAuthMiddleware_RequireRole_ValidRole(t *testing.T) {
	// Given
	mockJWT := &MockJWTService{}
	middleware := NewAuthMiddleware(mockJWT)
	
	claims := &models.CustomClaims{
		UserID: "user123",
		Email:  "admin@example.com",
		Roles:  []string{"user", "admin"},
	}
	
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
	
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req = req.WithContext(context.WithValue(req.Context(), ClaimsKey, claims))
	w := httptest.NewRecorder()

	// When
	middleware.RequireRole("admin")(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
}

// TestAuthMiddleware_RequireRole_InvalidRole tests role requirement with invalid role
func TestAuthMiddleware_RequireRole_InvalidRole(t *testing.T) {
	// Given
	mockJWT := &MockJWTService{}
	middleware := NewAuthMiddleware(mockJWT)
	
	claims := &models.CustomClaims{
		UserID: "user123",
		Email:  "user@example.com",
		Roles:  []string{"user"},
	}
	
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
	
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req = req.WithContext(context.WithValue(req.Context(), ClaimsKey, claims))
	w := httptest.NewRecorder()

	// When
	middleware.RequireRole("admin")(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusForbidden, w.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "INSUFFICIENT_PERMISSIONS", response.Code)
}

// TestAuthMiddleware_RequireAnyRole_ValidRole tests any role requirement with valid role
func TestAuthMiddleware_RequireAnyRole_ValidRole(t *testing.T) {
	// Given
	mockJWT := &MockJWTService{}
	middleware := NewAuthMiddleware(mockJWT)
	
	claims := &models.CustomClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{"user", "admin"},
	}
	
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
	
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req = req.WithContext(context.WithValue(req.Context(), ClaimsKey, claims))
	w := httptest.NewRecorder()

	// When
	middleware.RequireRole("admin", "moderator")(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
}

// TestAuthMiddleware_RequireAnyRole_NoValidRole tests any role requirement without valid role
func TestAuthMiddleware_RequireAnyRole_NoValidRole(t *testing.T) {
	// Given
	mockJWT := &MockJWTService{}
	middleware := NewAuthMiddleware(mockJWT)
	
	claims := &models.CustomClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{"user"},
	}
	
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
	
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req = req.WithContext(context.WithValue(req.Context(), ClaimsKey, claims))
	w := httptest.NewRecorder()

	// When
	middleware.RequireRole("admin", "moderator")(nextHandler).ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusForbidden, w.Code)
	
	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "INSUFFICIENT_PERMISSIONS", response.Code)
}

// TestAuthMiddleware_GetClaimsFromContext tests getting claims from context
func TestAuthMiddleware_GetClaimsFromContext(t *testing.T) {
	// Given
	expectedClaims := &models.CustomClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{"user"},
	}

	ctx := context.WithValue(context.Background(), ClaimsKey, expectedClaims)

	// When
	claims, ok := GetClaimsFromContext(ctx)

	// Then
	assert.True(t, ok)
	assert.NotNil(t, claims)
	assert.Equal(t, expectedClaims.UserID, claims.UserID)
	assert.Equal(t, expectedClaims.Email, claims.Email)
}

// TestAuthMiddleware_GetClaimsFromContext_NoClaims tests getting claims when none exist
func TestAuthMiddleware_GetClaimsFromContext_NoClaims(t *testing.T) {
	// Given
	ctx := context.Background()

	// When
	claims, ok := GetClaimsFromContext(ctx)

	// Then
	assert.False(t, ok)
	assert.Nil(t, claims)
}
