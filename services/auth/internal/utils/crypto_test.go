package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHashPassword_ValidPassword tests password hashing with valid password
func TestHashPassword_ValidPassword(t *testing.T) {
	// Given
	password := "testpassword123"

	// When
	hashedPassword, err := HashPassword(password)

	// Then
	require.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
	assert.True(t, len(hashedPassword) > 50) // bcrypt hashes are typically 60+ chars
}

// TestHashPassword_EmptyPassword tests password hashing with empty password
func TestHashPassword_EmptyPassword(t *testing.T) {
	// Given
	password := ""

	// When
	hashedPassword, err := HashPassword(password)

	// Then
	assert.Error(t, err)
	assert.Empty(t, hashedPassword)
	assert.Contains(t, err.Error(), "password cannot be empty")
}

// TestHashPassword_LongPassword tests password hashing with very long password (should fail)
func TestHashPassword_LongPassword(t *testing.T) {
	// Given
	password := string(make([]byte, 100)) // 100 character password (exceeds bcrypt 72 byte limit)

	// When
	hashedPassword, err := HashPassword(password)

	// Then
	require.Error(t, err)
	assert.Empty(t, hashedPassword)
	assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
}

// TestVerifyPassword_ValidPassword tests password verification with correct password
func TestVerifyPassword_ValidPassword(t *testing.T) {
	// Given
	password := "testpassword123"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	// When
	isValid := VerifyPassword(hashedPassword, password)

	// Then
	assert.True(t, isValid)
}

// TestVerifyPassword_InvalidPassword tests password verification with incorrect password
func TestVerifyPassword_InvalidPassword(t *testing.T) {
	// Given
	password := "testpassword123"
	wrongPassword := "wrongpassword"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	// When
	isValid := VerifyPassword(hashedPassword, wrongPassword)

	// Then
	assert.False(t, isValid)
}

// TestVerifyPassword_EmptyPassword tests password verification with empty password
func TestVerifyPassword_EmptyPassword(t *testing.T) {
	// Given
	password := "testpassword123"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	// When
	isValid := VerifyPassword(hashedPassword, "")

	// Then
	assert.False(t, isValid)
}

// TestVerifyPassword_EmptyHash tests password verification with empty hash
func TestVerifyPassword_EmptyHash(t *testing.T) {
	// Given
	password := "testpassword123"

	// When
	isValid := VerifyPassword("", password)

	// Then
	assert.False(t, isValid)
}

// TestVerifyPassword_InvalidHash tests password verification with invalid hash
func TestVerifyPassword_InvalidHash(t *testing.T) {
	// Given
	password := "testpassword123"
	invalidHash := "invalid.hash.format"

	// When
	isValid := VerifyPassword(invalidHash, password)

	// Then
	assert.False(t, isValid)
}

// TestGenerateRandomString_ValidLength tests random string generation with valid length
func TestGenerateRandomString_ValidLength(t *testing.T) {
	// Given
	length := 32

	// When
	randomString, err := GenerateRandomString(length)

	// Then
	require.NoError(t, err)
	assert.Len(t, randomString, length)
	assert.NotEmpty(t, randomString)
}

// TestGenerateRandomString_ZeroLength tests random string generation with zero length
func TestGenerateRandomString_ZeroLength(t *testing.T) {
	// Given
	length := 0

	// When
	randomString, err := GenerateRandomString(length)

	// Then
	assert.Error(t, err)
	assert.Empty(t, randomString)
	assert.Contains(t, err.Error(), "length must be positive")
}

// TestGenerateRandomString_NegativeLength tests random string generation with negative length
func TestGenerateRandomString_NegativeLength(t *testing.T) {
	// Given
	length := -5

	// When
	randomString, err := GenerateRandomString(length)

	// Then
	assert.Error(t, err)
	assert.Empty(t, randomString)
	assert.Contains(t, err.Error(), "length must be positive")
}

