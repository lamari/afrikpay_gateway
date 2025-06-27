//go:build integration
// +build integration

package integration

import (
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "github.com/afrikpay/gateway/services/crud/internal/handlers"
    "github.com/afrikpay/gateway/services/crud/internal/middleware"
    "github.com/afrikpay/gateway/services/crud/internal/repositories"
    "github.com/afrikpay/gateway/services/crud/internal/services"
    "github.com/gorilla/mux"
)

func buildRouter() *mux.Router {
    // in-memory repos
    userRepo := repositories.NewInMemoryUserRepository()
    walletRepo := repositories.NewInMemoryWalletRepository()
    trxRepo := repositories.NewInMemoryTransactionRepository()

    userSvc := &services.UserService{Repo: userRepo}
    walletSvc := &services.WalletService{Repo: walletRepo}
    trxSvc := &services.TransactionService{Repo: trxRepo}

    r := mux.NewRouter()
    r.Use(middleware.AuthMiddleware())
    (&handlers.UserHandler{Service: userSvc}).Register(r)
    (&handlers.WalletHandler{Service: walletSvc}).Register(r)
    (&handlers.TransactionHandler{Service: trxSvc}).Register(r)
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
    return r
}

func mockAuthServer() *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/auth/verify" {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        auth := r.Header.Get("Authorization")
        if auth == "Bearer valid" {
            w.WriteHeader(http.StatusOK)
        } else {
            w.WriteHeader(http.StatusUnauthorized)
        }
    }))
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
    authSrv := mockAuthServer()
    defer authSrv.Close()
    os.Setenv("AUTH_URL", authSrv.URL)

    router := buildRouter()
    srv := httptest.NewServer(router)
    defer srv.Close()

    req, _ := http.NewRequest(http.MethodGet, srv.URL+"/health", nil)
    req.Header.Set("Authorization", "Bearer valid")
    resp, err := srv.Client().Do(req)
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200, got %d", resp.StatusCode)
    }
}

func TestAuthMiddleware_MissingOrInvalidToken(t *testing.T) {
    authSrv := mockAuthServer()
    defer authSrv.Close()
    os.Setenv("AUTH_URL", authSrv.URL)

    router := buildRouter()
    srv := httptest.NewServer(router)
    defer srv.Close()

    // missing token
    resp, err := srv.Client().Get(srv.URL + "/health")
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    if resp.StatusCode != http.StatusUnauthorized {
        t.Fatalf("expected 401, got %d", resp.StatusCode)
    }

    // invalid token
    req, _ := http.NewRequest(http.MethodGet, srv.URL+"/health", nil)
    req.Header.Set("Authorization", "Bearer invalid")
    resp, err = srv.Client().Do(req)
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    if resp.StatusCode != http.StatusUnauthorized {
        t.Fatalf("expected 401, got %d", resp.StatusCode)
    }
}
