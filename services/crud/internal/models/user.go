package models

import (
    "fmt"
    "regexp"
    "strings"
    "time"
)

// User represents an end-user of Afrikpay.
type User struct {
    ID           string    `bson:"_id,omitempty" json:"id"`
    Email        string    `bson:"email" json:"email"`
    PasswordHash string    `bson:"password_hash" json:"password_hash"`
    Role         string    `bson:"role" json:"role"`
    CreatedAt    time.Time `bson:"created_at" json:"created_at"`
    UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}

var (
    emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    roleSet    = map[string]struct{}{ // allowed roles
        "admin":    {},
        "user":     {},
        "operator": {},
        "viewer":   {},
    }
)

// Validate performs business-level validation and returns an error if the entity is not valid.
func (u *User) Validate() error {
    if strings.TrimSpace(u.Email) == "" {
        return fmt.Errorf("email cannot be empty")
    }
    if !emailRegex.MatchString(u.Email) {
        return fmt.Errorf("invalid email format")
    }
    if _, ok := roleSet[u.Role]; !ok {
        return fmt.Errorf("invalid role: %s", u.Role)
    }
    if u.PasswordHash == "" {
        return fmt.Errorf("password hash cannot be empty")
    }
    return nil
}
