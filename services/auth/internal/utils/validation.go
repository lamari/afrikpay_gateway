package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	
	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	
	// Check length constraints
	if len(email) > 254 {
		return fmt.Errorf("email too long")
	}
	
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	
	if len(password) > 128 {
		return fmt.Errorf("password too long")
	}
	
	// Check for at least one letter and one number
	hasLetter := false
	hasNumber := false
	
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}
	
	if !hasLetter {
		return fmt.Errorf("password must contain at least one letter")
	}
	
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	
	return nil
}

// ValidateUserID validates user ID format
func ValidateUserID(userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	
	// Check if it's a valid UUID format or alphanumeric
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	
	if !uuidRegex.MatchString(userID) && !alphanumericRegex.MatchString(userID) {
		return fmt.Errorf("invalid user ID format")
	}
	
	if len(userID) > 50 {
		return fmt.Errorf("user ID too long")
	}
	
	return nil
}

// ValidateRole validates role name
func ValidateRole(role string) error {
	if role == "" {
		return fmt.Errorf("role cannot be empty")
	}
	
	// Valid roles for the system
	validRoles := map[string]bool{
		"admin":    true,
		"user":     true,
		"operator": true,
		"viewer":   true,
	}
	
	if !validRoles[role] {
		return fmt.Errorf("invalid role: %s", role)
	}
	
	return nil
}

// SanitizeInput sanitizes user input
func SanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Replace multiple whitespace with single space
	re := regexp.MustCompile(`\s+`)
	input = re.ReplaceAllString(input, " ")
	
	// Remove control characters except newlines and tabs
	var result strings.Builder
	for _, r := range input {
		if unicode.IsPrint(r) || r == '\n' || r == '\t' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// ValidateJSONStructure validates if a string is valid JSON
func ValidateJSONStructure(data string) error {
	if data == "" {
		return fmt.Errorf("JSON data cannot be empty")
	}
	
	var js interface{}
	if err := json.Unmarshal([]byte(data), &js); err != nil {
		return fmt.Errorf("invalid JSON structure: %w", err)
	}
	
	return nil
}

// NormalizeWhitespace normalizes whitespace in a string
func NormalizeWhitespace(input string) string {
	// Replace all whitespace sequences with single spaces
	re := regexp.MustCompile(`\s+`)
	normalized := re.ReplaceAllString(input, " ")
	
	// Trim leading and trailing whitespace
	return strings.TrimSpace(normalized)
}

// ValidateRequiredFields validates that required fields are not empty
func ValidateRequiredFields(fields map[string]string) error {
	var missingFields []string
	
	for fieldName, fieldValue := range fields {
		if strings.TrimSpace(fieldValue) == "" {
			missingFields = append(missingFields, fieldName)
		}
	}
	
	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}
	
	return nil
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(uuid)
}

// IsValidEmail checks if an email is valid
func IsValidEmail(email string) bool {
	return ValidateEmail(email) == nil
}

// IsValidPassword checks if a password is valid
func IsValidPassword(password string) bool {
	return ValidatePassword(password) == nil
}

// IsValidUserID checks if a user ID is valid
func IsValidUserID(userID string) bool {
	return ValidateUserID(userID) == nil
}

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	return ValidateRole(role) == nil
}
