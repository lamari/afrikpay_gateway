package repositories

import (
    "context"
    "testing"

    "github.com/afrikpay/gateway/services/crud/internal/models"
)

func TestInMemoryUserRepository_CRUD(t *testing.T) {
    repo := NewInMemoryUserRepository()
    ctx := context.Background()

    user := &models.User{ID: "u1", Email: "test@example.com", PasswordHash: "hash", Role: "user"}
    if err := repo.Create(ctx, user); err != nil {
        t.Fatalf("create error: %v", err)
    }

    got, err := repo.GetByID(ctx, "u1")
    if err != nil || got.Email != "test@example.com" {
        t.Fatalf("get error: %v", err)
    }

    user.Email = "changed@example.com"
    if err := repo.Update(ctx, user); err != nil {
        t.Fatalf("update error: %v", err)
    }
    got, _ = repo.GetByID(ctx, "u1")
    if got.Email != "changed@example.com" {
        t.Fatalf("update not applied")
    }

    if err := repo.Delete(ctx, "u1"); err != nil {
        t.Fatalf("delete error: %v", err)
    }
    if _, err := repo.GetByID(ctx, "u1"); err != ErrNotFound {
        t.Fatalf("expected not found, got %v", err)
    }
}
