package models

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestCustomClaims_Valid_ValidClaims tests custom claims validation with valid data
func TestCustomClaims_Valid_ValidClaims(t *testing.T) {
	// Given
	claims := CustomClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			Subject:   "user123",
			Audience:  []string{"test-audience"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "token123",
		},
	}

	// When
	err := claims.Valid()

	// Then
	assert.NoError(t, err)
}

// TestCustomClaims_Valid_EmptyUserID tests custom claims validation with empty user ID
func TestCustomClaims_Valid_EmptyUserID(t *testing.T) {
	// Given
	claims := CustomClaims{
		UserID: "",
		Email:  "test@example.com",
		Roles:  []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			Subject:   "user123",
			Audience:  []string{"test-audience"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "token123",
		},
	}

	// When
	err := claims.Valid()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID cannot be empty")
}

// TestCustomClaims_Valid_EmptyEmail tests custom claims validation with empty email
func TestCustomClaims_Valid_EmptyEmail(t *testing.T) {
	// Given
	claims := CustomClaims{
		UserID: "user123",
		Email:  "",
		Roles:  []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			Subject:   "user123",
			Audience:  []string{"test-audience"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "token123",
		},
	}

	// When
	err := claims.Valid()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email cannot be empty")
}

// TestCustomClaims_Valid_EmptyRoles tests custom claims validation with empty roles
func TestCustomClaims_Valid_EmptyRoles(t *testing.T) {
	// Given
	claims := CustomClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			Subject:   "user123",
			Audience:  []string{"test-audience"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "token123",
		},
	}

	// When
	err := claims.Valid()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "roles cannot be empty")
}

// TestCustomClaims_Valid_ExpiredToken tests custom claims validation with expired token
func TestCustomClaims_Valid_ExpiredToken(t *testing.T) {
	// Given
	claims := CustomClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test-issuer",
			Subject:   "user123",
			Audience:  []string{"test-audience"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ID:        "token123",
		},
	}

	// When
	err := claims.Valid()

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

// TestCustomClaims_HasRole_ExistingRole tests role check with existing role
func TestCustomClaims_HasRole_ExistingRole(t *testing.T) {
	// Given
	claims := CustomClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{"user", "admin"},
	}

	// When
	hasRole := claims.HasRole("admin")

	// Then
	assert.True(t, hasRole)
}

// TestCustomClaims_HasRole_NonExistingRole tests role check with non-existing role
func TestCustomClaims_HasRole_NonExistingRole(t *testing.T) {
	// Given
	claims := CustomClaims{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{"user"},
	}

	// When
	hasRole := claims.HasRole("admin")

	// Then
	assert.False(t, hasRole)
}

// TestCustomClaims_IsAdmin_AdminRole tests admin check with admin role
func TestCustomClaims_IsAdmin_AdminRole(t *testing.T) {
	// Given
	claims := CustomClaims{
		UserID: "user123",
		Email:  "admin@example.com",
		Roles:  []string{"user", "admin"},
	}

	// When
	isAdmin := claims.IsAdmin()

	// Then
	assert.True(t, isAdmin)
}

// TestCustomClaims_IsAdmin_NoAdminRole tests admin check without admin role
func TestCustomClaims_IsAdmin_NoAdminRole(t *testing.T) {
	// Given
	claims := CustomClaims{
		UserID: "user123",
		Email:  "user@example.com",
		Roles:  []string{"user"},
	}

	// When
	isAdmin := claims.IsAdmin()

	// Then
	assert.False(t, isAdmin)
}

// TestCustomClaims_IsExpired_ExpiredClaims tests expiration check with expired claims
func TestCustomClaims_IsExpired_ExpiredClaims(t *testing.T) {
	// Given
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}

	// When
	isExpired := claims.IsExpired()

	// Then
	assert.True(t, isExpired)
}

// TestCustomClaims_IsExpired_ValidClaims tests expiration check with valid claims
func TestCustomClaims_IsExpired_ValidClaims(t *testing.T) {
	// Given
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	// When
	isExpired := claims.IsExpired()

	// Then
	assert.False(t, isExpired)
}

// TestCustomClaims_GetExpirationTime tests getting expiration time from claims
func TestCustomClaims_GetExpirationTime(t *testing.T) {
	// Given
	expirationTime := time.Now().Add(time.Hour)
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// When
	expTime := claims.ExpiresAt.Time

	// Then
	assert.True(t, expTime.Equal(expirationTime.Truncate(time.Second)))
}

// TestCustomClaims_GetIssuedAt tests getting issued at time from claims
func TestCustomClaims_GetIssuedAt(t *testing.T) {
	// Given
	issuedTime := time.Now()
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(issuedTime),
		},
	}

	// When
	iat := claims.IssuedAt.Time

	// Then
	assert.True(t, iat.Equal(issuedTime.Truncate(time.Second)))
}

// TestCustomClaims_GetNotBefore tests getting not before time from claims
func TestCustomClaims_GetNotBefore(t *testing.T) {
	// Given
	notBeforeTime := time.Now()
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(notBeforeTime),
		},
	}

	// When
	nbf := claims.NotBefore.Time

	// Then
	assert.True(t, nbf.Equal(notBeforeTime.Truncate(time.Second)))
}

// TestNewCustomClaims tests creating new custom claims
func TestNewCustomClaims(t *testing.T) {
	// Given
	userID := "user123"
	email := "test@example.com"
	roles := []string{"user", "admin"}
	issuer := "test-issuer"
	audience := "test-audience"
	expiration := 24 * time.Hour

	// When
	claims := NewCustomClaims(userID, email, roles, issuer, audience, expiration)

	// Then
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, issuer, claims.Issuer)
	assert.Equal(t, []string{audience}, []string(claims.Audience))
	assert.Equal(t, userID, claims.Subject)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
	assert.True(t, claims.IssuedAt.Before(time.Now().Add(time.Second)))
	assert.True(t, claims.NotBefore.Before(time.Now().Add(time.Second)))
	assert.NotEmpty(t, claims.ID)
}
