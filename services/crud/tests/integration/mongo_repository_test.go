//go:build integration
// +build integration

package integration

import (
    "context"
    "os"
    "testing"
    "time"

    "github.com/afrikpay/gateway/services/crud/internal/models"
    "github.com/afrikpay/gateway/services/crud/internal/repositories"
)

func TestMongoRepositories_CRUD(t *testing.T) {
    // Use MongoDB from Docker Compose (should be running on localhost:27017)
    uri := os.Getenv("MONGO_URI")
    if uri == "" {
        uri = "mongodb://admin:password123@localhost:27017/afrikpay?authSource=admin"
    }

    ctx := context.Background()
    client := repositories.MustNewMongoClient(uri)
    db := client.Database("afrikpay_test")

    // Clean up before test
    db.Drop(ctx)
    
    // Wait a bit for cleanup
    time.Sleep(100 * time.Millisecond)

    userRepo := repositories.NewMongoUserRepository(db.Collection("users"))
    walletRepo := repositories.NewMongoWalletRepository(db.Collection("wallets"))
    trxRepo := repositories.NewMongoTransactionRepository(db.Collection("transactions"))

    // Test User CRUD
    user := &models.User{ID: "u1", Email: "test@example.com", Role: "user", PasswordHash: "hash123"}
    if err := userRepo.Create(ctx, user); err != nil {
        t.Fatalf("user create: %v", err)
    }
    
    retrievedUser, err := userRepo.GetByID(ctx, "u1")
    if err != nil {
        t.Fatalf("user get: %v", err)
    }
    if retrievedUser.Email != "test@example.com" {
        t.Errorf("expected email test@example.com, got %s", retrievedUser.Email)
    }

    // Test Wallet CRUD
    wallet := &models.Wallet{ID: "w1", UserID: "u1", Balance: 100.50, Currency: "USDT"}
    if err := walletRepo.Create(ctx, wallet); err != nil {
        t.Fatalf("wallet create: %v", err)
    }
    
    retrievedWallet, err := walletRepo.GetByID(ctx, "w1")
    if err != nil {
        t.Fatalf("wallet get: %v", err)
    }
    if retrievedWallet.Balance != 100.50 {
        t.Errorf("expected balance 100.50, got %f", retrievedWallet.Balance)
    }

    // Test Transaction CRUD
    trx := &models.Transaction{ID: "t1", WalletID: "w1", Amount: 10.25, Type: "deposit", Status: "completed"}
    if err := trxRepo.Create(ctx, trx); err != nil {
        t.Fatalf("transaction create: %v", err)
    }
    
    retrievedTrx, err := trxRepo.GetByID(ctx, "t1")
    if err != nil {
        t.Fatalf("transaction get: %v", err)
    }
    if retrievedTrx.Amount != 10.25 {
        t.Errorf("expected amount 10.25, got %f", retrievedTrx.Amount)
    }
    
    t.Log("All MongoDB repository tests passed successfully!")
}
