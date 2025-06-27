package services

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"afrikpay/services/auth/internal/models"
)

// JWTService interface defines JWT operations
type JWTService interface {
	LoadKeys(privateKeyPath, publicKeyPath string) error
	GenerateToken(userID, email string, roles []string) (string, error)
	GenerateTokenPair(userID, email string, roles []string) (*models.TokenPair, error)
	ValidateToken(tokenString string) (*models.CustomClaims, error)
	GetClaims(tokenString string) (*models.CustomClaims, error)
	IsTokenExpired(tokenString string) bool
	RefreshToken(refreshTokenString string) (*models.TokenPair, error)
	RevokeToken(tokenString string) error
	GetPublicKey() *rsa.PublicKey
	GetPrivateKey() *rsa.PrivateKey
}

// jwtService implements JWTService
type jwtService struct {
	privateKey       *rsa.PrivateKey
	publicKey        *rsa.PublicKey
	issuer           string
	audience         string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	revokedTokens    map[string]bool // In production, use Redis or database
}

// NewJWTService creates a new JWT service
func NewJWTService(issuer, audience string, accessTokenTTL time.Duration) JWTService {
	refreshTokenTTL := accessTokenTTL * 24 // Refresh token lasts 24x longer
	return &jwtService{
		issuer:          issuer,
		audience:        audience,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		revokedTokens:   make(map[string]bool),
	}
}

// LoadKeys loads RSA private and public keys from PEM files
func (s *jwtService) LoadKeys(privateKeyPath, publicKeyPath string) error {
	// Load private key
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyData)
	if privateKeyBlock == nil {
		return fmt.Errorf("failed to decode private key PEM")
	}

	// Try PKCS8 first, then PKCS1
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		// Fallback to PKCS1
		privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
		s.privateKey = privateKey
	} else {
		// PKCS8 parsed successfully, cast to RSA private key
		rsaPrivateKey, ok := privateKeyInterface.(*rsa.PrivateKey)
		if !ok {
			return fmt.Errorf("private key is not RSA")
		}
		s.privateKey = rsaPrivateKey
	}

	// Load public key
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key file: %w", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyData)
	if publicKeyBlock == nil {
		return fmt.Errorf("failed to decode public key PEM")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key is not RSA")
	}

	s.publicKey = rsaPublicKey
	return nil
}

// GenerateToken generates a single access token (for backward compatibility with tests)
func (s *jwtService) GenerateToken(userID, email string, roles []string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("user ID cannot be empty")
	}
	
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}
	
	if len(roles) == 0 {
		return "", fmt.Errorf("roles cannot be empty")
	}

	if s.privateKey == nil {
		return "", fmt.Errorf("private key not loaded")
	}
	
	// Generate access token
	accessClaims := models.NewCustomClaims(userID, email, roles, s.issuer, s.audience, s.accessTokenTTL)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	
	return accessTokenString, nil
}

// GenerateTokenPair generates access and refresh token pair
func (s *jwtService) GenerateTokenPair(userID, email string, roles []string) (*models.TokenPair, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	
	if len(roles) == 0 {
		return nil, fmt.Errorf("roles cannot be empty")
	}
	
	// Generate access token
	accessClaims := models.NewCustomClaims(userID, email, roles, s.issuer, s.audience, s.accessTokenTTL)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}
	
	// Generate refresh token
	refreshClaims := models.NewRefreshTokenClaims(userID, s.issuer, s.audience, s.refreshTokenTTL)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}
	
	return models.NewTokenPair(accessTokenString, refreshTokenString, int64(s.accessTokenTTL.Seconds())), nil
}

// ValidateToken validates and parses a JWT token
func (s *jwtService) ValidateToken(tokenString string) (*models.CustomClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}
	
	// Parse token with custom claims
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	// Extract claims
	claims, ok := token.Claims.(*models.CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	
	// Check if token is revoked
	if s.revokedTokens[claims.ID] {
		return nil, fmt.Errorf("token has been revoked")
	}
	
	// Validate custom claims
	if err := claims.Valid(); err != nil {
		return nil, fmt.Errorf("invalid claims: %w", err)
	}
	
	return claims, nil
}

// GetClaims gets the claims from a JWT token
func (s *jwtService) GetClaims(tokenString string) (*models.CustomClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}
	
	// Parse token with custom claims
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	// Extract claims
	claims, ok := token.Claims.(*models.CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	
	return claims, nil
}

// IsTokenExpired checks if a JWT token is expired
func (s *jwtService) IsTokenExpired(tokenString string) bool {
	claims, err := s.GetClaims(tokenString)
	if err != nil {
		return true
	}
	
	return claims.IsExpired()
}

