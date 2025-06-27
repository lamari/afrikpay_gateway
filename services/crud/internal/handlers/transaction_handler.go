package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/afrikpay/gateway/services/crud/internal/models"
    "github.com/afrikpay/gateway/services/crud/internal/services"
    "github.com/gorilla/mux"
)

type TransactionHandler struct {
    Service *services.TransactionService
}

func (h *TransactionHandler) Register(r *mux.Router) {
    r.HandleFunc("/transactions", h.create).Methods(http.MethodPost)
    r.HandleFunc("/transactions/{id}", h.get).Methods(http.MethodGet)
    r.HandleFunc("/transactions/{id}", h.update).Methods(http.MethodPut)
    r.HandleFunc("/transactions/{id}", h.delete).Methods(http.MethodDelete)
}

func (h *TransactionHandler) create(w http.ResponseWriter, r *http.Request) {
    var req models.Transaction
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

func (h *TransactionHandler) get(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    tr, err := h.Service.Get(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    _ = json.NewEncoder(w).Encode(tr)
}

func (h *TransactionHandler) update(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    var req models.Transaction
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

func (h *TransactionHandler) delete(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    if err := h.Service.Delete(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
