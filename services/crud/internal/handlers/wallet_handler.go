package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/afrikpay/gateway/services/crud/internal/models"
    "github.com/afrikpay/gateway/services/crud/internal/services"
    "github.com/gorilla/mux"
)

type WalletHandler struct {
    Service *services.WalletService
}

func (h *WalletHandler) Register(r *mux.Router) {
    // API v1 prefix
    api := r.PathPrefix("/api/v1").Subrouter()
    
    // Standard CRUD operations
    api.HandleFunc("/wallets", h.create).Methods(http.MethodPost)
    api.HandleFunc("/wallets/{id}", h.get).Methods(http.MethodGet)
    api.HandleFunc("/wallets/{id}", h.update).Methods(http.MethodPut)
    api.HandleFunc("/wallets/{id}", h.delete).Methods(http.MethodDelete)
    
    // Additional endpoints for internal CRUD client
    api.HandleFunc("/wallets/{id}/balance", h.updateBalance).Methods(http.MethodPatch)
    api.HandleFunc("/users/{userId}/wallets", h.getUserWallet).Methods(http.MethodGet)
    
    // For backward compatibility, keep the old routes without prefix
    r.HandleFunc("/wallets", h.create).Methods(http.MethodPost)
    r.HandleFunc("/wallets/{id}", h.get).Methods(http.MethodGet)
    r.HandleFunc("/wallets/{id}", h.update).Methods(http.MethodPut)
    r.HandleFunc("/wallets/{id}", h.delete).Methods(http.MethodDelete)
}

func (h *WalletHandler) create(w http.ResponseWriter, r *http.Request) {
    var req models.Wallet
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := h.Service.Create(r.Context(), &req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(req)
}

func (h *WalletHandler) get(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    wallet, err := h.Service.Get(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    _ = json.NewEncoder(w).Encode(wallet)
}

func (h *WalletHandler) update(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    var req models.Wallet
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    req.ID = id
    if err := h.Service.Update(r.Context(), &req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    _ = json.NewEncoder(w).Encode(req)
}

func (h *WalletHandler) delete(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    if err := h.Service.Delete(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// updateBalance updates only the balance of a wallet
func (h *WalletHandler) updateBalance(w http.ResponseWriter, r *http.Request) {
    type BalanceUpdate struct {
        Amount   float64 `json:"amount"`
        Currency string  `json:"currency"`
    }

    id := mux.Vars(r)["id"]
    var update BalanceUpdate
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    wallet, err := h.Service.Get(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    if wallet.Currency != update.Currency {
        http.Error(w, "currency mismatch", http.StatusBadRequest)
        return
    }

    wallet.Balance += update.Amount
    if err := h.Service.Update(r.Context(), wallet); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _ = json.NewEncoder(w).Encode(wallet)
}

// getUserWallet gets a wallet by user ID and currency
func (h *WalletHandler) getUserWallet(w http.ResponseWriter, r *http.Request) {
    userId := mux.Vars(r)["userId"]
    currency := r.URL.Query().Get("currency")
    
    if currency == "" {
        http.Error(w, "currency parameter is required", http.StatusBadRequest)
        return
    }

    // We need to extend the WalletService to support this query
    wallet, err := h.Service.GetByUserIdAndCurrency(r.Context(), userId, currency)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    _ = json.NewEncoder(w).Encode(wallet)
}