// RefreshToken generates new token pair using refresh token
func (s *jwtService) RefreshToken(refreshTokenString string) (*models.TokenPair, error) {
	if refreshTokenString == "" {
		return nil, fmt.Errorf("refresh token cannot be empty")
	}
	
	// Parse refresh token
	token, err := jwt.ParseWithClaims(refreshTokenString, &models.RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}
	
	refreshClaims, ok := token.Claims.(*models.RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token claims")
	}
	
	// Check if refresh token is revoked
	if s.revokedTokens[refreshClaims.ID] {
		return nil, fmt.Errorf("refresh token has been revoked")
	}
	
	// Validate refresh token claims
	if err := refreshClaims.Valid(); err != nil {
		return nil, fmt.Errorf("invalid refresh token claims: %w", err)
	}
	
	// For this example, we'll need to get user info from somewhere
	// In a real implementation, you'd fetch from database
	userID := refreshClaims.UserID
	email := "user@example.com" // This should come from database
	roles := []string{"user"}   // This should come from database
	
	// Generate new token pair
	return s.GenerateTokenPair(userID, email, roles)
}

// RevokeToken revokes a token by adding it to revoked list
func (s *jwtService) RevokeToken(tokenString string) error {
	if tokenString == "" {
		return fmt.Errorf("token cannot be empty")
	}
	
	// Parse token to get ID
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to parse token for revocation: %w", err)
	}
	
	claims, ok := token.Claims.(*models.CustomClaims)
	if !ok {
		return fmt.Errorf("invalid token claims for revocation")
	}
	
	// Add to revoked tokens
	s.revokedTokens[claims.ID] = true
	
	return nil
}

// GetPublicKey returns the public key for token verification
func (s *jwtService) GetPublicKey() *rsa.PublicKey {
	return s.publicKey
}

// GetPrivateKey returns the private key (for testing purposes only)
func (s *jwtService) GetPrivateKey() *rsa.PrivateKey {
	return s.privateKey
}

// MockJWTService for testing
type MockJWTService struct {
	LoadKeysFunc        func(privateKeyPath, publicKeyPath string) error
	GenerateTokenFunc   func(userID, email string, roles []string) (string, error)
	GenerateTokenPairFunc func(userID, email string, roles []string) (*models.TokenPair, error)
	ValidateTokenFunc     func(tokenString string) (*models.CustomClaims, error)
	GetClaimsFunc         func(tokenString string) (*models.CustomClaims, error)
	IsTokenExpiredFunc    func(tokenString string) bool
	RefreshTokenFunc      func(refreshTokenString string) (*models.TokenPair, error)
	RevokeTokenFunc       func(tokenString string) error
	GetPublicKeyFunc      func() *rsa.PublicKey
	GetPrivateKeyFunc      func() *rsa.PrivateKey
}

func (m *MockJWTService) LoadKeys(privateKeyPath, publicKeyPath string) error {
	if m.LoadKeysFunc != nil {
		return m.LoadKeysFunc(privateKeyPath, publicKeyPath)
	}
	return nil
}

func (m *MockJWTService) GenerateToken(userID, email string, roles []string) (string, error) {
	if m.GenerateTokenFunc != nil {
		return m.GenerateTokenFunc(userID, email, roles)
	}
	return "mock-token", nil
}

func (m *MockJWTService) GenerateTokenPair(userID, email string, roles []string) (*models.TokenPair, error) {
	if m.GenerateTokenPairFunc != nil {
		return m.GenerateTokenPairFunc(userID, email, roles)
	}
	return &models.TokenPair{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		ExpiresIn:    3600,
	}, nil
}

func (m *MockJWTService) ValidateToken(tokenString string) (*models.CustomClaims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(tokenString)
	}
	return &models.CustomClaims{
		UserID: "test-user-id",
		Email:  "test@example.com",
		Roles:  []string{"user"},
	}, nil
}

func (m *MockJWTService) GetClaims(tokenString string) (*models.CustomClaims, error) {
	if m.GetClaimsFunc != nil {
		return m.GetClaimsFunc(tokenString)
	}
	return &models.CustomClaims{
		UserID: "test-user-id",
		Email:  "test@example.com",
		Roles:  []string{"user"},
	}, nil
}

func (m *MockJWTService) IsTokenExpired(tokenString string) bool {
	if m.IsTokenExpiredFunc != nil {
		return m.IsTokenExpiredFunc(tokenString)
	}
	return false
}

func (m *MockJWTService) RefreshToken(refreshTokenString string) (*models.TokenPair, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(refreshTokenString)
	}
	return &models.TokenPair{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresIn:    3600,
	}, nil
}

func (m *MockJWTService) RevokeToken(tokenString string) error {
	if m.RevokeTokenFunc != nil {
		return m.RevokeTokenFunc(tokenString)
	}
	return nil
}

func (m *MockJWTService) GetPublicKey() *rsa.PublicKey {
	if m.GetPublicKeyFunc != nil {
		return m.GetPublicKeyFunc()
	}
	return nil
}

func (m *MockJWTService) GetPrivateKey() *rsa.PrivateKey {
	if m.GetPrivateKeyFunc != nil {
		return m.GetPrivateKeyFunc()
	}
	return nil
}
