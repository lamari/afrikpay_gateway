package models

import "testing"

func TestUser_Validate(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        wantErr bool
    }{
        {
            name: "valid user",
            user: User{
                Email:        "jane@example.com",
                PasswordHash: "hashed_pwd",
                Role:         "user",
            },
            wantErr: false,
        },
        {
            name: "missing email",
            user: User{
                PasswordHash: "hashed_pwd",
                Role:         "user",
            },
            wantErr: true,
        },
        {
            name: "invalid role",
            user: User{
                Email:        "john@example.com",
                PasswordHash: "hashed_pwd",
                Role:         "superadmin",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        if err := tt.user.Validate(); (err != nil) != tt.wantErr {
            t.Errorf("%s: expected error=%v, got %v", tt.name, tt.wantErr, err)
        }
    }
}
