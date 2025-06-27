package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomClaims represents JWT claims with custom fields
type CustomClaims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// Valid validates the custom claims
func (c *CustomClaims) Valid() error {
	// Validate standard claims manually since RegisteredClaims.Valid() doesn't exist in jwt/v5
	now := time.Now()
	
	// Check expiration
	if c.ExpiresAt != nil && c.ExpiresAt.Before(now) {
		return fmt.Errorf("token is expired")
	}
	
	// Check not before
	if c.NotBefore != nil && c.NotBefore.After(now) {
		return fmt.Errorf("token used before valid")
	}
	
	// Check issued at (shouldn't be in the future)
	if c.IssuedAt != nil && c.IssuedAt.After(now.Add(time.Minute)) {
		return fmt.Errorf("token issued in the future")
	}
	
	// Validate custom fields
	if c.UserID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	
	if c.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	
	if len(c.Roles) == 0 {
		return fmt.Errorf("roles cannot be empty")
	}
	
	return nil
}

// HasRole checks if the claims contain a specific role
func (c *CustomClaims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// IsAdmin checks if the claims contain admin role
func (c *CustomClaims) IsAdmin() bool {
	return c.HasRole("admin")
}

// IsExpired checks if the token is expired
func (c *CustomClaims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return c.ExpiresAt.Before(time.Now())
}

// NewCustomClaims creates new custom claims
func NewCustomClaims(userID, email string, roles []string, issuer, audience string, expiration time.Duration) *CustomClaims {
	now := time.Now()
	tokenID := uuid.New().String()
	
	return &CustomClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			Audience:  []string{audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        tokenID,
		},
	}
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// NewTokenPair creates a new token pair
func NewTokenPair(accessToken, refreshToken string, expiresIn int64) *TokenPair {
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}

// RefreshTokenClaims represents claims for refresh tokens
type RefreshTokenClaims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// Valid validates the refresh token claims
func (r *RefreshTokenClaims) Valid() error {
	// Validate standard claims manually since RegisteredClaims.Valid() doesn't exist in jwt/v5
	now := time.Now()
	
	// Check expiration
	if r.ExpiresAt != nil && r.ExpiresAt.Before(now) {
		return fmt.Errorf("refresh token is expired")
	}
	
	// Check not before
	if r.NotBefore != nil && r.NotBefore.After(now) {
		return fmt.Errorf("refresh token used before valid")
	}
	
	// Check issued at (shouldn't be in the future)
	if r.IssuedAt != nil && r.IssuedAt.After(now.Add(time.Minute)) {
		return fmt.Errorf("refresh token issued in the future")
	}
	
	// Validate custom fields
	if r.UserID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	
	if r.TokenType != "refresh" {
		return fmt.Errorf("invalid token type")
	}
	
	return nil
}

// NewRefreshTokenClaims creates new refresh token claims
func NewRefreshTokenClaims(userID, issuer, audience string, expiration time.Duration) *RefreshTokenClaims {
	now := time.Now()
	tokenID := uuid.New().String()
	
	return &RefreshTokenClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			Audience:  []string{audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        tokenID,
		},
	}
}
