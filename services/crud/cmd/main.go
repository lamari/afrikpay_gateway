package main

import (
    "log"
    "net/http"
    "os"

    "github.com/afrikpay/gateway/services/crud/internal/config"
    "github.com/afrikpay/gateway/services/crud/internal/handlers"
    "github.com/afrikpay/gateway/services/crud/internal/middleware"
    "github.com/afrikpay/gateway/services/crud/internal/repositories"
    "github.com/afrikpay/gateway/services/crud/internal/services"
    "github.com/gorilla/mux"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8002"
    }

        var (
        userRepo repositories.UserRepository
        walletRepo repositories.WalletRepository
        trxRepo   repositories.TransactionRepository
    )
    if os.Getenv("MONGO_URI") != "" {
        cli := config.MongoClient()
        db := cli.Database("afrikpay")
        userRepo = repositories.NewMongoUserRepository(db.Collection("users"))
        walletRepo = repositories.NewMongoWalletRepository(db.Collection("wallets"))
        trxRepo = repositories.NewMongoTransactionRepository(db.Collection("transactions"))
    } else {
        userRepo = repositories.NewInMemoryUserRepository()
        walletRepo = repositories.NewInMemoryWalletRepository()
        trxRepo = repositories.NewInMemoryTransactionRepository()
    }

    userSvc := &services.UserService{Repo: userRepo}
    walletSvc := &services.WalletService{Repo: walletRepo}
    trxSvc := &services.TransactionService{Repo: trxRepo}

    r := mux.NewRouter()
    r.Use(middleware.AuthMiddleware())
    (&handlers.UserHandler{Service: userSvc}).Register(r)
    (&handlers.WalletHandler{Service: walletSvc}).Register(r)
    (&handlers.TransactionHandler{Service: trxSvc}).Register(r)

    r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("healthy")) })

    log.Printf("CRUD service listening on :%s", port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatal(err)
    }
}
