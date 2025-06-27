package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	
	// Use bcrypt with cost 12 for good security/performance balance
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(hashedPassword, password string) bool {
	if hashedPassword == "" || password == "" {
		return false
	}
	
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}
	
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		result[i] = charset[num.Int64()]
	}
	
	return string(result), nil
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}
	
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	
	return hex.EncodeToString(bytes), nil
}

// ConstantTimeCompare performs constant-time string comparison to prevent timing attacks
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
