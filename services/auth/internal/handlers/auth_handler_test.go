package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"afrikpay/services/auth/internal/models"
	"crypto/rsa"
)

// MockJWTService is a mock implementation of JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) LoadKeys(privateKeyPath, publicKeyPath string) error {
	args := m.Called(privateKeyPath, publicKeyPath)
	return args.Error(0)
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

func (m *MockJWTService) RevokeToken(tokenString string) error {
	args := m.Called(tokenString)
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

// TestAuthHandler_Login_Success tests successful login
func TestAuthHandler_Login_Success(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	
	expectedTokenPair := &models.TokenPair{
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_123",
		ExpiresIn:    3600,
	}
	
	mockJWTService.On("GenerateTokenPair", "user123", "test@example.com", []string{"user"}).Return(expectedTokenPair, nil)
	
	handler := NewAuthHandler(mockJWTService)
	
	loginReq := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// When
	handler.Login(rr, req)
	
	// Then
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response models.LoginResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTokenPair.AccessToken, response.AccessToken)
	assert.Equal(t, expectedTokenPair.RefreshToken, response.RefreshToken)
	
	mockJWTService.AssertExpectations(t)
}

// TestAuthHandler_Login_InvalidCredentials tests login with invalid credentials
func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	handler := NewAuthHandler(mockJWTService)
	
	loginReq := models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrong_password",
	}
	
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// When
	handler.Login(rr, req)
	
	// Then
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var errorResp models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_CREDENTIALS", errorResp.Code)
}

// TestAuthHandler_Login_InvalidRequest tests login with invalid request body
func TestAuthHandler_Login_InvalidRequest(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	handler := NewAuthHandler(mockJWTService)
	
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// When
	handler.Login(rr, req)
	
	// Then
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	
	var errorResp models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_REQUEST", errorResp.Code)
}

// TestAuthHandler_Login_ValidationError tests login with validation errors
func TestAuthHandler_Login_ValidationError(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	handler := NewAuthHandler(mockJWTService)
	
	loginReq := models.LoginRequest{
		Email:    "invalid-email",
		Password: "123", // Too short
	}
	
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// When
	handler.Login(rr, req)
	
	// Then
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	
	var errorResp models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_ERROR", errorResp.Code)
}

// TestAuthHandler_Login_TokenGenerationError tests login with token generation error
func TestAuthHandler_Login_TokenGenerationError(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	
	mockJWTService.On("GenerateTokenPair", "user123", "test@example.com", []string{"user"}).Return(nil, errors.New("token generation failed"))
	
	handler := NewAuthHandler(mockJWTService)
	
	loginReq := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// When
	handler.Login(rr, req)
	
	// Then
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	
	mockJWTService.AssertExpectations(t)
}

// TestAuthHandler_RefreshToken_Success tests successful token refresh
func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	
	expectedTokenPair := &models.TokenPair{
		AccessToken:  "new_access_token_123",
		RefreshToken: "new_refresh_token_123",
		ExpiresIn:    3600,
	}
	
	mockJWTService.On("RefreshToken", "refresh_token_123").Return(expectedTokenPair, nil)
	
	handler := NewAuthHandler(mockJWTService)
	
	refreshReq := models.RefreshTokenRequest{
		RefreshToken: "refresh_token_123",
	}
	
	reqBody, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// When
	handler.RefreshToken(rr, req)
	
	// Then
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response models.TokenPair
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTokenPair.AccessToken, response.AccessToken)
	assert.Equal(t, expectedTokenPair.RefreshToken, response.RefreshToken)
	
	mockJWTService.AssertExpectations(t)
}

// TestAuthHandler_RefreshToken_InvalidToken tests refresh with invalid token
func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	
	mockJWTService.On("RefreshToken", "invalid_token").Return(nil, errors.New("invalid refresh token"))
	
	handler := NewAuthHandler(mockJWTService)
	
	refreshReq := models.RefreshTokenRequest{
		RefreshToken: "invalid_token",
	}
	
	reqBody, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// When
	handler.RefreshToken(rr, req)
	
	// Then
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var errorResp models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_REFRESH_TOKEN", errorResp.Code)
	
	mockJWTService.AssertExpectations(t)
}

// TestAuthHandler_VerifyToken_Success tests successful token verification
func TestAuthHandler_VerifyToken_Success(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	
	expectedClaims := &models.CustomClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{"user"},
	}
	
	mockJWTService.On("ValidateToken", "valid_token").Return(expectedClaims, nil)
	
	handler := NewAuthHandler(mockJWTService)
	
	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	
	rr := httptest.NewRecorder()
	
	// When
	handler.VerifyToken(rr, req)
	
	// Then
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response models.CustomClaims
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedClaims.UserID, response.UserID)
	assert.Equal(t, expectedClaims.Email, response.Email)
	
	mockJWTService.AssertExpectations(t)
}

// TestAuthHandler_VerifyToken_MissingToken tests verification without token
func TestAuthHandler_VerifyToken_MissingToken(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	handler := NewAuthHandler(mockJWTService)
	
	req := httptest.NewRequest("GET", "/verify", nil)
	rr := httptest.NewRecorder()
	
	// When
	handler.VerifyToken(rr, req)
	
	// Then
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var errorResp models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "MISSING_TOKEN", errorResp.Code)
}

// TestAuthHandler_VerifyToken_InvalidFormat tests verification with invalid token format
func TestAuthHandler_VerifyToken_InvalidFormat(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	handler := NewAuthHandler(mockJWTService)
	
	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	rr := httptest.NewRecorder()
	
	// When
	handler.VerifyToken(rr, req)
	
	// Then
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var errorResp models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_TOKEN_FORMAT", errorResp.Code)
}

// TestAuthHandler_VerifyToken_InvalidToken tests verification with invalid token
func TestAuthHandler_VerifyToken_InvalidToken(t *testing.T) {
	// Given
	mockJWTService := new(MockJWTService)
	
	mockJWTService.On("ValidateToken", "invalid_token").Return(nil, errors.New("invalid token"))
	
	handler := NewAuthHandler(mockJWTService)
	
	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	rr := httptest.NewRecorder()
	
	// When
	handler.VerifyToken(rr, req)
	
	// Then
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var errorResp models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_TOKEN", errorResp.Code)
	
	mockJWTService.AssertExpectations(t)
}
