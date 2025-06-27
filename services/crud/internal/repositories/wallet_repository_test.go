package repositories

import (
    "context"
    "testing"

    "github.com/afrikpay/gateway/services/crud/internal/models"
)

func TestInMemoryWalletRepository_CRUD(t *testing.T) {
    repo := NewInMemoryWalletRepository()
    ctx := context.Background()

    w := &models.Wallet{ID: "w1", UserID: "u1", Balance: 0, Currency: "USDT"}
    if err := repo.Create(ctx, w); err != nil {
        t.Fatalf("create error: %v", err)
    }

    got, err := repo.GetByID(ctx, "w1")
    if err != nil || got.UserID != "u1" {
        t.Fatalf("get error: %v", err)
    }

    w.Balance = 10
    if err := repo.Update(ctx, w); err != nil {
        t.Fatalf("update error: %v", err)
    }
    got, _ = repo.GetByID(ctx, "w1")
    if got.Balance != 10 {
        t.Fatalf("balance not updated")
    }

    if err := repo.Delete(ctx, "w1"); err != nil {
        t.Fatalf("delete error: %v", err)
    }
    if _, err := repo.GetByID(ctx, "w1"); err != ErrNotFound {
        t.Fatalf("expected not found, got %v", err)
    }
}