// TestGenerateRandomString_Uniqueness tests that generated strings are unique
func TestGenerateRandomString_Uniqueness(t *testing.T) {
	// Given
	length := 16
	iterations := 100
	generated := make(map[string]bool)

	// When
	for i := 0; i < iterations; i++ {
		randomString, err := GenerateRandomString(length)
		require.NoError(t, err)
		
		// Then
		assert.False(t, generated[randomString], "Generated string should be unique")
		generated[randomString] = true
		assert.Len(t, randomString, length)
	}
}

// TestGenerateRandomString_CharacterSet tests that generated string contains valid characters
func TestGenerateRandomString_CharacterSet(t *testing.T) {
	// Given
	length := 100
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// When
	randomString, err := GenerateRandomString(length)

	// Then
	require.NoError(t, err)
	assert.Len(t, randomString, length)

	// Verify all characters are from valid set
	for _, char := range randomString {
		assert.Contains(t, validChars, string(char))
	}
}

// TestGenerateSecureToken_ValidLength tests secure token generation
func TestGenerateSecureToken_ValidLength(t *testing.T) {
	// Given
	length := 32

	// When
	token, err := GenerateSecureToken(length)

	// Then
	require.NoError(t, err)
	assert.Len(t, token, length*2) // Hex encoding doubles the length
	assert.NotEmpty(t, token)
}

// TestGenerateSecureToken_ZeroLength tests secure token generation with zero length
func TestGenerateSecureToken_ZeroLength(t *testing.T) {
	// Given
	length := 0

	// When
	token, err := GenerateSecureToken(length)

	// Then
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "length must be positive")
}

// TestGenerateSecureToken_Uniqueness tests that generated tokens are unique
func TestGenerateSecureToken_Uniqueness(t *testing.T) {
	// Given
	length := 16
	iterations := 100
	generated := make(map[string]bool)

	// When
	for i := 0; i < iterations; i++ {
		token, err := GenerateSecureToken(length)
		require.NoError(t, err)
		
		// Then
		assert.False(t, generated[token], "Generated token should be unique")
		generated[token] = true
		assert.Len(t, token, length*2) // Hex encoding
	}
}

// TestGenerateSecureToken_HexFormat tests that generated token is valid hex
func TestGenerateSecureToken_HexFormat(t *testing.T) {
	// Given
	length := 16

	// When
	token, err := GenerateSecureToken(length)

	// Then
	require.NoError(t, err)
	assert.Len(t, token, length*2)

	// Verify all characters are valid hex
	validHexChars := "0123456789abcdef"
	for _, char := range token {
		assert.Contains(t, validHexChars, string(char))
	}
}

// TestConstantTimeCompare_EqualStrings tests constant time comparison with equal strings
func TestConstantTimeCompare_EqualStrings(t *testing.T) {
	// Given
	str1 := "test_string_123"
	str2 := "test_string_123"

	// When
	isEqual := ConstantTimeCompare(str1, str2)

	// Then
	assert.True(t, isEqual)
}

// TestConstantTimeCompare_DifferentStrings tests constant time comparison with different strings
func TestConstantTimeCompare_DifferentStrings(t *testing.T) {
	// Given
	str1 := "test_string_123"
	str2 := "test_string_456"

	// When
	isEqual := ConstantTimeCompare(str1, str2)

	// Then
	assert.False(t, isEqual)
}

// TestConstantTimeCompare_DifferentLengths tests constant time comparison with different lengths
func TestConstantTimeCompare_DifferentLengths(t *testing.T) {
	// Given
	str1 := "short"
	str2 := "much_longer_string"

	// When
	isEqual := ConstantTimeCompare(str1, str2)

	// Then
	assert.False(t, isEqual)
}

// TestConstantTimeCompare_EmptyStrings tests constant time comparison with empty strings
func TestConstantTimeCompare_EmptyStrings(t *testing.T) {
	// Given
	str1 := ""
	str2 := ""

	// When
	isEqual := ConstantTimeCompare(str1, str2)

	// Then
	assert.True(t, isEqual)
}

// TestConstantTimeCompare_OneEmpty tests constant time comparison with one empty string
func TestConstantTimeCompare_OneEmpty(t *testing.T) {
	// Given
	str1 := "not_empty"
	str2 := ""

	// When
	isEqual := ConstantTimeCompare(str1, str2)

	// Then
	assert.False(t, isEqual)
}
